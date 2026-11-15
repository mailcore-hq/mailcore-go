package mailcore

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// OpcoesAssinatura permite injetar um relógio no teste e reduzir a janela temporal.
type OpcoesAssinatura struct {
	Agora      time.Time
	Tolerancia time.Duration
}

// EventoWebhook mantém os dados JSON originais após a autenticação do corpo bruto.
type EventoWebhook struct {
	Tipo     string          `json:"type"`
	CriadoEm string          `json:"created_at"`
	Dados    json.RawMessage `json:"data"`
}

func (e *EventoWebhook) String() string   { return "Mailcore.EventoWebhook(conteúdo protegido)" }
func (e *EventoWebhook) GoString() string { return e.String() }

// VerificarAssinatura autentica bytes brutos, HMAC-SHA256 e janela de até 300s.
// A deduplicação durável de eventos continua pertencendo à aplicação.
func VerificarAssinatura(corpo []byte, cabecalhos http.Header, segredo string, opcoes ...OpcoesAssinatura) (*EventoWebhook, error) {
	falha := func() (*EventoWebhook, error) {
		return nil, &Erro{Codigo: "assinatura_invalida", Mensagem: "Não foi possível autenticar o webhook."}
	}
	if len(opcoes) > 1 || len(corpo) > 8*1024*1024 || !strings.HasPrefix(segredo, "whsec_") || len(segredo) > 1024 {
		return falha()
	}
	config := OpcoesAssinatura{}
	if len(opcoes) == 1 {
		config = opcoes[0]
	}
	if config.Agora.IsZero() {
		config.Agora = time.Now()
	}
	if config.Tolerancia == 0 {
		config.Tolerancia = 300 * time.Second
	}
	if config.Tolerancia < 0 || config.Tolerancia > 300*time.Second {
		return falha()
	}
	normalizados := map[string]string{}
	for nome, valores := range cabecalhos {
		nome = strings.ToLower(strings.ReplaceAll(nome, "_", "-"))
		nome = strings.TrimPrefix(nome, "http-")
		if nome != "svix-id" && nome != "svix-timestamp" && nome != "svix-signature" {
			continue
		}
		if _, repetido := normalizados[nome]; repetido || len(valores) != 1 {
			return falha()
		}
		normalizados[nome] = valores[0]
	}
	id, instante, assinaturas := normalizados["svix-id"], normalizados["svix-timestamp"], normalizados["svix-signature"]
	if !chaveValida(id) || len(id) > 256 || len(instante) > 16 || len(instante) == 0 || instante[0] == '0' || len(assinaturas) > 4096 {
		return falha()
	}
	for _, c := range instante {
		if c < '0' || c > '9' {
			return falha()
		}
	}
	segundos, err := strconv.ParseInt(instante, 10, 64)
	if err != nil {
		return falha()
	}
	diferenca := config.Agora.Sub(time.Unix(segundos, 0))
	if diferenca > config.Tolerancia || diferenca < -config.Tolerancia {
		return falha()
	}
	chave, err := base64.StdEncoding.Strict().DecodeString(strings.TrimPrefix(segredo, "whsec_"))
	if err != nil || len(chave) < 16 || base64.StdEncoding.EncodeToString(chave) != strings.TrimPrefix(segredo, "whsec_") {
		return falha()
	}
	mac := hmac.New(sha256.New, chave)
	mac.Write([]byte(id + "." + instante + "."))
	mac.Write(corpo)
	esperado := mac.Sum(nil)
	valido := false
	for i, parte := range strings.Split(assinaturas, " ") {
		if i >= 32 {
			break
		}
		versao, valor, ok := strings.Cut(parte, ",")
		if !ok || versao != "v1" {
			continue
		}
		assinatura, err := base64.StdEncoding.Strict().DecodeString(valor)
		if err == nil && base64.StdEncoding.EncodeToString(assinatura) == valor {
			valido = hmac.Equal(esperado, assinatura) || valido
		}
	}
	if !valido || !utf8.Valid(corpo) {
		return falha()
	}
	var evento EventoWebhook
	if json.Unmarshal(corpo, &evento) != nil || evento.Tipo == "" || evento.CriadoEm == "" || len(evento.Dados) == 0 {
		return falha()
	}
	return &evento, nil
}
