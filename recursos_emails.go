// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import (
	"context"
	"strconv"
)

// Emails reúne operações do contrato público Mailcore.
type Emails struct {
	transporte executor
}

func novoEmails(transporte executor) *Emails {
	return &Emails{
		transporte: transporte,
	}
}

// Enviar aceita a intenção de envio; acompanhe a entrega separadamente.
func (recurso *Emails) Enviar(ctx context.Context, dados any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(dados)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/emails", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Emails) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/emails", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Emails) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	idCodificado, err := segmento(id)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/emails/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Reagendar executa a operação correspondente do contrato público.
func (recurso *Emails) Reagendar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	idCodificado, err := segmento(id)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "PATCH", "/emails/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Cancelar executa a operação correspondente do contrato público.
func (recurso *Emails) Cancelar(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	idCodificado, err := segmento(id)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/emails/"+idCodificado+"/cancel", nil, nil, "json", false, configuracao, nil)
}

// ObterPorIdempotencia executa a operação correspondente do contrato público.
func (recurso *Emails) ObterPorIdempotencia(ctx context.Context, chave string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	chaveCodificado, err := segmento(chave)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/emails/idempotency/"+chaveCodificado, nil, nil, "json", false, configuracao, nil)
}

// Eventos executa a operação correspondente do contrato público.
func (recurso *Emails) Eventos(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	idCodificado, err := segmento(id)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/emails/"+idCodificado+"/events", nil, nil, "json", false, configuracao, nil)
}

// BaixarAnexo executa a operação correspondente do contrato público.
func (recurso *Emails) BaixarAnexo(ctx context.Context, id string, indice int, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	if indice < 0 {
		return nil, erroConfiguracao("Informe um índice de anexo inteiro e não negativo.")
	}
	idCodificado, err := segmento(id)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/emails/"+idCodificado+"/attachments/"+strconv.Itoa(indice), nil, nil, "bytes", false, configuracao, nil)
}

// Lotes reúne operações do contrato público Mailcore.
type Lotes struct {
	transporte executor
}

func novoLotes(transporte executor) *Lotes {
	return &Lotes{
		transporte: transporte,
	}
}

// Enviar executa a operação correspondente do contrato público.
func (recurso *Lotes) Enviar(ctx context.Context, dados any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	corpo, err := lista(dados)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/emails/batch", corpo, nil, "json", false, configuracao, nil)
}

// ObterPorIdempotencia executa a operação correspondente do contrato público.
func (recurso *Lotes) ObterPorIdempotencia(ctx context.Context, chave string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	chaveCodificado, err := segmento(chave)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/emails/batch/idempotency/"+chaveCodificado, nil, nil, "json", false, configuracao, nil)
}
