package mailcore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func (c *Cliente) executar(ctx context.Context, metodo, caminho string, corpo any, consulta Consulta, formato string, aceitarVazio bool, opcoes OpcoesRequisicao, consumidor func([]byte) error) (*Resposta, error) {
	if c == nil || c.http == nil || c.fechado.Load() {
		return nil, erroConfiguracao("Construa um cliente e mantenha-o aberto durante as operações.")
	}
	if ctx == nil {
		return nil, erroConfiguracao("Forneça um contexto não nulo.")
	}
	if err := ctx.Err(); err != nil {
		return nil, falhaTransporte(ctx, err)
	}
	if !strings.HasPrefix(caminho, "/") || strings.HasPrefix(caminho, "//") || strings.ContainsAny(caminho, "?#\\\r\n\t ") {
		return nil, erroConfiguracao("Caminho de recurso inválido.")
	}
	for _, parte := range strings.Split(caminho, "/") {
		if parte == "." || parte == ".." {
			return nil, erroConfiguracao("Caminho de recurso inválido.")
		}
	}
	switch metodo {
	case "GET", "POST", "PATCH", "PUT", "DELETE":
	default:
		return nil, erroConfiguracao("Método HTTP inválido.")
	}
	switch formato {
	case "json", "texto", "bytes", "fluxo":
	default:
		return nil, erroConfiguracao("Formato de resposta inválido.")
	}
	if formato == "fluxo" && consumidor == nil {
		return nil, erroConfiguracao("Forneça um consumidor para o download em fluxo.")
	}
	if opcoes.TempoLimite < 0 || opcoes.TempoLimite > 24*time.Hour || opcoes.LimiteRespostaBytes < 0 {
		return nil, erroConfiguracao("Prazos e limites devem ser positivos.")
	}
	tempo, limite := c.tempo, c.limite
	if opcoes.TempoLimite != 0 {
		tempo = opcoes.TempoLimite
	}
	if opcoes.LimiteRespostaBytes != 0 {
		limite = opcoes.LimiteRespostaBytes
	}
	ctx, cancelar := context.WithTimeout(ctx, tempo)
	defer cancelar()
	filtros, err := codificarConsulta(consulta)
	if err != nil {
		return nil, err
	}
	endereco := c.base + caminho
	if filtros != "" {
		endereco += "?" + filtros
	}
	cabecalhos, err := cabecalhosPedido(opcoes.Cabecalhos, opcoes.ChaveIdempotencia)
	if err != nil {
		return nil, err
	}
	var dados []byte
	if corpo != nil {
		dados, err = codificarCorpo(corpo)
		if err != nil {
			return nil, err
		}
		cabecalhos.Set("Content-Type", "application/json")
	}
	chave, err := c.resolverChave(ctx)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, falhaTransporte(ctx, err)
	}
	cabecalhos.Set("Authorization", "Bearer "+chave)
	cabecalhos.Set("User-Agent", "mailcore-go/"+Versao)
	cabecalhos.Set("Accept-Encoding", "identity")
	aceito := "application/json"
	if formato == "texto" {
		aceito = "text/html"
	} else if formato == "bytes" || formato == "fluxo" {
		aceito = "application/octet-stream"
	}
	cabecalhos.Set("Accept", aceito)
	// Não expomos GetBody para o Transport repetir uma mutação. Mesmo mutações
	// sem corpo usam um reader opaco: Body=nil tornaria Idempotency-Key replayable.
	var leitor io.Reader
	if dados != nil || metodo != "GET" {
		leitor = io.NopCloser(bytes.NewReader(dados))
	}
	pedido, err := http.NewRequestWithContext(ctx, metodo, endereco, leitor)
	if err != nil {
		return nil, erroConfiguracao("Não foi possível construir o pedido HTTP.")
	}
	pedido.GetBody = nil
	if len(dados) > 0 {
		pedido.ContentLength = int64(len(dados))
	}
	pedido.Header = cabecalhos
	retorno, err := c.http.Do(pedido)
	if err != nil {
		if retorno != nil && retorno.Body != nil {
			retorno.Body.Close()
		}
		return nil, falhaTransporte(ctx, err)
	}
	if retorno == nil || retorno.Body == nil {
		return nil, &Erro{Codigo: "resposta_invalida", Mensagem: "A API não forneceu uma resposta HTTP válida."}
	}
	defer retorno.Body.Close()
	cabecalhosResposta := limparCabecalhos(retorno.Header)
	if retorno.StatusCode < 100 || retorno.StatusCode > 599 {
		return nil, &Erro{Codigo: "resposta_invalida", Mensagem: "A API retornou status HTTP inválido."}
	}
	if retorno.StatusCode >= 300 && retorno.StatusCode < 400 {
		return nil, &Erro{Codigo: "redirecionamento_recusado", Mensagem: "A API retornou um redirecionamento recusado.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
	}
	tamanhoCabecalhos := 0
	for nome, valores := range retorno.Header {
		for _, valor := range valores {
			tamanhoCabecalhos += len(nome) + len(valor) + 4
		}
	}
	if tamanhoCabecalhos > 64*1024 {
		return nil, &Erro{Codigo: "resposta_excedida", Mensagem: "Cabeçalhos excederam o limite.", Status: retorno.StatusCode}
	}
	for nome, valores := range retorno.Header {
		if !strings.EqualFold(nome, "Content-Encoding") {
			continue
		}
		for _, valor := range valores {
			for _, token := range strings.Split(valor, ",") {
				codificacao := strings.ToLower(strings.TrimSpace(token))
				if codificacao != "" && codificacao != "identity" {
					return nil, &Erro{Codigo: "codificacao_inesperada", Mensagem: "A resposta usa codificação não negociada.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
				}
			}
		}
	}
	anunciado := retorno.ContentLength
	tamanhoDefinido := false
	for nome, valores := range retorno.Header {
		if !strings.EqualFold(nome, "Content-Length") {
			continue
		}
		for _, valor := range valores {
			tamanho, erroTamanho := strconv.ParseInt(valor, 10, 64)
			if erroTamanho != nil || tamanho < 0 || strings.Trim(valor, "0123456789") != "" || (tamanhoDefinido && anunciado != tamanho) {
				return nil, &Erro{Codigo: "resposta_invalida", Mensagem: "Tamanho HTTP inválido.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
			}
			anunciado, tamanhoDefinido = tamanho, true
		}
	}
	if anunciado > limite {
		return nil, &Erro{Codigo: "resposta_excedida", Mensagem: "A resposta ultrapassa o limite de bytes.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
	}
	sucesso := retorno.StatusCode >= 200 && retorno.StatusCode <= 299
	fluxo := sucesso && formato == "fluxo"
	var acumulado bytes.Buffer
	buffer := make([]byte, 32*1024)
	var total int64
	semProgresso := 0
	for {
		if err := ctx.Err(); err != nil {
			return nil, falhaTransporte(ctx, err)
		}
		lidos, erroLeitura := retorno.Body.Read(buffer)
		if lidos < 0 || lidos > len(buffer) {
			return nil, &Erro{Codigo: "resposta_invalida", Mensagem: "O provedor HTTP retornou uma leitura inválida.", Status: retorno.StatusCode}
		}
		if int64(lidos) > limite-total {
			return nil, &Erro{Codigo: "resposta_excedida", Mensagem: "A resposta ultrapassa o limite de bytes.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
		}
		total += int64(lidos)
		if lidos > 0 {
			semProgresso = 0
			if fluxo {
				if err := consumidor(buffer[:lidos]); err != nil {
					return nil, err
				}
			} else {
				acumulado.Write(buffer[:lidos])
			}
		} else if erroLeitura == nil {
			semProgresso++
			if semProgresso > 100 {
				return nil, &Erro{Codigo: "fluxo_sem_progresso", Mensagem: "O provedor HTTP não avançou na leitura.", Status: retorno.StatusCode}
			}
		}
		if erroLeitura != nil {
			if erroLeitura != io.EOF {
				return nil, falhaLeitura(ctx, erroLeitura, retorno.StatusCode, cabecalhosResposta)
			}
			break
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, falhaTransporte(ctx, err)
	}
	if anunciado >= 0 && anunciado != total {
		return nil, &Erro{Codigo: "resposta_incompleta", Mensagem: "A resposta terminou com tamanho incompleto.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
	}
	dadosResposta := acumulado.Bytes()
	if !sucesso {
		falha := &Erro{Codigo: "erro_http", Mensagem: "A API retornou um erro HTTP.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
		var envelope struct {
			Erro *struct {
				Codigo   string `json:"codigo"`
				Mensagem string `json:"mensagem"`
			} `json:"erro"`
		}
		if utf8.Valid(dadosResposta) && json.Unmarshal(dadosResposta, &envelope) == nil && envelope.Erro != nil && envelope.Erro.Codigo != "" {
			falha.Codigo, falha.Mensagem = envelope.Erro.Codigo, envelope.Erro.Mensagem
		}
		return nil, falha
	}
	if formato == "json" && !(len(dadosResposta) == 0 && (retorno.StatusCode == 204 || aceitarVazio)) {
		if !utf8.Valid(dadosResposta) || !json.Valid(dadosResposta) {
			return nil, &Erro{Codigo: "resposta_invalida", Mensagem: "A API retornou JSON inválido.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
		}
	}
	if formato == "texto" && !utf8.Valid(dadosResposta) {
		return nil, &Erro{Codigo: "resposta_invalida", Mensagem: "A API retornou texto UTF-8 inválido.", Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta}
	}
	return &Resposta{Dados: dadosResposta, Status: retorno.StatusCode, Cabecalhos: cabecalhosResposta, TotalBytes: total}, nil
}

func falhaTransporte(ctx context.Context, erro error) *Erro {
	causa := ctx.Err()
	if causa == nil {
		if errors.Is(erro, context.Canceled) {
			causa = context.Canceled
		}
		if errors.Is(erro, context.DeadlineExceeded) {
			causa = context.DeadlineExceeded
		}
	}
	codigo := "falha_de_rede"
	if causa == context.Canceled {
		codigo = "cancelado"
	} else if causa == context.DeadlineExceeded {
		codigo = "prazo_excedido"
	}
	return &Erro{Codigo: codigo, Mensagem: "Não foi possível concluir a comunicação com a API.", causa: causa}
}

func falhaLeitura(ctx context.Context, erro error, status int, cabecalhos http.Header) *Erro {
	if ctx.Err() != nil || errors.Is(erro, context.Canceled) || errors.Is(erro, context.DeadlineExceeded) {
		return falhaTransporte(ctx, erro)
	}
	return &Erro{Codigo: "resposta_incompleta", Mensagem: "Não foi possível concluir a leitura da resposta.", Status: status, Cabecalhos: cabecalhos}
}

func limparCabecalhos(origem http.Header) http.Header {
	resultado := origem.Clone()
	for nome := range resultado {
		switch strings.ToLower(nome) {
		case "set-cookie", "authorization", "proxy-authorization":
			delete(resultado, nome)
		}
	}
	return resultado
}

func cabecalhosPedido(origem http.Header, idempotencia string) (http.Header, error) {
	resultado := http.Header{}
	for nome, valores := range origem {
		if nome == "" {
			return nil, erroConfiguracao("Cabeçalho sem nome.")
		}
		for _, c := range nome {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || strings.ContainsRune("!#$%&'*+-.^_`|~", c)) {
				return nil, erroConfiguracao("Nome de cabeçalho inválido.")
			}
		}
		switch strings.ToLower(nome) {
		case "authorization", "proxy-authorization", "cookie", "host", "content-length", "content-type", "accept", "connection", "transfer-encoding", "idempotency-key", "user-agent", "accept-encoding", "trailer", "te", "upgrade":
			return nil, erroConfiguracao("Cabeçalho reservado pelo SDK.")
		}
		for _, valor := range valores {
			for _, c := range valor {
				if c < 32 || c > 126 {
					return nil, erroConfiguracao("Valor de cabeçalho inválido.")
				}
			}
			resultado.Add(nome, valor)
		}
	}
	if idempotencia != "" {
		if len(idempotencia) > 256 || !chaveValida(idempotencia) {
			return nil, erroConfiguracao("Use uma chave de idempotência com até 256 caracteres ASCII visíveis.")
		}
		resultado.Set("Idempotency-Key", idempotencia)
	}
	return resultado, nil
}
