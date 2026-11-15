// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// EmailsCaixa reúne operações do contrato público Mailcore.
type EmailsCaixa struct {
	transporte executor
}

func novoEmailsCaixa(transporte executor) *EmailsCaixa {
	return &EmailsCaixa{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *EmailsCaixa) Listar(ctx context.Context, caixaID string, conversaID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado+"/emails", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *EmailsCaixa) Obter(ctx context.Context, caixaID string, conversaID string, emailID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	emailIDCodificado, err := segmento(emailID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado+"/emails/"+emailIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// Responder executa a operação correspondente do contrato público.
func (recurso *EmailsCaixa) Responder(ctx context.Context, caixaID string, conversaID string, emailID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	emailIDCodificado, err := segmento(emailID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado+"/emails/"+emailIDCodificado+"/reply", corpo, nil, "json", false, configuracao, nil)
}

// Encaminhar executa a operação correspondente do contrato público.
func (recurso *EmailsCaixa) Encaminhar(ctx context.Context, caixaID string, conversaID string, emailID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	emailIDCodificado, err := segmento(emailID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado+"/emails/"+emailIDCodificado+"/forward", corpo, nil, "json", false, configuracao, nil)
}

// EtiquetasConversaCaixa reúne operações do contrato público Mailcore.
type EtiquetasConversaCaixa struct {
	transporte executor
}

func novoEtiquetasConversaCaixa(transporte executor) *EtiquetasConversaCaixa {
	return &EtiquetasConversaCaixa{
		transporte: transporte,
	}
}

// Adicionar executa a operação correspondente do contrato público.
func (recurso *EtiquetasConversaCaixa) Adicionar(ctx context.Context, caixaID string, conversaID string, ids []string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = make([]string, 0)
	}
	return recurso.transporte.executar(ctx, "POST", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado+"/labels", Documento{"ids": ids}, nil, "json", false, configuracao, nil)
}

// Substituir executa a operação correspondente do contrato público.
func (recurso *EtiquetasConversaCaixa) Substituir(ctx context.Context, caixaID string, conversaID string, ids []string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = make([]string, 0)
	}
	return recurso.transporte.executar(ctx, "PUT", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado+"/labels", Documento{"ids": ids}, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *EtiquetasConversaCaixa) Remover(ctx context.Context, caixaID string, conversaID string, etiquetaID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	etiquetaIDCodificado, err := segmento(etiquetaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "DELETE", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado+"/labels/"+etiquetaIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// ConversasCaixa reúne operações do contrato público Mailcore.
type ConversasCaixa struct {
	transporte executor
	Emails     *EmailsCaixa
	Etiquetas  *EtiquetasConversaCaixa
}

func novoConversasCaixa(transporte executor) *ConversasCaixa {
	return &ConversasCaixa{
		transporte: transporte,
		Emails:     novoEmailsCaixa(transporte),
		Etiquetas:  novoEtiquetasConversaCaixa(transporte),
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *ConversasCaixa) Listar(ctx context.Context, caixaID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/threads", nil, filtros, "json", false, configuracao, nil)
}

// Buscar executa a operação correspondente do contrato público.
func (recurso *ConversasCaixa) Buscar(ctx context.Context, caixaID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/threads/search", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *ConversasCaixa) Obter(ctx context.Context, caixaID string, conversaID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *ConversasCaixa) Atualizar(ctx context.Context, caixaID string, conversaID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "PATCH", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *ConversasCaixa) Remover(ctx context.Context, caixaID string, conversaID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	conversaIDCodificado, err := segmento(conversaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "DELETE", "/inboxes/"+caixaIDCodificado+"/threads/"+conversaIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// EtiquetasCaixa reúne operações do contrato público Mailcore.
type EtiquetasCaixa struct {
	transporte executor
}

func novoEtiquetasCaixa(transporte executor) *EtiquetasCaixa {
	return &EtiquetasCaixa{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *EtiquetasCaixa) Listar(ctx context.Context, caixaID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/labels", nil, nil, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *EtiquetasCaixa) Obter(ctx context.Context, caixaID string, etiquetaID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	etiquetaIDCodificado, err := segmento(etiquetaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/labels/"+etiquetaIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// Criar executa a operação correspondente do contrato público.
func (recurso *EtiquetasCaixa) Criar(ctx context.Context, caixaID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/inboxes/"+caixaIDCodificado+"/labels", corpo, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *EtiquetasCaixa) Atualizar(ctx context.Context, caixaID string, etiquetaID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	etiquetaIDCodificado, err := segmento(etiquetaID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "PATCH", "/inboxes/"+caixaIDCodificado+"/labels/"+etiquetaIDCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *EtiquetasCaixa) Remover(ctx context.Context, caixaID string, etiquetaID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	etiquetaIDCodificado, err := segmento(etiquetaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "DELETE", "/inboxes/"+caixaIDCodificado+"/labels/"+etiquetaIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// RascunhosCaixa reúne operações do contrato público Mailcore.
type RascunhosCaixa struct {
	transporte executor
}

func novoRascunhosCaixa(transporte executor) *RascunhosCaixa {
	return &RascunhosCaixa{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *RascunhosCaixa) Listar(ctx context.Context, caixaID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/drafts", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *RascunhosCaixa) Obter(ctx context.Context, caixaID string, rascunhoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	rascunhoIDCodificado, err := segmento(rascunhoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+caixaIDCodificado+"/drafts/"+rascunhoIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// Criar executa a operação correspondente do contrato público.
func (recurso *RascunhosCaixa) Criar(ctx context.Context, caixaID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/inboxes/"+caixaIDCodificado+"/drafts", corpo, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *RascunhosCaixa) Atualizar(ctx context.Context, caixaID string, rascunhoID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	rascunhoIDCodificado, err := segmento(rascunhoID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "PATCH", "/inboxes/"+caixaIDCodificado+"/drafts/"+rascunhoIDCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *RascunhosCaixa) Remover(ctx context.Context, caixaID string, rascunhoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	rascunhoIDCodificado, err := segmento(rascunhoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "DELETE", "/inboxes/"+caixaIDCodificado+"/drafts/"+rascunhoIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// Enviar executa a operação correspondente do contrato público.
func (recurso *RascunhosCaixa) Enviar(ctx context.Context, caixaID string, rascunhoID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	caixaIDCodificado, err := segmento(caixaID)
	if err != nil {
		return nil, err
	}
	rascunhoIDCodificado, err := segmento(rascunhoID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/inboxes/"+caixaIDCodificado+"/drafts/"+rascunhoIDCodificado+"/send", corpo, nil, "json", false, configuracao, nil)
}

// Caixas reúne operações do contrato público Mailcore.
type Caixas struct {
	transporte executor
	Conversas  *ConversasCaixa
	Etiquetas  *EtiquetasCaixa
	Rascunhos  *RascunhosCaixa
}

func novoCaixas(transporte executor) *Caixas {
	return &Caixas{
		transporte: transporte,
		Conversas:  novoConversasCaixa(transporte),
		Etiquetas:  novoEtiquetasCaixa(transporte),
		Rascunhos:  novoRascunhosCaixa(transporte),
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *Caixas) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/inboxes", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Caixas) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/inboxes", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Caixas) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/inboxes/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *Caixas) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/inboxes/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *Caixas) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/inboxes/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}
