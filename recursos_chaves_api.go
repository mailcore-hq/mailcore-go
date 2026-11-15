// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// ChavesAPI reúne operações do contrato público Mailcore.
type ChavesAPI struct {
	transporte executor
}

func novoChavesAPI(transporte executor) *ChavesAPI {
	return &ChavesAPI{
		transporte: transporte,
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *ChavesAPI) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/api-keys", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *ChavesAPI) Listar(ctx context.Context, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/api-keys", nil, nil, "json", false, configuracao, nil)
}

// Renomear executa a operação correspondente do contrato público.
func (recurso *ChavesAPI) Renomear(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/api-keys/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Revogar executa a operação correspondente do contrato público.
func (recurso *ChavesAPI) Revogar(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/api-keys/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}
