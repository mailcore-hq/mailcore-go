// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// ExecucoesAutomacoes reúne operações do contrato público Mailcore.
type ExecucoesAutomacoes struct {
	transporte executor
}

func novoExecucoesAutomacoes(transporte executor) *ExecucoesAutomacoes {
	return &ExecucoesAutomacoes{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *ExecucoesAutomacoes) Listar(ctx context.Context, automacaoID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	automacaoIDCodificado, err := segmento(automacaoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/automations/"+automacaoIDCodificado+"/runs", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *ExecucoesAutomacoes) Obter(ctx context.Context, automacaoID string, execucaoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	automacaoIDCodificado, err := segmento(automacaoID)
	if err != nil {
		return nil, err
	}
	execucaoIDCodificado, err := segmento(execucaoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/automations/"+automacaoIDCodificado+"/runs/"+execucaoIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// Automacoes reúne operações do contrato público Mailcore.
type Automacoes struct {
	transporte executor
	Execucoes  *ExecucoesAutomacoes
}

func novoAutomacoes(transporte executor) *Automacoes {
	return &Automacoes{
		transporte: transporte,
		Execucoes:  novoExecucoesAutomacoes(transporte),
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *Automacoes) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/automations", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Automacoes) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/automations", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Automacoes) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/automations/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *Automacoes) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/automations/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *Automacoes) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/automations/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Duplicar executa a operação correspondente do contrato público.
func (recurso *Automacoes) Duplicar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/automations/"+idCodificado+"/duplicate", corpo, nil, "json", false, configuracao, nil)
}

// Parar executa a operação correspondente do contrato público.
func (recurso *Automacoes) Parar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/automations/"+idCodificado+"/stop", corpo, nil, "json", false, configuracao, nil)
}
