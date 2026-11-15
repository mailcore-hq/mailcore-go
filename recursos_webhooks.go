// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// EventosWebhook reúne operações do contrato público Mailcore.
type EventosWebhook struct {
	transporte executor
}

func novoEventosWebhook(transporte executor) *EventosWebhook {
	return &EventosWebhook{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *EventosWebhook) Listar(ctx context.Context, webhookID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	webhookIDCodificado, err := segmento(webhookID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/webhooks/"+webhookIDCodificado+"/events", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *EventosWebhook) Obter(ctx context.Context, webhookID string, eventoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	webhookIDCodificado, err := segmento(webhookID)
	if err != nil {
		return nil, err
	}
	eventoIDCodificado, err := segmento(eventoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/webhooks/"+webhookIDCodificado+"/events/"+eventoIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// ListarTentativas executa a operação correspondente do contrato público.
func (recurso *EventosWebhook) ListarTentativas(ctx context.Context, webhookID string, eventoID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	webhookIDCodificado, err := segmento(webhookID)
	if err != nil {
		return nil, err
	}
	eventoIDCodificado, err := segmento(eventoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/webhooks/"+webhookIDCodificado+"/events/"+eventoIDCodificado+"/attempts", nil, filtros, "json", false, configuracao, nil)
}

// Reenviar executa a operação correspondente do contrato público.
func (recurso *EventosWebhook) Reenviar(ctx context.Context, webhookID string, eventoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	webhookIDCodificado, err := segmento(webhookID)
	if err != nil {
		return nil, err
	}
	eventoIDCodificado, err := segmento(eventoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/webhooks/"+webhookIDCodificado+"/events/"+eventoIDCodificado+"/replay", nil, nil, "json", false, configuracao, nil)
}

// Webhooks reúne operações do contrato público Mailcore.
type Webhooks struct {
	transporte executor
	Eventos    *EventosWebhook
}

func novoWebhooks(transporte executor) *Webhooks {
	return &Webhooks{
		transporte: transporte,
		Eventos:    novoEventosWebhook(transporte),
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *Webhooks) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/webhooks", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Webhooks) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/webhooks", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Webhooks) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/webhooks/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *Webhooks) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/webhooks/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *Webhooks) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/webhooks/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}
