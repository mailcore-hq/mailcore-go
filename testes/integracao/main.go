// Cliente do laboratório Docker. A saída contém somente resultados sanitizados.
package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	mailcore "github.com/mailcore-hq/mailcore-go"
)

type objeto = map[string]any

type fixture struct {
	Token              string `json:"token"`
	CSRF               string `json:"csrf"`
	Organizacao        string `json:"organizacao"`
	Dominio            string `json:"dominio"`
	NomeDominio        string `json:"nome_dominio"`
	Rota               string `json:"rota"`
	CertificadoCA      string `json:"certificado_ca"`
	CallbackCorpo      string `json:"callback_corpo"`
	CallbackAssinatura string `json:"callback_assinatura"`
	MIMEEsperado       string `json:"mime_esperado"`
}

type falhaIntegracao struct{}

func exigir(condicao bool) {
	if !condicao {
		panic(falhaIntegracao{})
	}
}

func resposta(retorno *mailcore.Resposta, erro error) *mailcore.Resposta {
	if erro != nil {
		panic(erro)
	}
	exigir(retorno != nil)
	return retorno
}

func dados(retorno *mailcore.Resposta, erro error) objeto {
	retorno = resposta(retorno, erro)
	var documento objeto
	if erro := retorno.Decodificar(&documento); erro != nil {
		panic(erro)
	}
	exigir(documento != nil)
	return documento
}

func objetoEm(documento objeto, chave string) objeto {
	valor, ok := documento[chave].(map[string]any)
	exigir(ok)
	return valor
}

func listaEm(documento objeto, chave string) []any {
	valor, ok := documento[chave].([]any)
	exigir(ok)
	return valor
}

func textoEm(documento objeto, chave string) string {
	valor, ok := documento[chave].(string)
	exigir(ok)
	return valor
}

func numeroIgual(valor any, esperado string) bool {
	numero, ok := valor.(json.Number)
	return ok && string(numero) == esperado
}

func esperarErro(status int, executar func() (*mailcore.Resposta, error)) {
	_, erro := executar()
	var api *mailcore.Erro
	exigir(errors.As(erro, &api) && api.Status == status)
}

func ids(documento objeto) []string {
	identificadores := []string{}
	for _, valor := range listaEm(documento, "dados") {
		item, ok := valor.(map[string]any)
		exigir(ok)
		identificadores = append(identificadores, textoEm(item, "id"))
	}
	sort.Strings(identificadores)
	return identificadores
}

type integracao struct {
	fixture     fixture
	certificado string
	http        *http.Client
	clientes    []*mailcore.Cliente
	etapa       string
	cenarios    []string
}

func clienteHTTP(certificado string) *http.Client {
	ca, erro := os.ReadFile(certificado)
	exigir(erro == nil)
	raizes := x509.NewCertPool()
	exigir(raizes.AppendCertsFromPEM(ca))
	transporte := http.DefaultTransport.(*http.Transport).Clone()
	transporte.Proxy = nil
	transporte.DisableCompression = true
	transporte.TLSClientConfig = &tls.Config{RootCAs: raizes, MinVersion: tls.VersionTLS12}
	transporte.DialContext = (&net.Dialer{Timeout: 2 * time.Second}).DialContext
	return &http.Client{
		Transport: transporte,
		Timeout:   10 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (i *integracao) cliente(chave string, injetado bool) *mailcore.Cliente {
	configuracao := mailcore.Configuracao{
		URLBase: "https://mailcore-api:3443/v1", CertificadoCA: i.certificado,
		TempoLimite: 10 * time.Second, TempoConexao: 2 * time.Second,
	}
	if injetado {
		configuracao.ClienteHTTP = i.http
	}
	cliente, erro := mailcore.NovoCliente(chave, configuracao)
	if erro != nil {
		panic(erro)
	}
	i.clientes = append(i.clientes, cliente)
	return cliente
}

// Sessão e callback pertencem apenas à preparação da fixture, fora do SDK.
func (i *integracao) requisitar(ctx context.Context, caminho, metodo string, cabecalhos http.Header, corpo []byte) (int, objeto) {
	pedido, erro := http.NewRequestWithContext(ctx, metodo, "https://mailcore-api:3443"+caminho, bytes.NewReader(corpo))
	exigir(erro == nil)
	pedido.GetBody = nil
	pedido.Header = cabecalhos.Clone()
	retorno, erro := i.http.Do(pedido)
	exigir(erro == nil)
	defer retorno.Body.Close()
	conteudo, erro := io.ReadAll(io.LimitReader(retorno.Body, 1024*1024+1))
	exigir(erro == nil && len(conteudo) <= 1024*1024)
	if len(conteudo) == 0 {
		return retorno.StatusCode, nil
	}
	var documento objeto
	decoder := json.NewDecoder(bytes.NewReader(conteudo))
	decoder.UseNumber()
	exigir(decoder.Decode(&documento) == nil)
	return retorno.StatusCode, documento
}

func (i *integracao) concluir() {
	i.cenarios = append(i.cenarios, i.etapa)
}

func agendamento() string {
	return time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
}

func copia(documento objeto, chave string, valor any) objeto {
	novo := make(objeto, len(documento)+1)
	for nome, original := range documento {
		novo[nome] = original
	}
	novo[chave] = valor
	return novo
}

func (i *integracao) executarCenarios(ctx context.Context) objeto {
	prazo := time.Now().Add(30 * time.Second)
	pronta := false
	for time.Now().Before(prazo) && !pronta {
		func() {
			defer func() { _ = recover() }()
			status, _ := i.requisitar(ctx, "/health/ready", http.MethodGet, http.Header{}, nil)
			pronta = status == http.StatusOK
		}()
		if !pronta {
			select {
			case <-time.After(200 * time.Millisecond):
			case <-ctx.Done():
				panic(ctx.Err())
			}
		}
	}
	exigir(pronta)

	i.etapa = "criar_full_sessao_sintetica"
	corpo, erro := json.Marshal(objeto{"nome": "Fixture Go Full", "permissao": "full_access"})
	exigir(erro == nil)
	status, full := i.requisitar(ctx, "/v1/api-keys", http.MethodPost, http.Header{
		"Content-Type": {"application/json"}, "Origin": {"https://localhost:17443"},
		"Cookie": {"__Host-mailcore=" + i.fixture.Token}, "X-Csrf-Token": {i.fixture.CSRF},
	}, corpo)
	i.fixture.Token, i.fixture.CSRF = "", ""
	exigir(status == http.StatusCreated)
	sdk := i.cliente(textoEm(full, "chave"), false)
	sdkInjetado := i.cliente(textoEm(full, "chave"), true)
	remetente := "avisos@" + i.fixture.NomeDominio

	i.etapa = "chaves_full_sending_escopo_dump"
	chaveEnvio := dados(sdk.ChavesAPI.Criar(ctx, objeto{"nome": "Fixture Go Sending", "permissao": "sending_access", "dominio_id": i.fixture.Dominio}))
	envio := i.cliente(textoEm(chaveEnvio, "chave"), false)
	esperarErro(403, func() (*mailcore.Resposta, error) { return envio.Dominios.Listar(ctx, nil) })
	exigir(!strings.Contains(fmt.Sprintf("%+v", sdk), textoEm(full, "chave")))
	exigir(!strings.Contains(fmt.Sprintf("%#v", sdk), textoEm(full, "chave")))
	encontrada := false
	for _, valor := range listaEm(dados(sdk.ChavesAPI.Listar(ctx)), "dados") {
		item, ok := valor.(map[string]any)
		exigir(ok)
		encontrada = encontrada || item["id"] == chaveEnvio["id"]
	}
	exigir(encontrada)
	i.concluir()

	pedido := objeto{
		"remetente": "Equipe Go <" + remetente + ">", "destinatarios": []string{"principal@example.test"},
		"cc": []string{"copia@example.test"}, "cco": []string{"oculto@example.test"}, "assunto": "Validação Go",
		"texto": "Dados sintéticos", "html": `<p>Dados sintéticos<img src="cid:marca@example.test"></p>`,
		"cabecalhos": objeto{"X-Teste-SDK": "sintetico"}, "agendada_em": agendamento(),
		"anexos": []any{objeto{"nome": "marca.png", "tipo_conteudo": "image/png", "conteudo_base64": base64.StdEncoding.EncodeToString([]byte{0x89, 'P', 'N', 'G'}), "id_conteudo": "marca@example.test"}},
	}
	i.etapa = "envio_composto_idempotencia_409_download_cancelamento_concorrencia"
	opcoes := mailcore.OpcoesRequisicao{ChaveIdempotencia: "go-composto"}
	aceita := dados(envio.Emails.Enviar(ctx, pedido, opcoes))
	exigir(aceita["estado"] == "agendada")
	exigir(reflect.DeepEqual(dados(envio.Emails.Enviar(ctx, pedido, opcoes)), aceita))
	esperarErro(409, func() (*mailcore.Resposta, error) {
		return envio.Emails.Enviar(ctx, copia(pedido, "assunto", "Outra intenção"), opcoes)
	})
	exigir(dados(sdk.Emails.ObterPorIdempotencia(ctx, "go-composto"))["id"] == aceita["id"])
	id := textoEm(aceita, "id")
	exigir(len(listaEm(dados(sdk.Emails.Obter(ctx, id)), "destinatarios")) == 3)
	anexo := resposta(sdk.Emails.BaixarAnexo(ctx, id, 0))
	exigir(bytes.Equal(anexo.Dados, []byte{0x89, 'P', 'N', 'G'}))
	var tarefas sync.WaitGroup
	falhas := make(chan bool, 8)
	for indice := 0; indice < 8; indice++ {
		tarefas.Add(1)
		go func() {
			defer tarefas.Done()
			valido := false
			defer func() { _ = recover(); falhas <- !valido }()
			valido = dados(sdk.Emails.Obter(ctx, id))["id"] == aceita["id"]
		}()
	}
	tarefas.Wait()
	close(falhas)
	for falhou := range falhas {
		exigir(!falhou)
	}
	exigir(dados(sdk.Emails.Cancelar(ctx, id))["estado"] == "cancelada")
	i.concluir()

	i.etapa = "lote_atomico_repeticao_conflito"
	item := objeto{"remetente": remetente, "destinatarios": []string{"lote@example.test"}, "assunto": "Lote Go", "texto": "Dados sintéticos", "agendada_em": agendamento()}
	antes := ids(dados(sdk.Emails.Listar(ctx, mailcore.Consulta{"limite": 100})))
	esperarErro(422, func() (*mailcore.Resposta, error) {
		return envio.Lotes.Enviar(ctx, []any{item, copia(item, "destinatarios", []string{"endereco-invalido"})}, mailcore.OpcoesRequisicao{ChaveIdempotencia: "go-lote-invalido"})
	})
	exigir(reflect.DeepEqual(antes, ids(dados(sdk.Emails.Listar(ctx, mailcore.Consulta{"limite": 100})))))
	lote := []any{item, copia(item, "assunto", "Segundo item Go")}
	opcoesLote := mailcore.OpcoesRequisicao{ChaveIdempotencia: "go-lote"}
	aceitaLote := dados(envio.Lotes.Enviar(ctx, lote, opcoesLote))
	exigir(len(listaEm(aceitaLote, "dados")) == 2)
	exigir(reflect.DeepEqual(dados(envio.Lotes.Enviar(ctx, lote, opcoesLote)), aceitaLote))
	exigir(reflect.DeepEqual(dados(sdk.Lotes.ObterPorIdempotencia(ctx, "go-lote")), aceitaLote))
	esperarErro(409, func() (*mailcore.Resposta, error) {
		return envio.Lotes.Enviar(ctx, []any{copia(item, "assunto", "Divergente"), lote[1]}, opcoesLote)
	})
	exigir(len(ids(dados(sdk.Emails.Listar(ctx, mailcore.Consulta{"limite": 100})))) == len(antes)+2)
	for _, valor := range listaEm(aceitaLote, "dados") {
		email, ok := valor.(map[string]any)
		exigir(ok)
		resposta(sdk.Emails.Cancelar(ctx, textoEm(email, "id")))
	}
	i.concluir()

	i.etapa = "dominios_capacidades_tracking_crud"
	dominio := dados(sdk.Dominios.Atualizar(ctx, i.fixture.Dominio, objeto{"capacidades": objeto{"envio": "enabled", "recebimento": "enabled"}, "rastreamento_aberturas": true, "rastreamento_cliques": true}))
	exigir(objetoEm(dominio, "capacidades")["recebimento"] == "enabled" && dominio["rastreamento_aberturas"] == true)
	novo := dados(sdk.Dominios.Criar(ctx, objeto{"nome": "go-" + strings.ReplaceAll(i.fixture.Organizacao, "-", "") + ".example.com", "capacidades": objeto{"envio": "enabled", "recebimento": "disabled"}}))
	exigir(dados(sdk.Dominios.Remover(ctx, textoEm(novo, "id")))["estado"] == "suspenso")
	i.concluir()

	i.etapa = "contatos_objetos_listas_null_false_importacao_csv_fluxos"
	segmento := dados(sdk.Segmentos.Criar(ctx, objeto{"nome": "Segmento Go"}))
	topico := dados(sdk.Topicos.Criar(ctx, objeto{"nome": "Tópico Go", "visibilidade": "publico", "assinatura_padrao": "opt_in"}))
	contato := dados(sdk.Contatos.Criar(ctx, objeto{"email": "contato@example.test", "primeiro_nome": "Sintético", "descadastrado": false, "segmentos": []string{textoEm(segmento, "id")}, "topicos": []any{objeto{"id": textoEm(topico, "id"), "assinatura": "opt_in"}}}))
	exigir(dados(sdk.Contatos.Obter(ctx, textoEm(contato, "id")))["id"] == contato["id"] && contato["descadastrado"] == false)
	importacao := dados(sdk.Contatos.Importacoes.Criar(ctx, objeto{"arquivo_nome": "fixture.csv", "conteudo": "email,primeiro_nome\ncontato@example.test,Importado\n", "mapeamento": objeto{"0": "email", "1": "primeiro_nome"}}, mailcore.OpcoesRequisicao{ChaveIdempotencia: "go-importacao"}))
	// Sem worker: a operação é aceita e conserva os contadores pendentes.
	exigir(importacao["estado"] == "pendente" && numeroIgual(importacao["total"], "1") && numeroIgual(importacao["processados"], "0"))
	exigir(dados(sdk.Contatos.Atualizar(ctx, textoEm(contato, "id"), objeto{"primeiro_nome": nil}))["primeiro_nome"] == nil)
	exportacao := dados(sdk.Contatos.Exportacoes.Criar(ctx, objeto{}, mailcore.OpcoesRequisicao{ChaveIdempotencia: "go-exportacao"}))
	exigir(exportacao["estado"] == "concluida" && numeroIgual(exportacao["total"], "1"))
	var csv bytes.Buffer
	fluxoCSV := resposta(sdk.Contatos.Exportacoes.BaixarFluxo(ctx, textoEm(exportacao, "id"), func(parte []byte) error { _, erro := csv.Write(parte); return erro }))
	exigir(strings.Contains(csv.String(), "contato@example.test") && fluxoCSV.TotalBytes == int64(csv.Len()))
	var csvInjetado bytes.Buffer
	fluxoCSVInjetado := resposta(sdkInjetado.Contatos.Exportacoes.BaixarFluxo(ctx, textoEm(exportacao, "id"), func(parte []byte) error { _, erro := csvInjetado.Write(parte); return erro }))
	exigir(bytes.Equal(csvInjetado.Bytes(), csv.Bytes()) && fluxoCSVInjetado.TotalBytes == int64(csvInjetado.Len()))
	exigir(len(listaEm(dados(sdk.Contatos.Exportacoes.Listar(ctx, nil)), "dados")) == 1)
	i.concluir()

	i.etapa = "modelos_pastas_versoes_previa_publicacao_cas"
	pasta := dados(sdk.Modelos.Pastas.Criar(ctx, objeto{"nome": "Pasta Go"}))
	modelo := dados(sdk.Modelos.Criar(ctx, objeto{"nome": "Modelo Go", "html": "<p>Olá, {{{NOME}}}.</p>", "pasta_id": pasta["id"], "remetente": remetente, "assunto": "Modelo sintético", "variaveis": []any{objeto{"chave": "NOME", "tipo": "string"}}}))
	versao, ok := listaEm(dados(sdk.Modelos.Versoes.Listar(ctx, textoEm(modelo, "id"), nil)), "dados")[0].(map[string]any)
	exigir(ok)
	versaoObtida := dados(sdk.Modelos.Versoes.Obter(ctx, textoEm(modelo, "id"), textoEm(versao, "id")))
	exigir(versaoObtida["versao_id"] == versao["id"] && versaoObtida["id"] == modelo["id"])
	previa := dados(sdk.Modelos.Prever(ctx, textoEm(modelo, "id"), objeto{"revisao": modelo["revisao"], "variaveis": objeto{"NOME": "Go"}}))
	exigir(strings.Contains(textoEm(previa, "html"), "Go"))
	endereco, erro := url.Parse(textoEm(previa, "url"))
	exigir(erro == nil)
	previaID := path.Base(endereco.Path)
	conteudoPrevia := resposta(sdk.Modelos.ObterPrevia(ctx, i.fixture.Organizacao, previaID))
	exigir(bytes.Contains(conteudoPrevia.Dados, []byte("Go")))
	antiga := modelo["revisao"]
	modelo = dados(sdk.Modelos.Atualizar(ctx, textoEm(modelo, "id"), objeto{"revisao": antiga, "nome": "Modelo editado"}))
	esperarErro(409, func() (*mailcore.Resposta, error) {
		return sdk.Modelos.Atualizar(ctx, textoEm(modelo, "id"), objeto{"revisao": antiga, "nome": "Conflito"})
	})
	modelo = dados(sdk.Modelos.Publicar(ctx, textoEm(modelo, "id"), objeto{"revisao": modelo["revisao"]}))
	exigir(modelo["status"] == "published")
	modelo = dados(sdk.Modelos.Versoes.Restaurar(ctx, textoEm(modelo, "id"), textoEm(versao, "id"), objeto{"revisao": modelo["revisao"]}))
	i.concluir()

	i.etapa = "marketing_objeto_vazio_crud_sem_disparo"
	campanha := dados(sdk.Campanhas.Criar(ctx, objeto{"nome": "Campanha Go"}))
	campanha = dados(sdk.Campanhas.Atualizar(ctx, textoEm(campanha, "id"), objeto{"revisao": campanha["revisao"], "nome": "Campanha editada"}))
	exigir(dados(sdk.Campanhas.Obter(ctx, textoEm(campanha, "id")))["id"] == campanha["id"])
	resposta(sdk.Campanhas.Remover(ctx, textoEm(campanha, "id")))
	evento := dados(sdk.Eventos.Criar(ctx, objeto{"nome": "go.qa", "esquema": objeto{}}))
	exigir(len(objetoEm(evento, "esquema")) == 0)
	exigir(dados(sdk.Eventos.Atualizar(ctx, textoEm(evento, "id"), objeto{"esquema": nil}))["esquema"] == nil)
	automacao := dados(sdk.Automacoes.Criar(ctx, objeto{"nome": "Automação Go", "estado": "inativa"}))
	automacao = dados(sdk.Automacoes.Atualizar(ctx, textoEm(automacao, "id"), objeto{"revisao": automacao["revisao"], "nome": "Automação editada"}))
	exigir(dados(sdk.Automacoes.Obter(ctx, textoEm(automacao, "id")))["id"] == automacao["id"])
	resposta(sdk.Automacoes.Remover(ctx, textoEm(automacao, "id")))
	resposta(sdk.Eventos.Remover(ctx, textoEm(evento, "id")))
	i.concluir()

	i.etapa = "caixas_etiquetas_rascunhos_cas_sem_envio"
	caixa := dados(sdk.Caixas.Criar(ctx, objeto{"endereco": "suporte@" + i.fixture.NomeDominio, "nome": "Suporte Go"}))
	exigir(dados(sdk.Caixas.Obter(ctx, textoEm(caixa, "id")))["id"] == caixa["id"])
	etiqueta := dados(sdk.Caixas.Etiquetas.Criar(ctx, textoEm(caixa, "id"), objeto{"nome": "Pedidos", "cor": "#224466"}))
	rascunho := dados(sdk.Caixas.Rascunhos.Criar(ctx, textoEm(caixa, "id"), objeto{"para": []string{"cliente@example.test"}, "assunto": "Pedido", "texto": "Olá"}))
	antiga = rascunho["revisao"]
	rascunho = dados(sdk.Caixas.Rascunhos.Atualizar(ctx, textoEm(caixa, "id"), textoEm(rascunho, "id"), objeto{"revisao": antiga, "texto": "Atualizado"}))
	esperarErro(409, func() (*mailcore.Resposta, error) {
		return sdk.Caixas.Rascunhos.Atualizar(ctx, textoEm(caixa, "id"), textoEm(rascunho, "id"), objeto{"revisao": antiga, "texto": "Conflito"})
	})
	exigir(dados(sdk.Caixas.Rascunhos.Obter(ctx, textoEm(caixa, "id"), textoEm(rascunho, "id")))["html"] == nil)
	resposta(sdk.Caixas.Rascunhos.Remover(ctx, textoEm(caixa, "id"), textoEm(rascunho, "id")))
	resposta(sdk.Caixas.Etiquetas.Remover(ctx, textoEm(caixa, "id"), textoEm(etiqueta, "id")))
	resposta(sdk.Caixas.Remover(ctx, textoEm(caixa, "id")))
	i.concluir()

	i.etapa = "recebimento_mime_real_pendente_bytes_fluxos_padrao_injetado"
	exigir(len(listaEm(dados(sdk.Recebimento.Listar(ctx, nil)), "dados")) == 0)
	caminho := "/interno/postal/recebidos/" + i.fixture.Rota
	cabecalhos := http.Header{"Content-Type": {"application/json"}, "X-Postal-Signature-256": {i.fixture.CallbackAssinatura}}
	status, callback := i.requisitar(ctx, caminho, http.MethodPost, cabecalhos, []byte(i.fixture.CallbackCorpo))
	exigir(status == http.StatusOK)
	recebido := textoEm(callback, "id")
	status, repetida := i.requisitar(ctx, caminho, http.MethodPost, cabecalhos, []byte(i.fixture.CallbackCorpo))
	exigir(status == http.StatusOK && repetida["id"] == recebido)
	mime, erro := base64.StdEncoding.Strict().DecodeString(i.fixture.MIMEEsperado)
	exigir(erro == nil)
	detalhe := dados(sdk.Recebimento.Obter(ctx, recebido, mailcore.Consulta{"formato_html": "cid"}))
	exigir(detalhe["estado_processamento"] == "pendente" && reflect.DeepEqual(listaEm(detalhe, "para"), []any{"real@" + i.fixture.NomeDominio}))
	original := resposta(sdk.Recebimento.BaixarOriginal(ctx, recebido))
	exigir(original.Status == http.StatusOK && original.Cabecalhos.Get("Content-Type") == "message/rfc822" && bytes.Equal(original.Dados, mime))
	var bruto bytes.Buffer
	fluxo := resposta(sdk.Recebimento.BaixarOriginalFluxo(ctx, recebido, func(parte []byte) error { _, erro := bruto.Write(parte); return erro }))
	exigir(bytes.Equal(bruto.Bytes(), mime) && fluxo.TotalBytes == int64(len(mime)))
	var brutoInjetado bytes.Buffer
	fluxoInjetado := resposta(sdkInjetado.Recebimento.BaixarOriginalFluxo(ctx, recebido, func(parte []byte) error { _, erro := brutoInjetado.Write(parte); return erro }))
	exigir(bytes.Equal(brutoInjetado.Bytes(), mime) && fluxoInjetado.TotalBytes == int64(len(mime)) && fluxoInjetado.Cabecalhos.Get("Content-Type") == "message/rfc822")
	exigir(len(listaEm(dados(sdk.Recebimento.Listar(ctx, nil)), "dados")) == 1)
	esperarErro(403, func() (*mailcore.Resposta, error) { return envio.Recebimento.Listar(ctx, nil) })
	i.concluir()

	i.etapa = "cli_revogacao_e_limpeza_crud"
	exigir(len(listaEm(dados(sdk.CLI.Concessoes.Listar(ctx, mailcore.Consulta{"limite": 1})), "dados")) == 0)
	esperarErro(403, func() (*mailcore.Resposta, error) { return envio.CLI.Concessoes.Listar(ctx, nil) })
	resposta(sdk.Contatos.Remover(ctx, textoEm(contato, "id")))
	resposta(sdk.Segmentos.Remover(ctx, textoEm(segmento, "id")))
	resposta(sdk.Topicos.Remover(ctx, textoEm(topico, "id")))
	resposta(sdk.Modelos.Remover(ctx, textoEm(modelo, "id")))
	resposta(sdk.Modelos.Pastas.Remover(ctx, textoEm(pasta, "id")))
	resposta(sdk.ChavesAPI.Revogar(ctx, textoEm(chaveEnvio, "id")))
	esperarErro(401, func() (*mailcore.Resposta, error) {
		return envio.Emails.Enviar(ctx, pedido, mailcore.OpcoesRequisicao{ChaveIdempotencia: "go-revogada"})
	})
	i.concluir()

	return objeto{
		"resultado": "aprovado", "go": runtime.Version(), "cenarios": i.cenarios,
		"recebimento_processado_por_worker": false, "importacao_processada_por_worker": false,
		"callbacks_fluxo_padrao": 2, "callbacks_fluxo_http_injetado": 2,
		"http_injetado_provedor": "net/http", "leituras_concorrentes_mesmo_cliente": 8,
		"politica_http_injetado": "tls_ca_explicita_sem_proxy_redirect_retry_sdk_timeout_total",
	}
}

func (i *integracao) executar(ctx context.Context) (resultado objeto) {
	defer func() {
		for indice := len(i.clientes) - 1; indice >= 0; indice-- {
			i.clientes[indice].Fechar()
		}
		i.http.CloseIdleConnections()
		if falha := recover(); falha != nil {
			resultado = objeto{"resultado": "falhou", "go": runtime.Version(), "etapa": i.etapa, "cenarios": i.cenarios, "falha": "Contrato de integração não satisfeito; detalhes privados omitidos."}
			if erro, ok := falha.(error); ok {
				var api *mailcore.Erro
				if errors.As(erro, &api) {
					resultado["erro_http"] = objeto{"codigo": api.Codigo, "status": api.Status}
				}
			}
		}
	}()
	return i.executarCenarios(ctx)
}

func iniciar() (resultado objeto) {
	defer func() {
		if recover() != nil {
			resultado = objeto{"resultado": "falhou", "etapa": "configuracao", "falha": "Não foi possível iniciar o cliente sintético."}
		}
	}()
	var entrada fixture
	decoder := json.NewDecoder(io.LimitReader(os.Stdin, 1024*1024+1))
	exigir(decoder.Decode(&entrada) == nil)
	exigir(entrada.CertificadoCA != "" && entrada.Organizacao != "" && entrada.NomeDominio != "")
	certificado, erro := os.CreateTemp("", "mailcore-go-ca-")
	exigir(erro == nil)
	defer os.Remove(certificado.Name())
	_, erro = certificado.WriteString(entrada.CertificadoCA)
	exigir(erro == nil)
	exigir(certificado.Close() == nil)
	entrada.CertificadoCA = ""
	httpCliente := clienteHTTP(certificado.Name())
	ensaio := &integracao{fixture: entrada, certificado: certificado.Name(), http: httpCliente, etapa: "prontidao", cenarios: []string{}}
	ctx, cancelar := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelar()
	return ensaio.executar(ctx)
}

func main() {
	resultado := iniciar()
	if erro := json.NewEncoder(os.Stdout).Encode(resultado); erro != nil {
		os.Exit(1)
	}
	if resultado["resultado"] != "aprovado" {
		os.Exit(1)
	}
}
