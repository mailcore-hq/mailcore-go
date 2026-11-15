package mailcore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Versao identifica esta distribuição local, ainda sem publicação do módulo.
const Versao = "0.1.0-alpha.1"

// Documento mantém os campos do contrato JSON, inclusive null e false.
// Um Documento{} é {}; uma lista vazia deve ser construída com []T{}.
type Documento map[string]any

// Consulta representa filtros e cursores; valores nil são omitidos da URL.
type Consulta map[string]any

// Configuracao é copiada na construção. Não a altere durante NovoCliente.
type Configuracao struct {
	URLBase             string
	TempoLimite         time.Duration
	TempoConexao        time.Duration
	LimiteRespostaBytes int64
	CertificadoCA       string
	PermitirHTTPLocal   bool
	ClienteHTTP         *http.Client
	FornecedorChave     func(context.Context) (string, error)
}

// OpcoesRequisicao permite no máximo um conjunto de opções por chamada.
// Zero em prazos/limites utiliza o padrão do cliente; negativos são recusados.
type OpcoesRequisicao struct {
	TempoLimite         time.Duration
	LimiteRespostaBytes int64
	ChaveIdempotencia   string
	Cabecalhos          http.Header
}

// Resposta mantém os bytes originais. Dados contém JSON, texto ou bytes conforme
// o método; nas variantes Fluxo, somente TotalBytes é preenchido com o tamanho.
// Após o retorno, os dados pertencem ao chamador.
type Resposta struct {
	Dados      json.RawMessage
	Status     int
	Cabecalhos http.Header
	TotalBytes int64
}

// Decodificar preserva números em json.Number ao utilizar destinos dinâmicos.
// Para objetos tipados, os campos seguem as tags JSON do destino.
func (r *Resposta) Decodificar(destino any) error {
	if r == nil || len(r.Dados) == 0 {
		return &Erro{Codigo: "resposta_invalida", Mensagem: "Não há JSON para decodificar."}
	}
	decodificador := json.NewDecoder(bytes.NewReader(r.Dados))
	decodificador.UseNumber()
	if err := decodificador.Decode(destino); err != nil {
		return &Erro{Codigo: "resposta_invalida", Mensagem: "O destino não é compatível com o JSON recebido.", Status: r.Status}
	}
	if decodificador.Decode(new(any)) != io.EOF {
		return &Erro{Codigo: "resposta_invalida", Mensagem: "A resposta contém JSON adicional.", Status: r.Status}
	}
	return nil
}

func (r *Resposta) String() string {
	if r == nil {
		return "Mailcore.Resposta(nil)"
	}
	return fmt.Sprintf("Mailcore.Resposta(status=%d, bytes=%d)", r.Status, r.TotalBytes)
}
func (r *Resposta) GoString() string { return r.String() }

// Erro não inclui pedidos, URLs ou dados privados em sua representação textual.
// Mensagem pode conter o erro de negócio da API e deve ser tratado como privado.
type Erro struct {
	Codigo     string
	Mensagem   string
	Status     int
	Cabecalhos http.Header
	causa      error
}

func (e *Erro) Error() string {
	if e == nil {
		return "Mailcore.Erro(nil)"
	}
	return fmt.Sprintf("Mailcore: operação recusada (HTTP %d).", e.Status)
}
func (e *Erro) GoString() string { return e.Error() }
func (e *Erro) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.causa
}

func erroConfiguracao(mensagem string) *Erro {
	return &Erro{Codigo: "configuracao_invalida", Mensagem: mensagem}
}

// AnexoEmail utiliza Base64 padrão no corpo JSON. IDConteudo permite imagens CID.
type AnexoEmail struct {
	Nome           string  `json:"nome"`
	TipoConteudo   string  `json:"tipo_conteudo"`
	ConteudoBase64 string  `json:"conteudo_base64"`
	IDConteudo     *string `json:"id_conteudo,omitempty"`
}

// ModeloEmail referencia uma versão publicada e suas variáveis.
type ModeloEmail struct {
	ID        string    `json:"id"`
	Variaveis Documento `json:"variaveis,omitempty"`
}

// PedidoEmail oferece campos tipados para o envio. A API valida as combinações
// de conteúdo direto/modelo e as capacidades específicas de envio e lote.
// Use Documento para distinguir omissão, null e valores vazios explicitamente.
type PedidoEmail struct {
	Remetente     string            `json:"remetente,omitempty"`
	Destinatarios []string          `json:"destinatarios"`
	Assunto       string            `json:"assunto,omitempty"`
	Texto         string            `json:"texto,omitempty"`
	HTML          string            `json:"html,omitempty"`
	CC            []string          `json:"cc,omitempty"`
	CCO           []string          `json:"cco,omitempty"`
	ResponderPara *string           `json:"responder_para,omitempty"`
	AgendadaEm    *string           `json:"agendada_em,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	Cabecalhos    map[string]string `json:"cabecalhos,omitempty"`
	Anexos        []AnexoEmail      `json:"anexos,omitempty"`
	Modelo        *ModeloEmail      `json:"modelo,omitempty"`
}

// RespostaEnvio indica aceitação; a entrega é acompanhada por eventos/webhooks.
type RespostaEnvio struct {
	ID     string `json:"id"`
	Estado string `json:"estado"`
}
