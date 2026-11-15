package mailcore

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type executor interface {
	executar(context.Context, string, string, any, Consulta, string, bool, OpcoesRequisicao, func([]byte) error) (*Resposta, error)
}

// Cliente mantém configuração imutável e conexões compartilháveis entre goroutines.
// Construa com NovoCliente. Não altere os recursos públicos durante operações.
type Cliente struct {
	Emails               *Emails
	Lotes                *Lotes
	Dominios             *Dominios
	ChavesAPI            *ChavesAPI
	Contatos             *Contatos
	Segmentos            *Segmentos
	Topicos              *Topicos
	PropriedadesContatos *PropriedadesContatos
	Modelos              *Modelos
	Webhooks             *Webhooks
	Campanhas            *Campanhas
	Automacoes           *Automacoes
	Eventos              *Eventos
	Recebimento          *Recebimento
	Caixas               *Caixas
	CLI                  *CLIMailcore
	base                 string
	chave                string
	fornecedor           func(context.Context) (string, error)
	tempo                time.Duration
	limite               int64
	http                 *http.Client
	proprio              *http.Transport
	fechado              atomic.Bool
}

// NovoCliente aceita uma chave explícita ou, com chave vazia, MAILCORE_API_KEY
// resolvida por operação. FornecedorChave substitui ambos e deve ser concorrente.
// Não são iniciadas goroutines ou conexões de rede durante a construção.
func NovoCliente(chave string, configuracoes ...Configuracao) (*Cliente, error) {
	if len(configuracoes) > 1 {
		return nil, erroConfiguracao("Forneça no máximo uma configuração.")
	}
	var config Configuracao
	if len(configuracoes) == 1 {
		config = configuracoes[0]
	}
	if chave != "" && !chaveValida(chave) {
		return nil, erroConfiguracao("Informe uma chave Mailcore válida.")
	}
	if chave != "" && config.FornecedorChave != nil {
		return nil, erroConfiguracao("Escolha chave explícita ou fornecedor, sem ambiguidade.")
	}
	if config.URLBase == "" {
		config.URLBase = "https://mailcore.com.br/v1"
	}
	u, err := url.Parse(config.URLBase)
	if err != nil || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" || strings.ContainsAny(config.URLBase, "\\\r\n\t ") {
		return nil, erroConfiguracao("Use URL base HTTPS sem credenciais, consulta ou fragmento.")
	}
	for _, c := range config.URLBase {
		if c < 32 || c == 127 {
			return nil, erroConfiguracao("A URL base contém controles.")
		}
	}
	local := u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" || u.Hostname() == "::1"
	if u.Scheme != "https" && !(u.Scheme == "http" && local && config.PermitirHTTPLocal) {
		return nil, erroConfiguracao("HTTPS é obrigatório; HTTP explícito de teste se limita a loopback.")
	}
	if config.TempoLimite == 0 {
		config.TempoLimite = 30 * time.Second
	}
	if config.TempoConexao == 0 {
		config.TempoConexao = 10 * time.Second
	}
	if config.TempoLimite < 0 || config.TempoConexao < 0 || config.TempoLimite > 24*time.Hour || config.TempoConexao > 24*time.Hour || config.LimiteRespostaBytes < 0 {
		return nil, erroConfiguracao("Prazos devem ser positivos, até 24h; limites não podem ser negativos.")
	}
	if config.LimiteRespostaBytes == 0 {
		config.LimiteRespostaBytes = 16 * 1024 * 1024
	}
	c := &Cliente{base: strings.TrimRight(config.URLBase, "/"), chave: chave, fornecedor: config.FornecedorChave, tempo: config.TempoLimite, limite: config.LimiteRespostaBytes}
	if config.ClienteHTTP != nil {
		// A cópia não altera a instância do chamador. O Transport injetado ainda
		// controla TLS, proxy, buffering e eventuais tentativas próprias.
		copia := *config.ClienteHTTP
		copia.Jar = nil
		copia.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
		c.http = &copia
	} else {
		tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
		if config.CertificadoCA != "" {
			pem, erro := os.ReadFile(config.CertificadoCA)
			if erro != nil {
				return nil, erroConfiguracao("Não foi possível ler a autoridade certificadora.")
			}
			cas, erro := x509.SystemCertPool()
			if erro != nil {
				cas = x509.NewCertPool()
			}
			if !cas.AppendCertsFromPEM(pem) {
				return nil, erroConfiguracao("A autoridade certificadora não contém certificados válidos.")
			}
			tlsConfig.RootCAs = cas
		}
		c.proprio = &http.Transport{
			Proxy: nil, DialContext: (&net.Dialer{Timeout: config.TempoConexao, KeepAlive: 30 * time.Second}).DialContext,
			TLSClientConfig: tlsConfig, TLSHandshakeTimeout: config.TempoConexao,
			ForceAttemptHTTP2: true, DisableCompression: true, MaxIdleConns: 100,
			MaxIdleConnsPerHost: 100, IdleConnTimeout: 90 * time.Second,
			MaxResponseHeaderBytes: 64 * 1024, ExpectContinueTimeout: time.Second,
		}
		c.http = &http.Client{Transport: c.proprio, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	c.Emails = novoEmails(c)
	c.Lotes = novoLotes(c)
	c.Dominios = novoDominios(c)
	c.ChavesAPI = novoChavesAPI(c)
	c.Contatos = novoContatos(c)
	c.Segmentos = novoSegmentos(c)
	c.Topicos = novoTopicos(c)
	c.PropriedadesContatos = novoPropriedadesContatos(c)
	c.Modelos = novoModelos(c)
	c.Webhooks = novoWebhooks(c)
	c.Campanhas = novoCampanhas(c)
	c.Automacoes = novoAutomacoes(c)
	c.Eventos = novoEventos(c)
	c.Recebimento = novoRecebimento(c)
	c.Caixas = novoCaixas(c)
	c.CLI = novoCLIMailcore(c)
	return c, nil
}

// Fechar recusa novas operações e libera conexões ociosas próprias. Operações
// iniciadas antes da chamada terminam pelo seu contexto; o provedor injetado
// continua pertencendo à aplicação.
func (c *Cliente) Fechar() {
	if c != nil && c.fechado.CompareAndSwap(false, true) && c.proprio != nil {
		c.proprio.CloseIdleConnections()
	}
}
func (c *Cliente) String() string   { return "Mailcore.Cliente(credenciais protegidas)" }
func (c *Cliente) GoString() string { return c.String() }

func chaveValida(chave string) bool {
	if chave == "" {
		return false
	}
	for _, c := range chave {
		if c < 33 || c > 126 {
			return false
		}
	}
	return true
}

func (c *Cliente) resolverChave(ctx context.Context) (string, error) {
	chave := c.chave
	if c.fornecedor != nil {
		var err error
		chave, err = c.fornecedor(ctx)
		if err != nil {
			return "", falhaTransporte(ctx, err)
		}
	} else if chave == "" {
		chave = os.Getenv("MAILCORE_API_KEY")
	}
	if !chaveValida(chave) {
		return "", erroConfiguracao("Informe uma chave Mailcore válida ou configure MAILCORE_API_KEY.")
	}
	return chave, nil
}
