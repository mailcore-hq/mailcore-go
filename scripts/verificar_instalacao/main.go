// Consumidor isolado do módulo instalado. Nenhuma operação utiliza a rede.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	mailcore "github.com/mailcore-hq/mailcore-go"
)

const (
	moduloSDK      = "github.com/mailcore-hq/mailcore-go"
	versaoSDK      = "v0.1.0-alpha.1"
	chaveSintetica = "mc_instalacao_sintetica"
	chaveIntencao  = "instalacao-local"
	identificador  = "id/ação+?@example.test"
	inteiroGrande  = "9007199254740993123456789"
)

var (
	bytesEsperados = []byte{'p', 'a', 'r', 't', 'e', 0, 255, '\r', '\n'}
	tipoContexto   = reflect.TypeFor[context.Context]()
	tipoConsulta   = reflect.TypeFor[mailcore.Consulta]()
	tipoOpcoes     = reflect.TypeFor[[]mailcore.OpcoesRequisicao]()
	tipoDocumentos = reflect.TypeFor[[]mailcore.Documento]()
	tipoConsumidor = reflect.TypeFor[func([]byte) error]()
	tipoResposta   = reflect.TypeFor[*mailcore.Resposta]()
	tipoErro       = reflect.TypeFor[error]()
)

type falhaInstalacao struct{ status int }

func exigir(condicao bool) {
	if !condicao {
		panic(falhaInstalacao{})
	}
}

func conferirResposta(retorno *mailcore.Resposta, erro error) *mailcore.Resposta {
	if erro != nil {
		status := 0
		var api *mailcore.Erro
		if errors.As(erro, &api) {
			status = api.Status
		}
		panic(falhaInstalacao{status: status})
	}
	exigir(retorno != nil)
	return retorno
}

func documentoResposta() mailcore.Documento {
	return mailcore.Documento{
		"objeto": mailcore.Documento{}, "lista": []any{}, "nulo": nil,
		"ativo": false, "inteiro": json.Number(inteiroGrande),
	}
}

func decodificar(corpo []byte) map[string]any {
	var documento map[string]any
	leitor := json.NewDecoder(bytes.NewReader(corpo))
	leitor.UseNumber()
	exigir(leitor.Decode(&documento) == nil && documento != nil)
	exigir(leitor.Decode(new(any)) == io.EOF)
	return documento
}

func conferirDocumento(documento map[string]any) {
	objeto, objetoValido := documento["objeto"].(map[string]any)
	lista, listaValida := documento["lista"].([]any)
	nulo, nuloPresente := documento["nulo"]
	ativo, ativoValido := documento["ativo"].(bool)
	inteiro, inteiroValido := documento["inteiro"].(json.Number)
	exigir(objetoValido && objeto != nil && len(objeto) == 0)
	exigir(listaValida && lista != nil && len(lista) == 0)
	exigir(nuloPresente && nulo == nil && ativoValido && !ativo)
	exigir(inteiroValido && string(inteiro) == inteiroGrande)
}

type transporteFalso struct {
	modo     string
	chamadas int
	consulta url.Values
	caminho  string
	corpo    []byte
}

// Somente este RoundTripper recebe os pedidos: não há encaminhamento à rede.
func (t *transporteFalso) RoundTrip(pedido *http.Request) (*http.Response, error) {
	t.chamadas++
	exigir(pedido.Context().Err() == nil)
	exigir(pedido.URL.Scheme == "https" && pedido.URL.Host == "api-mailcore.example.test")
	exigir(strings.HasPrefix(pedido.URL.EscapedPath(), "/v1/"))
	exigir(pedido.Header.Get("Authorization") == "Bearer "+chaveSintetica)
	exigir(pedido.Header.Get("Idempotency-Key") == chaveIntencao)
	exigir(pedido.Header.Get("Cookie") == "" && pedido.Header.Get("Accept-Encoding") == "identity")
	exigir(pedido.GetBody == nil)
	switch pedido.Method {
	case "GET", "POST", "PATCH", "PUT", "DELETE":
	default:
		panic(falhaInstalacao{})
	}
	t.consulta = pedido.URL.Query()
	t.caminho = pedido.URL.EscapedPath()
	t.corpo = nil
	if pedido.Body != nil {
		defer pedido.Body.Close()
		conteudo, erro := io.ReadAll(io.LimitReader(pedido.Body, 1024*1024+1))
		exigir(erro == nil && len(conteudo) <= 1024*1024)
		t.corpo = conteudo
	}
	if len(t.corpo) > 0 {
		exigir(pedido.Header.Get("Content-Type") == "application/json" && json.Valid(t.corpo))
		exigir(t.corpo[0] == '{' || t.corpo[0] == '[')
	}
	status := http.StatusOK
	tipo := pedido.Header.Get("Accept")
	var conteudo []byte
	switch t.modo {
	case "tipado":
		documento := decodificar(t.corpo)
		exigir(documento["remetente"] == "equipe@example.test" && documento["assunto"] == "Instalação Go")
		destinatarios, ok := documento["destinatarios"].([]any)
		exigir(ok && reflect.DeepEqual(destinatarios, []any{"destino@example.test"}))
		exigir(documento["html"] == "<p>Dados sintéticos</p>")
		status = http.StatusAccepted
		conteudo = []byte(`{"id":"sintetico","estado":"aceita"}`)
	case "documento":
		conferirDocumento(decodificar(t.corpo))
		var erro error
		conteudo, erro = json.Marshal(documentoResposta())
		exigir(erro == nil)
	case "erro":
		status = http.StatusServiceUnavailable
		var erro error
		conteudo, erro = json.Marshal(mailcore.Documento{"erro": mailcore.Documento{"codigo": "indisponivel", "mensagem": chaveSintetica}})
		exigir(erro == nil)
	default:
		switch tipo {
		case "application/octet-stream":
			conteudo = bytes.Clone(bytesEsperados)
		case "text/html":
			conteudo = []byte("<p>prévia sintética</p>")
		case "application/json":
			conteudo = []byte(`{}`)
		default:
			panic(falhaInstalacao{})
		}
	}
	return &http.Response{
		StatusCode: status, ContentLength: int64(len(conteudo)),
		Header: http.Header{
			"Content-Type": {tipo}, "Content-Length": {fmt.Sprint(len(conteudo))},
			"Set-Cookie": {"sintetico=privado"}, "Authorization": {"Bearer " + chaveSintetica},
		},
		Body: io.NopCloser(bytes.NewReader(conteudo)), Request: pedido,
	}, nil
}

func conferirMetadados(resposta *mailcore.Resposta) {
	exigir(resposta.Cabecalhos.Get("Set-Cookie") == "")
	exigir(resposta.Cabecalhos.Get("Authorization") == "")
	conferirRepresentacoes(resposta)
}

func conferirRepresentacoes(valor any) {
	for _, texto := range []string{fmt.Sprint(valor), fmt.Sprintf("%+v", valor), fmt.Sprintf("%#v", valor)} {
		exigir(!strings.Contains(texto, chaveSintetica))
	}
}

func argumento(tipo reflect.Type, recurso, metodo string, callbacks *int) reflect.Value {
	switch tipo {
	case tipoContexto:
		return reflect.ValueOf(context.Background())
	case tipoConsulta:
		return reflect.ValueOf(mailcore.Consulta{"limite": 2, "busca": "ação +?", "opcional": nil, "ativo": false})
	case tipoOpcoes:
		return reflect.ValueOf([]mailcore.OpcoesRequisicao{{ChaveIdempotencia: chaveIntencao}})
	case tipoDocumentos:
		return reflect.ValueOf([]mailcore.Documento{{"id": "topico", "assinatura": "opt_in"}})
	case tipoConsumidor:
		return reflect.ValueOf(func(parte []byte) error {
			exigir(bytes.Equal(parte, bytesEsperados))
			*callbacks = *callbacks + 1
			return nil
		})
	}
	switch tipo.Kind() {
	case reflect.String:
		return reflect.ValueOf(identificador).Convert(tipo)
	case reflect.Int:
		return reflect.Zero(tipo)
	case reflect.Slice:
		if tipo.Elem().Kind() == reflect.String {
			valor := reflect.MakeSlice(tipo, 1, 1)
			valor.Index(0).SetString("etiqueta")
			return valor
		}
	case reflect.Interface:
		if tipo.NumMethod() == 0 {
			var corpo any = mailcore.Documento{"texto": "sintético", "opcional": nil, "ativo": false}
			if recurso == "Lotes" && metodo == "Enviar" {
				corpo = []mailcore.Documento{{"remetente": "equipe@example.test", "destinatarios": []string{"destino@example.test"}, "assunto": "Instalação", "texto": "sintético"}}
			}
			return reflect.ValueOf(corpo)
		}
	}
	panic(falhaInstalacao{})
}

func refletir(cliente *mailcore.Cliente, falso *transporteFalso) (recursos, metodos, callbacks int) {
	objetos := []reflect.Value{}
	campos := reflect.ValueOf(cliente).Elem()
	for indice := 0; indice < campos.NumField(); indice++ {
		if campos.Type().Field(indice).PkgPath == "" {
			objetos = append(objetos, campos.Field(indice))
		}
	}
	vistos := map[uintptr]bool{}
	for len(objetos) > 0 {
		recurso := objetos[0]
		objetos = objetos[1:]
		exigir(recurso.Kind() == reflect.Pointer && !recurso.IsNil() && recurso.Elem().Kind() == reflect.Struct)
		exigir(recurso.Elem().Type().PkgPath() == moduloSDK && !vistos[recurso.Pointer()])
		vistos[recurso.Pointer()] = true
		recursos++
		camposRecurso := recurso.Elem()
		for indice := 0; indice < camposRecurso.NumField(); indice++ {
			if camposRecurso.Type().Field(indice).PkgPath == "" {
				objetos = append(objetos, camposRecurso.Field(indice))
			}
		}
		for indice := 0; indice < recurso.NumMethod(); indice++ {
			metodo := recurso.Method(indice)
			tipo := metodo.Type()
			if tipo.NumIn() == 0 || tipo.In(0) != tipoContexto {
				continue
			}
			exigir(tipo.IsVariadic() && tipo.In(tipo.NumIn()-1) == tipoOpcoes)
			exigir(tipo.NumOut() == 2 && tipo.Out(0) == tipoResposta && tipo.Out(1) == tipoErro)
			nome := recurso.Type().Method(indice).Name
			argumentos := make([]reflect.Value, tipo.NumIn())
			consultaPresente, identificadorPresente, fluxo := false, false, false
			for argumentoIndice := range argumentos {
				tipoArgumento := tipo.In(argumentoIndice)
				argumentos[argumentoIndice] = argumento(tipoArgumento, recurso.Elem().Type().Name(), nome, &callbacks)
				consultaPresente = consultaPresente || tipoArgumento == tipoConsulta
				identificadorPresente = identificadorPresente || tipoArgumento.Kind() == reflect.String
				fluxo = fluxo || tipoArgumento == tipoConsumidor
			}
			antes := falso.chamadas
			retornos := metodo.CallSlice(argumentos)
			exigir(retornos[1].IsNil() && !retornos[0].IsNil() && falso.chamadas == antes+1)
			retorno := retornos[0].Interface().(*mailcore.Resposta)
			exigir(retorno.Status == http.StatusOK)
			conferirMetadados(retorno)
			if consultaPresente {
				exigir(falso.consulta.Get("limite") == "2" && falso.consulta.Get("ativo") == "false")
				exigir(falso.consulta.Get("busca") == "ação +?" && !falso.consulta.Has("opcional"))
			}
			if identificadorPresente {
				// Identificadores hostis permanecem dados do caminho, sem redirecionar o pedido.
				esperado := strings.ReplaceAll(url.QueryEscape(identificador), "+", "%20")
				exigir(strings.Contains(falso.caminho, esperado))
			}
			if fluxo {
				exigir(retorno.TotalBytes == int64(len(bytesEsperados)) && len(retorno.Dados) == 0)
			}
			metodos++
		}
	}
	exigir(recursos == 30 && metodos == 151 && callbacks == 3 && falso.chamadas == 151)
	return
}

func conferirModulo() {
	info, ok := debug.ReadBuildInfo()
	exigir(ok)
	encontrado := false
	for _, dependencia := range info.Deps {
		if dependencia.Path == moduloSDK {
			exigir(dependencia.Version == versaoSDK && dependencia.Replace == nil)
			encontrado = true
		}
	}
	exigir(encontrado)
}

func executar() (resultado map[string]any) {
	etapa := "modulo_instalado"
	defer func() {
		if falha := recover(); falha != nil {
			status := 0
			if conhecida, ok := falha.(falhaInstalacao); ok {
				status = conhecida.status
			}
			resultado = map[string]any{"resultado": "falhou", "etapa": etapa, "status": status, "go": runtime.Version()}
		}
	}()
	conferirModulo()
	etapa = "construir_cliente"
	falso := &transporteFalso{}
	cliente, erro := mailcore.NovoCliente(chaveSintetica, mailcore.Configuracao{
		URLBase: "https://api-mailcore.example.test/v1", ClienteHTTP: &http.Client{Transport: falso},
	})
	exigir(erro == nil && cliente != nil)
	defer cliente.Fechar()
	conferirRepresentacoes(cliente)
	etapa = "operacoes_instaladas"
	recursos, metodos, callbacks := refletir(cliente, falso)
	ctx := context.Background()
	opcoes := mailcore.OpcoesRequisicao{ChaveIdempotencia: chaveIntencao}
	etapa = "pedido_email_tipado"
	falso.modo = "tipado"
	retorno := conferirResposta(cliente.Emails.Enviar(ctx, mailcore.PedidoEmail{
		Remetente: "equipe@example.test", Destinatarios: []string{"destino@example.test"},
		Assunto: "Instalação Go", HTML: "<p>Dados sintéticos</p>",
	}, opcoes))
	var envio mailcore.RespostaEnvio
	exigir(retorno.Status == http.StatusAccepted && retorno.Decodificar(&envio) == nil)
	exigir(envio.ID == "sintetico" && envio.Estado == "aceita")
	conferirMetadados(retorno)
	etapa = "json_preservado"
	falso.modo = "documento"
	retorno = conferirResposta(cliente.Eventos.Criar(ctx, documentoResposta(), opcoes))
	var documento map[string]any
	exigir(retorno.Decodificar(&documento) == nil)
	conferirDocumento(documento)
	etapa = "erro_protegido"
	falso.modo = "erro"
	_, erro = cliente.Emails.Listar(ctx, nil, opcoes)
	var api *mailcore.Erro
	exigir(errors.As(erro, &api) && api.Status == http.StatusServiceUnavailable && api.Codigo == "indisponivel")
	conferirRepresentacoes(api)
	exigir(api.Cabecalhos.Get("Set-Cookie") == "" && api.Cabecalhos.Get("Authorization") == "")
	etapa = "contextos_encerrados"
	antes := falso.chamadas
	cancelado, cancelar := context.WithCancel(ctx)
	cancelar()
	_, erro = cliente.Emails.Listar(cancelado, nil, opcoes)
	exigir(errors.Is(erro, context.Canceled))
	conferirRepresentacoes(erro)
	expirado, encerrar := context.WithDeadline(ctx, time.Now().Add(-time.Second))
	defer encerrar()
	_, erro = cliente.Emails.Listar(expirado, nil, opcoes)
	exigir(errors.Is(erro, context.DeadlineExceeded) && falso.chamadas == antes)
	conferirRepresentacoes(erro)
	return map[string]any{
		"resultado": "aprovado", "go": runtime.Version(), "modulo": moduloSDK, "versao_modulo": versaoSDK,
		"recursos": recursos, "metodos": metodos, "callbacks_fluxo": callbacks,
		"chamadas_operacoes": metodos, "chamadas_provas_adicionais": falso.chamadas - metodos,
		"pedido_email_tipado": true, "numeros_json_preservados": true, "objetos_listas_null_false": true,
		"contextos_sem_chamada_http": true, "representacoes_protegidas": true,
		"http_falso": true, "rede_utilizada": false, "publicado": false,
	}
}

func main() {
	resultado := executar()
	if json.NewEncoder(os.Stdout).Encode(resultado) != nil {
		os.Exit(1)
	}
	if resultado["resultado"] != "aprovado" {
		os.Exit(1)
	}
}
