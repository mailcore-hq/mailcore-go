package mailcore

import (
	"encoding/json"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

const limiteCorpoJSON = 50 * 1024 * 1024

type corpoObjeto struct{ dados any }

type corpoLista struct{ dados any }

func lista(dados any) (any, error) {
	if err := validarLista(dados); err != nil {
		return nil, err
	}
	return corpoLista{dados: dados}, nil
}

func objeto(dados any) (any, error) {
	if dados == nil {
		dados = Documento{}
	}
	return corpoObjeto{dados: dados}, nil
}

func validarLista(dados any) error {
	v := reflect.ValueOf(dados)
	for v.IsValid() && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			return erroConfiguracao("Forneça uma lista JSON explícita.")
		}
		v = v.Elem()
	}
	if !v.IsValid() || (v.Kind() != reflect.Slice && v.Kind() != reflect.Array) || (v.Kind() == reflect.Slice && v.IsNil()) {
		return erroConfiguracao("Forneça uma lista JSON explícita.")
	}
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 && v.Type() != reflect.TypeFor[json.RawMessage]() {
		return erroConfiguracao("O lote exige uma lista JSON, sem codificação Base64 implícita.")
	}
	return nil
}

func codificarCorpo(corpo any) ([]byte, error) {
	objetoRaiz, listaRaiz := false, false
	if objeto, ok := corpo.(corpoObjeto); ok {
		corpo, objetoRaiz = objeto.dados, true
	}
	if lista, ok := corpo.(corpoLista); ok {
		corpo, listaRaiz = lista.dados, true
	}
	if err := validarJSON(reflect.ValueOf(corpo), 0); err != nil {
		return nil, err
	}
	dados, err := json.Marshal(corpo)
	if err != nil || !utf8.Valid(dados) {
		return nil, erroConfiguracao("Use somente valores JSON e texto UTF-8 válidos.")
	}
	if objetoRaiz {
		if string(dados) == "null" {
			dados = []byte("{}")
		}
		if len(dados) == 0 || dados[0] != '{' {
			return nil, erroConfiguracao("O pedido exige um objeto JSON.")
		}
	}
	if listaRaiz && (len(dados) == 0 || dados[0] != '[') {
		return nil, erroConfiguracao("O lote exige uma lista JSON.")
	}
	if len(dados) > limiteCorpoJSON {
		return nil, erroConfiguracao("O corpo JSON ultrapassa 50 MiB.")
	}
	return dados, nil
}

// A biblioteca padrão substitui UTF-8 inválido em strings. A verificação evita
// perder informação antes de serializar; ciclos são recusados pelo teto de profundidade.
func validarJSON(v reflect.Value, profundidade int) error {
	if profundidade > 100 {
		return erroConfiguracao("O corpo excede a profundidade JSON permitida.")
	}
	if !v.IsValid() {
		return nil
	}
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
		if !v.IsNil() {
			return validarJSON(v.Elem(), profundidade+1)
		}
	case reflect.String:
		if !utf8.ValidString(v.String()) {
			return erroConfiguracao("Use somente texto UTF-8 válido.")
		}
	case reflect.Float32, reflect.Float64:
		if math.IsNaN(v.Float()) || math.IsInf(v.Float(), 0) {
			return erroConfiguracao("Números JSON devem ser finitos.")
		}
	case reflect.Map:
		it := v.MapRange()
		for it.Next() {
			if err := validarJSON(it.Key(), profundidade+1); err != nil {
				return err
			}
			if err := validarJSON(it.Value(), profundidade+1); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		// []byte é codificado como Base64; RawMessage precisa de UTF-8 próprio.
		if v.Type() == reflect.TypeFor[json.RawMessage]() && !utf8.Valid(v.Bytes()) {
			return erroConfiguracao("JSON bruto precisa usar UTF-8 válido.")
		}
		if v.Type().Elem().Kind() != reflect.Uint8 {
			for i := 0; i < v.Len(); i++ {
				if err := validarJSON(v.Index(i), profundidade+1); err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).PkgPath == "" && t.Field(i).Tag.Get("json") != "-" {
				if err := validarJSON(v.Field(i), profundidade+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func segmento(valor string) (string, error) {
	if valor == "" || valor == "." || valor == ".." || !utf8.ValidString(valor) {
		return "", erroConfiguracao("Informe um identificador UTF-8 válido.")
	}
	for _, c := range valor {
		if c < 32 || c == 127 {
			return "", erroConfiguracao("Identificadores não podem conter controles.")
		}
	}
	// PathEscape mantém ':' e '@' válidos no segmento. QueryEscape torna todos
	// os delimitadores opacos; espaço é %20, e não '+'.
	return strings.ReplaceAll(url.QueryEscape(valor), "+", "%20"), nil
}

func codificarConsulta(consulta Consulta) (string, error) {
	valores := url.Values{}
	for nome, valor := range consulta {
		if !utf8.ValidString(nome) || nome == "" {
			return "", erroConfiguracao("Use nomes UTF-8 válidos nos filtros.")
		}
		if valor == nil {
			continue
		}
		var texto string
		v := reflect.ValueOf(valor)
		switch v.Kind() {
		case reflect.String:
			texto = v.String()
		case reflect.Bool:
			texto = strconv.FormatBool(v.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			texto = strconv.FormatInt(v.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			texto = strconv.FormatUint(v.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			if math.IsNaN(v.Float()) || math.IsInf(v.Float(), 0) {
				return "", erroConfiguracao("Filtros numéricos devem ser finitos.")
			}
			texto = strconv.FormatFloat(v.Float(), 'g', -1, v.Type().Bits())
		default:
			return "", erroConfiguracao("Os filtros aceitam somente valores simples.")
		}
		if !utf8.ValidString(texto) {
			return "", erroConfiguracao("Use texto UTF-8 válido nos filtros.")
		}
		valores.Set(nome, texto)
	}
	return strings.ReplaceAll(valores.Encode(), "+", "%20"), nil
}

func opcao(opcoes []OpcoesRequisicao) (OpcoesRequisicao, error) {
	if len(opcoes) > 1 {
		return OpcoesRequisicao{}, erroConfiguracao("Forneça no máximo um conjunto de opções.")
	}
	if len(opcoes) == 0 {
		return OpcoesRequisicao{}, nil
	}
	return opcoes[0], nil
}
