package mailcore

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

var corpoAssinaturaTeste = []byte(`{"type":"email.received","created_at":"2026-11-14T00:00:00Z","data":{"email_id":"mensagem-sintetica","to":["destino@example.test"],"ativo":false,"nulo":null,"objeto":{},"lista":[]}}`)

func vetorAssinaturaTeste(corpo []byte, instante time.Time) (http.Header, string) {
	return vetorAssinaturaComChaveTeste(corpo, instante, []byte("0123456789abcdef0123456789abcdef"))
}

func vetorAssinaturaComChaveTeste(corpo []byte, instante time.Time, chave []byte) (http.Header, string) {
	id := "msg_sintetico_independente"
	tempo := strconv.FormatInt(instante.Unix(), 10)
	mac := hmac.New(sha256.New, chave)
	mac.Write([]byte(id + "." + tempo + "."))
	mac.Write(corpo)
	cabecalhos := http.Header{"Svix-Id": {id}, "Svix-Timestamp": {tempo}, "Svix-Signature": {"v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))}}
	return cabecalhos, "whsec_" + base64.StdEncoding.EncodeToString(chave)
}

func exigirFalhaAssinaturaTeste(t *testing.T, corpo []byte, cabecalhos http.Header, segredo string, opcoes ...OpcoesAssinatura) {
	t.Helper()
	evento, err := VerificarAssinatura(corpo, cabecalhos, segredo, opcoes...)
	var falha *Erro
	if evento != nil || !errors.As(err, &falha) || falha.Codigo != "assinatura_invalida" {
		t.Fatalf("Evento inválido não foi recusado com erro de assinatura: %v", err)
	}
	for _, texto := range []string{fmt.Sprintf("%v", err), fmt.Sprintf("%+v", err), fmt.Sprintf("%#v", err)} {
		if (segredo != "" && strings.Contains(texto, segredo)) || strings.Contains(texto, "destino@example.test") || strings.Contains(texto, "mensagem-sintetica") {
			t.Fatal("Erro de assinatura expôs segredo ou conteúdo.")
		}
	}
}

func TestAssinaturaAutenticaBytesEConservaDados(t *testing.T) {
	agora := time.Unix(1794614400, 0)
	cabecalhos, segredo := vetorAssinaturaTeste(corpoAssinaturaTeste, agora)
	evento, err := VerificarAssinatura(corpoAssinaturaTeste, cabecalhos, segredo, OpcoesAssinatura{Agora: agora})
	if err != nil {
		t.Fatal(err)
	}
	if evento.Tipo != "email.received" || evento.CriadoEm != "2026-11-14T00:00:00Z" {
		t.Fatal("Envelope do evento alterado.")
	}
	var dados map[string]any
	if err := json.Unmarshal(evento.Dados, &dados); err != nil {
		t.Fatal(err)
	}
	if dados["ativo"] != false || dados["nulo"] != nil {
		t.Fatal("Valores JSON alterados.")
	}
	if _, ok := dados["objeto"].(map[string]any); !ok {
		t.Fatal("Objeto vazio alterado.")
	}
	if lista, ok := dados["lista"].([]any); !ok || len(lista) != 0 {
		t.Fatal("Lista vazia alterada.")
	}
	for _, texto := range []string{fmt.Sprintf("%v", evento), fmt.Sprintf("%+v", evento), fmt.Sprintf("%#v", evento)} {
		if strings.Contains(texto, "destino@example.test") {
			t.Fatal("Representação textual expôs o payload.")
		}
	}
	alterado := append(bytes.Clone(corpoAssinaturaTeste), ' ')
	exigirFalhaAssinaturaTeste(t, alterado, cabecalhos, segredo, OpcoesAssinatura{Agora: agora})
}

func TestAssinaturaJanelaPassadoEFuturo(t *testing.T) {
	agora := time.Unix(1794614400, 0)
	for _, diferenca := range []int{-301, -300, -299, 0, 299, 300, 301} {
		t.Run(strconv.Itoa(diferenca), func(t *testing.T) {
			cabecalhos, segredo := vetorAssinaturaTeste(corpoAssinaturaTeste, agora.Add(time.Duration(diferenca)*time.Second))
			_, err := VerificarAssinatura(corpoAssinaturaTeste, cabecalhos, segredo, OpcoesAssinatura{Agora: agora})
			if (err == nil) != (diferenca >= -300 && diferenca <= 300) {
				t.Fatal("Janela temporal incompatível.")
			}
		})
	}
	cabecalhos, segredo := vetorAssinaturaTeste(corpoAssinaturaTeste, agora.Add(31*time.Second))
	exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, cabecalhos, segredo, OpcoesAssinatura{Agora: agora, Tolerancia: 30 * time.Second})
	for _, tolerancia := range []time.Duration{-time.Second, 301 * time.Second} {
		exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, cabecalhos, segredo, OpcoesAssinatura{Agora: agora, Tolerancia: tolerancia})
	}
	exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, cabecalhos, segredo, OpcoesAssinatura{Agora: agora}, OpcoesAssinatura{Agora: agora})
}

func TestAssinaturaRejeitaCabecalhosAusentesDuplicadosEAmbiguos(t *testing.T) {
	agora := time.Unix(1794614400, 0)
	original, segredo := vetorAssinaturaTeste(corpoAssinaturaTeste, agora)
	for _, nome := range []string{"Svix-Id", "Svix-Timestamp", "Svix-Signature"} {
		t.Run(nome, func(t *testing.T) {
			sem := original.Clone()
			sem.Del(nome)
			exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, sem, segredo, OpcoesAssinatura{Agora: agora})
			duplicado := original.Clone()
			duplicado.Add(nome, original.Get(nome))
			exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, duplicado, segredo, OpcoesAssinatura{Agora: agora})
			ambiguo := original.Clone()
			ambiguo[strings.ToLower(nome)] = []string{original.Get(nome)}
			exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, ambiguo, segredo, OpcoesAssinatura{Agora: agora})
		})
	}
	normalizados := http.Header{"HTTP_SVIX_ID": original.Values("Svix-Id"), "HTTP_SVIX_TIMESTAMP": original.Values("Svix-Timestamp"), "HTTP_SVIX_SIGNATURE": original.Values("Svix-Signature")}
	if _, err := VerificarAssinatura(corpoAssinaturaTeste, normalizados, segredo, OpcoesAssinatura{Agora: agora}); err != nil {
		t.Fatal("Cabeçalhos de ambiente não foram normalizados.")
	}
	for _, tempo := range []string{"0", "-1", "001794614400", "1794614400.0", "1e10", " 1794614400", "17946144000000000", "nao-numerico"} {
		alterado := original.Clone()
		alterado.Set("Svix-Timestamp", tempo)
		exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, alterado, segredo, OpcoesAssinatura{Agora: agora})
	}
	for _, id := range []string{"", strings.Repeat("m", 257), "msg\r\ninjecao"} {
		alterado := original.Clone()
		alterado.Set("Svix-Id", id)
		exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, alterado, segredo, OpcoesAssinatura{Agora: agora})
	}
}

func TestAssinaturaVersaoBase64ESegredoEstritos(t *testing.T) {
	agora := time.Unix(1794614400, 0)
	original, segredo := vetorAssinaturaTeste(corpoAssinaturaTeste, agora)
	assinatura := strings.TrimPrefix(original.Get("Svix-Signature"), "v1,")
	for _, valor := range []string{"v2," + assinatura, assinatura, "v1," + strings.TrimRight(assinatura, "="), "v1," + assinatura + "\n", "v1,%%%", "v1," + strings.Repeat("a", 4097)} {
		alterado := original.Clone()
		alterado.Set("Svix-Signature", valor)
		exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, alterado, segredo, OpcoesAssinatura{Agora: agora})
	}
	for _, valor := range []string{"", strings.TrimPrefix(segredo, "whsec_"), "whsec_%%%", segredo + "\n", strings.TrimRight(segredo, "="), "whsec_" + base64.StdEncoding.EncodeToString([]byte("curto")), "whsec_" + strings.Repeat("a", 1025)} {
		exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, original, valor, OpcoesAssinatura{Agora: agora})
	}
	rotacao := original.Clone()
	rotacao.Set("Svix-Signature", "v2,invalida v1,invalida "+original.Get("Svix-Signature"))
	if _, err := VerificarAssinatura(corpoAssinaturaTeste, rotacao, segredo, OpcoesAssinatura{Agora: agora}); err != nil {
		t.Fatal("Candidato válido de rotação foi recusado.")
	}
	foraDoLimite := original.Clone()
	foraDoLimite.Set("Svix-Signature", strings.Repeat("v1,invalida ", 32)+original.Get("Svix-Signature"))
	exigirFalhaAssinaturaTeste(t, corpoAssinaturaTeste, foraDoLimite, segredo, OpcoesAssinatura{Agora: agora})
	for _, tamanho := range []int{15, 16} {
		cabecalhos, chave := vetorAssinaturaComChaveTeste(corpoAssinaturaTeste, agora, bytes.Repeat([]byte("a"), tamanho))
		_, err := VerificarAssinatura(corpoAssinaturaTeste, cabecalhos, chave, OpcoesAssinatura{Agora: agora})
		if (err == nil) != (tamanho == 16) {
			t.Fatal("Comprimento mínimo do segredo não foi aplicado.")
		}
	}
}

func TestAssinaturaAutenticacaoNaoSubstituiValidacaoDoEvento(t *testing.T) {
	agora := time.Unix(1794614400, 0)
	for _, corpo := range [][]byte{[]byte(`[]`), []byte(`null`), []byte(`{}`), []byte(`{"type":"email.received","created_at":"data"}`), []byte(`{"type":"","created_at":"data","data":{}}`), []byte(`{"type":3,"created_at":"data","data":{}}`), []byte(`{"type":"email.received","created_at":"data","data":{"texto":"` + string([]byte{255}) + `"}}`), []byte(`{} {}`)} {
		cabecalhos, segredo := vetorAssinaturaTeste(corpo, agora)
		exigirFalhaAssinaturaTeste(t, corpo, cabecalhos, segredo, OpcoesAssinatura{Agora: agora})
	}
	corpo := []byte(`{"type":"evento","created_at":"data","data":null}`)
	cabecalhos, segredo := vetorAssinaturaTeste(corpo, agora)
	evento, err := VerificarAssinatura(corpo, cabecalhos, segredo, OpcoesAssinatura{Agora: agora})
	if err != nil || string(evento.Dados) != "null" {
		t.Fatal("Campo data presente como null foi perdido.")
	}
}

func TestAssinaturaLimiteDeOitoMiB(t *testing.T) {
	agora := time.Unix(1794614400, 0)
	prefixo := []byte(`{"type":"evento","created_at":"data","data":{"texto":"`)
	sufixo := []byte(`"}}`)
	corpo := append(prefixo, bytes.Repeat([]byte("a"), 8*1024*1024-len(prefixo)-len(sufixo))...)
	corpo = append(corpo, sufixo...)
	cabecalhos, segredo := vetorAssinaturaTeste(corpo, agora)
	if _, err := VerificarAssinatura(corpo, cabecalhos, segredo, OpcoesAssinatura{Agora: agora}); err != nil {
		t.Fatal("Corpo de8MiB foi recusado.")
	}
	corpo = append(corpo, ' ')
	cabecalhos, segredo = vetorAssinaturaTeste(corpo, agora)
	exigirFalhaAssinaturaTeste(t, corpo, cabecalhos, segredo, OpcoesAssinatura{Agora: agora})
}
