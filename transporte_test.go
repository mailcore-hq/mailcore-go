package mailcore

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func servidorTransporte(t *testing.T, atender http.HandlerFunc) *httptest.Server {
	t.Helper()
	s := httptest.NewUnstartedServer(atender)
	s.Config.ErrorLog = log.New(io.Discard, "", 0)
	s.StartTLS()
	t.Cleanup(s.Close)
	return s
}

func configuracaoTransporte(t *testing.T, s *httptest.Server) Configuracao {
	t.Helper()
	certificado := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: s.Certificate().Raw})
	caminho := t.TempDir() + "/ca-sintetica.crt"
	if err := os.WriteFile(caminho, certificado, 0600); err != nil {
		t.Fatal(err)
	}
	return Configuracao{URLBase: s.URL + "/v1", CertificadoCA: caminho}
}

func clienteTransporte(t *testing.T, config Configuracao) *Cliente {
	t.Helper()
	c, err := NovoCliente("mc_chave_sintetica_testes", config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c.Fechar)
	return c
}

func aceitarTransporte(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	io.WriteString(w, `{"id":"email-sintetico","estado":"na_fila"}`)
}

type transporteFuncao func(*http.Request) (*http.Response, error)

func (f transporteFuncao) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestTransporteTLSJSONEIdempotencia(t *testing.T) {
	var recebido Documento
	var autorizacao, idempotencia, codificacao string
	var protecao sync.Mutex
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		protecao.Lock()
		defer protecao.Unlock()
		autorizacao, idempotencia, codificacao = r.Header.Get("Authorization"), r.Header.Get("Idempotency-Key"), r.Header.Get("Accept-Encoding")
		if r.Method != "POST" || r.URL.Path != "/v1/emails" {
			t.Errorf("Rota inesperada: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&recebido); err != nil {
			t.Error(err)
		}
		aceitarTransporte(w)
	})
	c := clienteTransporte(t, configuracaoTransporte(t, s))
	resposta, err := c.Emails.Enviar(context.Background(), Documento{"remetente": "equipe@example.test", "destinatarios": []string{"destino@example.test"}, "esquema": Documento{}, "lista": []any{}, "nulo": nil, "ativo": false}, OpcoesRequisicao{ChaveIdempotencia: "pedido-estavel"})
	if err != nil {
		t.Fatal(err)
	}
	protecao.Lock()
	defer protecao.Unlock()
	if resposta.Status != 202 || autorizacao != "Bearer mc_chave_sintetica_testes" || idempotencia != "pedido-estavel" || codificacao != "identity" {
		t.Fatal("O contrato HTTP não foi preservado.")
	}
	if _, ok := recebido["esquema"].(map[string]any); !ok {
		t.Fatal("Objeto vazio alterado.")
	}
	if lista, ok := recebido["lista"].([]any); !ok || len(lista) != 0 {
		t.Fatal("Lista vazia alterada.")
	}
	if recebido["nulo"] != nil || recebido["ativo"] != false {
		t.Fatal("null ou false alterado.")
	}
	var dados RespostaEnvio
	if err := resposta.Decodificar(&dados); err != nil || dados.ID != "email-sintetico" {
		t.Fatal("DTO de aceitação incompatível.")
	}
}

func TestTransporteTLSExigeAutoridadeEIdentidade(t *testing.T) {
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, `{}`) })
	semCA := clienteTransporte(t, Configuracao{URLBase: s.URL + "/v1"})
	if _, err := semCA.Emails.Obter(context.Background(), "id"); err == nil {
		t.Fatal("Certificado não confiável foi aceito.")
	}
	config := configuracaoTransporte(t, s)
	config.URLBase = strings.Replace(config.URLBase, "127.0.0.1", "localhost", 1)
	if _, err := clienteTransporte(t, config).Emails.Obter(context.Background(), "id"); err == nil {
		t.Fatal("Hostname incorreto foi aceito.")
	}
	config.URLBase = s.URL + "/v1"
	if _, err := clienteTransporte(t, config).Emails.Obter(context.Background(), "id"); err != nil {
		t.Fatal(err)
	}
}

func TestTransporteNaoSegueRedirecionamentoNemArmazenaCookies(t *testing.T) {
	var destino, cookies atomic.Int64
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Cookie") != "" {
			cookies.Add(1)
		}
		w.Header().Add("Set-Cookie", "sessao=sintetica; Path=/")
		if strings.HasSuffix(r.URL.Path, "/redirecionar") {
			w.Header().Set("Location", "/destino")
			w.WriteHeader(307)
			return
		}
		if r.URL.Path == "/destino" {
			destino.Add(1)
		}
		io.WriteString(w, `{}`)
	})
	c := clienteTransporte(t, configuracaoTransporte(t, s))
	for i := 0; i < 2; i++ {
		resposta, err := c.Emails.Obter(context.Background(), "normal")
		if err != nil {
			t.Fatal(err)
		}
		if resposta.Cabecalhos.Get("Set-Cookie") != "" {
			t.Fatal("Cookie privado exposto nos metadados.")
		}
	}
	_, err := c.Emails.Obter(context.Background(), "redirecionar")
	var falha *Erro
	if !errors.As(err, &falha) || falha.Status != 307 || destino.Load() != 0 || cookies.Load() != 0 {
		t.Fatal("Política de cookies/redirecionamento não foi preservada.")
	}
	if c.proprio.Proxy != nil || c.http.Jar != nil || c.proprio.TLSClientConfig.InsecureSkipVerify || c.proprio.TLSClientConfig.MinVersion < tls.VersionTLS12 {
		t.Fatal("Política padrão de rede insegura.")
	}
}

func TestTransporteInjetadoNaoAlteraClienteERecebeCorpoNaoRepetivel(t *testing.T) {
	var observados atomic.Int64
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Cookie") != "" {
			t.Error("Cookie da aplicação foi enviado.")
		}
		aceitarTransporte(w)
	})
	original := s.Client()
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(s.URL)
	jar.SetCookies(u, []*http.Cookie{{Name: "sessao", Value: "sintetica"}})
	original.Jar = jar
	original.CheckRedirect = nil
	base := original.Transport
	original.Transport = transporteFuncao(func(r *http.Request) (*http.Response, error) {
		if r.GetBody != nil || r.Body == nil || r.Body == http.NoBody {
			t.Error("A mutação foi construída como repetível.")
		}
		observados.Add(1)
		return base.RoundTrip(r)
	})
	c := clienteTransporte(t, Configuracao{URLBase: s.URL + "/v1", ClienteHTTP: original})
	if _, err := c.Emails.Enviar(context.Background(), Documento{"texto": "Olá"}, OpcoesRequisicao{ChaveIdempotencia: "pedido"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Emails.Cancelar(context.Background(), "id", OpcoesRequisicao{ChaveIdempotencia: "cancelamento"}); err != nil {
		t.Fatal(err)
	}
	if observados.Load() != 2 || original.Jar != jar || original.CheckRedirect != nil {
		t.Fatal("Cliente injetado foi alterado ou chamadas inesperadas ocorreram.")
	}
}

func TestTransportePoolGETERecuperacaoDeConexao(t *testing.T) {
	for _, mutacao := range []string{"enviar", "cancelar"} {
		t.Run(mutacao, func(t *testing.T) {
			var conexoes, pedidos atomic.Int64
			var aquecida, mutada atomic.Value
			s := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					aquecida.Store(r.RemoteAddr)
					io.WriteString(w, `{}`)
					return
				}
				pedidos.Add(1)
				mutada.Store(r.RemoteAddr)
				io.Copy(io.Discard, r.Body)
				conexao, _, err := w.(http.Hijacker).Hijack()
				if err != nil {
					t.Error(err)
					return
				}
				conexao.Close()
			}))
			s.Config.ErrorLog = log.New(io.Discard, "", 0)
			s.Config.ConnState = func(_ net.Conn, estado http.ConnState) {
				if estado == http.StateNew {
					conexoes.Add(1)
				}
			}
			s.StartTLS()
			t.Cleanup(s.Close)
			c := clienteTransporte(t, configuracaoTransporte(t, s))
			for i := 0; i < 2; i++ {
				if _, err := c.Emails.Obter(context.Background(), "aquecimento"); err != nil {
					t.Fatal(err)
				}
			}
			if conexoes.Load() != 1 {
				t.Fatal("Consultas sequenciais não reutilizaram a conexão.")
			}
			var err error
			if mutacao == "enviar" {
				_, err = c.Emails.Enviar(context.Background(), Documento{"texto": "sintetico"}, OpcoesRequisicao{ChaveIdempotencia: "intencao"})
			} else {
				_, err = c.Emails.Cancelar(context.Background(), "id", OpcoesRequisicao{ChaveIdempotencia: "intencao"})
			}
			if err == nil || pedidos.Load() != 1 {
				t.Fatalf("Mutação ambígua repetida: pedidos=%d erro=%v", pedidos.Load(), err)
			}
			if aquecida.Load() != mutada.Load() {
				t.Fatal("A mutação não utilizou a conexão aquecida deste ensaio.")
			}
			if _, err := c.Emails.Obter(context.Background(), "recuperacao"); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTransporteContextoCancelaEExpira(t *testing.T) {
	for _, modo := range []string{"cancelar", "prazo"} {
		t.Run(modo, func(t *testing.T) {
			entrou := make(chan struct{})
			s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) { close(entrou); <-r.Context().Done() })
			c := clienteTransporte(t, configuracaoTransporte(t, s))
			ctx, cancelar := context.WithCancel(context.Background())
			esperado := context.Canceled
			if modo == "prazo" {
				cancelar()
				ctx, cancelar = context.WithTimeout(context.Background(), 500*time.Millisecond)
				esperado = context.DeadlineExceeded
			}
			defer cancelar()
			terminou := make(chan error, 1)
			go func() { _, err := c.Emails.Obter(ctx, "aguardando"); terminou <- err }()
			select {
			case <-entrou:
			case <-time.After(2 * time.Second):
				t.Fatal("O servidor não recebeu a operação.")
			}
			if modo == "cancelar" {
				cancelar()
			}
			select {
			case err := <-terminou:
				if !errors.Is(err, esperado) {
					t.Fatalf("Identidade de contexto perdida: %v", err)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("Contexto não interrompeu a espera.")
			}
		})
	}
}

func TestTransporteConcorrenteIsolaRespostasEFornecedor(t *testing.T) {
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Operacao", strings.TrimPrefix(r.URL.Path, "/v1/emails/"))
		io.WriteString(w, `{"id":`+strconv.Quote(r.Header.Get("X-Operacao"))+`}`)
	})
	config := configuracaoTransporte(t, s)
	var resolucoes atomic.Int64
	config.FornecedorChave = func(ctx context.Context) (string, error) { resolucoes.Add(1); return "mc_chave_resolvida", nil }
	c, err := NovoCliente("", config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c.Fechar)
	var grupo sync.WaitGroup
	falhas := make(chan error, 24)
	for i := 0; i < 24; i++ {
		grupo.Add(1)
		go func(i int) {
			defer grupo.Done()
			id := strconv.Itoa(i)
			resposta, erro := c.Emails.Obter(context.Background(), id, OpcoesRequisicao{Cabecalhos: http.Header{"X-Operacao": {id}}})
			if erro != nil {
				falhas <- erro
				return
			}
			var dados map[string]any
			if erro = resposta.Decodificar(&dados); erro != nil || dados["id"] != id || resposta.Cabecalhos.Get("X-Operacao") != id {
				falhas <- fmt.Errorf("resposta concorrente incompatível")
			}
		}(i)
	}
	grupo.Wait()
	close(falhas)
	for err := range falhas {
		t.Error(err)
	}
	if resolucoes.Load() != 24 {
		t.Fatal("Fornecedor não foi resolvido por operação.")
	}
}

func TestTransporteJSONRecusaPerdaDeInformacao(t *testing.T) {
	var chamadas atomic.Int64
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) { chamadas.Add(1); aceitarTransporte(w) })
	c := clienteTransporte(t, configuracaoTransporte(t, s))
	ciclo := Documento{}
	ciclo["ciclo"] = ciclo
	for _, pedido := range []any{Documento{"texto": string([]byte{255})}, Documento{string([]byte{255}): "valor"}, Documento{"numero": math.Inf(1)}, Documento{"numero": math.NaN()}, ciclo, json.RawMessage(`{"texto":"` + string([]byte{255}) + `"}`)} {
		if _, err := c.Emails.Enviar(context.Background(), pedido); err == nil {
			t.Fatal("JSON inválido foi aceito.")
		}
	}
	if chamadas.Load() != 0 {
		t.Fatal("Pedido inválido chegou à rede.")
	}
}

func TestTransporteJSONRawENumerosPrecisos(t *testing.T) {
	original := []byte(" {\"objeto\":{},\"lista\":[],\"nulo\":null,\"ativo\":false,\"grande\":123456789012345678901234567890} \n")
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) { w.Write(original) })
	resposta, err := clienteTransporte(t, configuracaoTransporte(t, s)).Emails.Obter(context.Background(), "id")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, resposta.Dados) {
		t.Fatal("Bytes JSON originais foram alterados.")
	}
	var dados map[string]any
	if err := resposta.Decodificar(&dados); err != nil {
		t.Fatal(err)
	}
	if dados["grande"] != json.Number("123456789012345678901234567890") || dados["nulo"] != nil || dados["ativo"] != false {
		t.Fatal("Precisão ou valores JSON alterados.")
	}
	for _, corpo := range [][]byte{[]byte(`{} {}`), []byte(`{"texto":"` + string([]byte{255}) + `"}`), []byte(`{"incompleto":`)} {
		t.Run(strconv.Itoa(len(corpo)), func(t *testing.T) {
			s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) { w.Write(corpo) })
			if _, err := clienteTransporte(t, configuracaoTransporte(t, s)).Emails.Obter(context.Background(), "id"); err == nil {
				t.Fatal("Resposta JSON inválida foi aceita.")
			}
		})
	}
}

func TestTransporteLimitesPadraoDeCorpoEResposta(t *testing.T) {
	var pedidos atomic.Int64
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		pedidos.Add(1)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"texto":"`)
		bloco := strings.Repeat("a", 32*1024)
		for i := 0; i < 17*1024/32; i++ {
			if _, err := io.WriteString(w, bloco); err != nil {
				return
			}
		}
		io.WriteString(w, `"}`)
	})
	c := clienteTransporte(t, configuracaoTransporte(t, s))
	if _, err := c.Emails.Enviar(context.Background(), Documento{"texto": strings.Repeat("a", 50*1024*1024)}); err == nil {
		t.Fatal("JSON serializado acima de50MiB foi aceito.")
	}
	if pedidos.Load() != 0 {
		t.Fatal("Corpo excessivo chegou à rede.")
	}
	_, err := c.Emails.Obter(context.Background(), "id")
	var falha *Erro
	if !errors.As(err, &falha) || falha.Codigo != "resposta_excedida" {
		t.Fatalf("Resposta acima de16MiB aceita: %v", err)
	}
}

func TestTransporteFalhasHTTPNaoRepetemNemExpoemCorpo(t *testing.T) {
	for _, status := range []int{422, 429, 503} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			var chamadas atomic.Int64
			s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
				chamadas.Add(1)
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(status)
				io.WriteString(w, `{"erro":{"codigo":"regra","mensagem":"dado_privado_sintetico"}}`)
			})
			_, err := clienteTransporte(t, configuracaoTransporte(t, s)).Emails.Enviar(context.Background(), Documento{"texto": "sintetico"}, OpcoesRequisicao{ChaveIdempotencia: "pedido"})
			var falha *Erro
			if !errors.As(err, &falha) || falha.Status != status || chamadas.Load() != 1 {
				t.Fatal("Status ou tentativas incompatíveis.")
			}
			for _, texto := range []string{fmt.Sprintf("%v", err), fmt.Sprintf("%+v", err), fmt.Sprintf("%#v", err)} {
				if strings.Contains(texto, "dado_privado_sintetico") {
					t.Fatal("Erro textual expôs dados privados.")
				}
			}
		})
	}
}

func TestTransporteFluxoOitoMiBSemAcumularEComBytesExatos(t *testing.T) {
	const tamanho = 8 * 1024 * 1024
	bloco := bytes.Repeat([]byte{0, 255, 13, 10}, 4096)
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "message/rfc822")
		w.Header().Set("Content-Length", strconv.Itoa(tamanho))
		for i := 0; i < tamanho/len(bloco); i++ {
			if _, err := w.Write(bloco); err != nil {
				return
			}
		}
	})
	c := clienteTransporte(t, configuracaoTransporte(t, s))
	runtime.GC()
	var antes, depois runtime.MemStats
	runtime.ReadMemStats(&antes)
	var total int64
	resposta, err := c.Recebimento.BaixarOriginalFluxo(context.Background(), "id", func(parte []byte) error {
		for i, b := range parte {
			if b != bloco[int(total+int64(i))%len(bloco)] {
				return errors.New("bytes alterados")
			}
		}
		total += int64(len(parte))
		return nil
	})
	runtime.ReadMemStats(&depois)
	if err != nil {
		t.Fatal(err)
	}
	if total != tamanho || resposta.TotalBytes != tamanho || len(resposta.Dados) != 0 {
		t.Fatal("Fluxo não preservou bytes/contagem.")
	}
	if alocado := depois.TotalAlloc - antes.TotalAlloc; alocado > 16*1024*1024 {
		t.Fatalf("Fluxo alocou %dbytes para conteúdo de8MiB.", alocado)
	}
}

func TestTransporteFluxoPropagaErroDoConsumidorECancelamento(t *testing.T) {
	for _, modo := range []string{"consumidor", "contexto"} {
		t.Run(modo, func(t *testing.T) {
			s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
				bloco := bytes.Repeat([]byte("a"), 32*1024)
				for i := 0; i < 256; i++ {
					if _, err := w.Write(bloco); err != nil {
						return
					}
					w.(http.Flusher).Flush()
				}
			})
			c := clienteTransporte(t, configuracaoTransporte(t, s))
			ctx, cancelar := context.WithCancel(context.Background())
			defer cancelar()
			sentinela := errors.New("consumidor interrompeu")
			_, err := c.Recebimento.BaixarOriginalFluxo(ctx, "id", func(parte []byte) error {
				if modo == "consumidor" {
					return sentinela
				}
				cancelar()
				return nil
			})
			esperado := sentinela
			if modo == "contexto" {
				esperado = context.Canceled
			}
			if !errors.Is(err, esperado) {
				t.Fatalf("Erro original perdido: %v", err)
			}
		})
	}
}

func TestTransporteTruncamentoETetoEmFluxo(t *testing.T) {
	for _, modo := range []string{"truncado", "chunked", "codificacao"} {
		t.Run(modo, func(t *testing.T) {
			s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
				if modo == "truncado" {
					w.Header().Set("Content-Length", "100")
					io.WriteString(w, "curto")
					return
				}
				if modo == "codificacao" {
					w.Header().Set("Content-Encoding", "gzip")
					io.WriteString(w, "nao-negociado")
					return
				}
				for i := 0; i < 4; i++ {
					io.WriteString(w, strings.Repeat("a", 1024))
					w.(http.Flusher).Flush()
				}
			})
			var consumido int
			_, err := clienteTransporte(t, configuracaoTransporte(t, s)).Recebimento.BaixarOriginalFluxo(context.Background(), "id", func(parte []byte) error { consumido += len(parte); return nil }, OpcoesRequisicao{LimiteRespostaBytes: 2048})
			if err == nil || consumido > 2048 {
				t.Fatal("Truncamento/codificação/limite não foi recusado.")
			}
		})
	}
}

func TestTransporteFechaCorpoMesmoQuandoConsumidorFalha(t *testing.T) {
	var fechamentos atomic.Int64
	corpo := &corpoContadoTransporte{Reader: strings.NewReader("conteudo"), fechamentos: &fechamentos}
	fornecedor := &http.Client{Transport: transporteFuncao(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: corpo, ContentLength: -1}, nil
	})}
	c := clienteTransporte(t, Configuracao{URLBase: "https://api.example.test/v1", ClienteHTTP: fornecedor})
	sentinela := errors.New("parar")
	_, err := c.Recebimento.BaixarOriginalFluxo(context.Background(), "id", func([]byte) error { return sentinela })
	if !errors.Is(err, sentinela) || fechamentos.Load() != 1 {
		t.Fatal("Erro do consumidor ou fechamento perdido.")
	}
}

type corpoContadoTransporte struct {
	io.Reader
	fechamentos *atomic.Int64
}

func (c *corpoContadoTransporte) Close() error { c.fechamentos.Add(1); return nil }

func TestTransporteConfiguracaoCabecalhosEFornecedorInvalidoNaoFazemRede(t *testing.T) {
	var chamadas atomic.Int64
	httpLocal := &http.Client{Transport: transporteFuncao(func(*http.Request) (*http.Response, error) {
		chamadas.Add(1)
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{}`)), ContentLength: -1}, nil
	})}
	for _, base := range []string{"http://api.example.test/v1", "https://usuario:segredo@example.test/v1", "https://api.example.test/v1?segredo=1", "https://api.example.test/v1#parte", "https://api.example.test/\r\n"} {
		if _, err := NovoCliente("chave", Configuracao{URLBase: base, PermitirHTTPLocal: true, ClienteHTTP: httpLocal}); err == nil {
			t.Fatal("URL insegura aceita.")
		}
	}
	c := clienteTransporte(t, Configuracao{URLBase: "https://api.example.test/v1", ClienteHTTP: httpLocal})
	for _, cabecalho := range []http.Header{{"Authorization": {"roubo"}}, {"Host": {"outro"}}, {"Cookie": {"sessao"}}, {"X-Teste": {"valor\r\ninjecao"}}, {"X\nNome": {"valor"}}} {
		if _, err := c.Emails.Obter(context.Background(), "id", OpcoesRequisicao{Cabecalhos: cabecalho}); err == nil {
			t.Fatal("Cabeçalho inseguro aceito.")
		}
	}
	t.Setenv("MAILCORE_API_KEY", "mc_chave_ambiente")
	semFallback, err := NovoCliente("", Configuracao{URLBase: "https://api.example.test/v1", ClienteHTTP: httpLocal, FornecedorChave: func(context.Context) (string, error) { return "", nil }})
	if err != nil {
		t.Fatal(err)
	}
	defer semFallback.Fechar()
	if _, err := semFallback.Emails.Obter(context.Background(), "id"); err == nil {
		t.Fatal("Fornecedor inválido recorreu ao ambiente.")
	}
	if chamadas.Load() != 0 {
		t.Fatal("Configuração inválida chegou à rede.")
	}
}

func mimeComprimidoTransporte(t *testing.T) []byte {
	t.Helper()
	var destino bytes.Buffer
	comprimir := gzip.NewWriter(&destino)
	if _, err := comprimir.Write([]byte("From: origem@example.test\r\n\r\nMensagem sintética.")); err != nil {
		t.Fatal(err)
	}
	if err := comprimir.Close(); err != nil {
		t.Fatal(err)
	}
	return destino.Bytes()
}

func TestTransporteRecusaCodificacaoEmDuasLinhasTLS(t *testing.T) {
	comprimido := mimeComprimidoTransporte(t)
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Encoding", "identity")
		w.Header().Add("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "message/rfc822")
		w.Write(comprimido)
	})
	var consumido int
	_, err := clienteTransporte(t, configuracaoTransporte(t, s)).Recebimento.BaixarOriginalFluxo(context.Background(), "id", func(parte []byte) error {
		consumido += len(parte)
		return nil
	})
	var falha *Erro
	if !errors.As(err, &falha) || falha.Codigo != "codificacao_inesperada" || consumido != 0 {
		t.Fatal("Duas linhas de codificação permitiram conteúdo não negociado ao consumidor.")
	}
}

func TestTransporteRecusaCodificacaoEmListaTLS(t *testing.T) {
	comprimido := mimeComprimidoTransporte(t)
	s := servidorTransporte(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "identity, gzip")
		w.Header().Set("Content-Type", "message/rfc822")
		w.Write(comprimido)
	})
	var consumido int
	_, err := clienteTransporte(t, configuracaoTransporte(t, s)).Recebimento.BaixarOriginalFluxo(context.Background(), "id", func(parte []byte) error {
		consumido += len(parte)
		return nil
	})
	var falha *Erro
	if !errors.As(err, &falha) || falha.Codigo != "codificacao_inesperada" || consumido != 0 {
		t.Fatal("Lista de codificações permitiu conteúdo não negociado ao consumidor.")
	}
}

func TestTransporteRecusaMetadadosAmbiguosDoProvedor(t *testing.T) {
	comprimido := mimeComprimidoTransporte(t)
	for _, caso := range []struct {
		nome       string
		cabecalhos http.Header
		codigo     string
	}{
		{"codificacao_minuscula", http.Header{"content-encoding": {"identity", "gzip"}}, "codificacao_inesperada"},
		{"tamanhos_distintos", http.Header{"content-length": {"3", "8"}}, "resposta_invalida"},
		{"tamanhos_em_caixas_distintas", http.Header{"Content-Length": {"3"}, "content-length": {"8"}}, "resposta_invalida"},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			var fechamentos atomic.Int64
			corpo := &corpoContadoTransporte{Reader: bytes.NewReader(comprimido), fechamentos: &fechamentos}
			fornecedor := &http.Client{Transport: transporteFuncao(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Header: caso.cabecalhos, Body: corpo, ContentLength: -1}, nil
			})}
			c := clienteTransporte(t, Configuracao{URLBase: "https://api.example.test/v1", ClienteHTTP: fornecedor})
			var consumido int
			_, err := c.Recebimento.BaixarOriginalFluxo(context.Background(), "id", func(parte []byte) error {
				consumido += len(parte)
				return nil
			})
			var falha *Erro
			if !errors.As(err, &falha) || falha.Codigo != caso.codigo || consumido != 0 || fechamentos.Load() != 1 {
				t.Fatal("Metadados inválidos não impediram o consumo ou o fechamento do corpo.")
			}
		})
	}
}
