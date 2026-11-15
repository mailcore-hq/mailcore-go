// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// Recebimento reúne operações do contrato público Mailcore.
type Recebimento struct {
	transporte executor
}

func novoRecebimento(transporte executor) *Recebimento {
	return &Recebimento{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Recebimento) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/receiving", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Recebimento) Obter(ctx context.Context, id string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/receiving/"+idCodificado, nil, filtros, "json", false, configuracao, nil)
}

// ListarAnexos executa a operação correspondente do contrato público.
func (recurso *Recebimento) ListarAnexos(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/receiving/"+idCodificado+"/attachments", nil, nil, "json", false, configuracao, nil)
}

// ObterAnexo executa a operação correspondente do contrato público.
func (recurso *Recebimento) ObterAnexo(ctx context.Context, id string, anexoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	anexoIDCodificado, err := segmento(anexoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/receiving/"+idCodificado+"/attachments/"+anexoIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// BaixarOriginal executa a operação correspondente do contrato público.
func (recurso *Recebimento) BaixarOriginal(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/receiving/"+idCodificado+"/raw", nil, nil, "bytes", false, configuracao, nil)
}

// BaixarAnexo executa a operação correspondente do contrato público.
func (recurso *Recebimento) BaixarAnexo(ctx context.Context, id string, anexoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	anexoIDCodificado, err := segmento(anexoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/receiving/"+idCodificado+"/attachments/"+anexoIDCodificado+"/download", nil, nil, "bytes", false, configuracao, nil)
}

// BaixarOriginalFluxo entrega blocos sem acumular todo o conteúdo.
// Copie o bloco no consumidor se precisar mantê-lo após o retorno.
func (recurso *Recebimento) BaixarOriginalFluxo(ctx context.Context, id string, consumir func([]byte) error, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	if consumir == nil {
		return nil, erroConfiguracao("Forneça um consumidor para o download em fluxo.")
	}
	idCodificado, err := segmento(id)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/receiving/"+idCodificado+"/raw", nil, nil, "fluxo", false, configuracao, consumir)
}

// BaixarAnexoFluxo entrega blocos sem acumular todo o conteúdo.
// Copie o bloco no consumidor se precisar mantê-lo após o retorno.
func (recurso *Recebimento) BaixarAnexoFluxo(ctx context.Context, id string, anexoID string, consumir func([]byte) error, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	if consumir == nil {
		return nil, erroConfiguracao("Forneça um consumidor para o download em fluxo.")
	}
	idCodificado, err := segmento(id)
	if err != nil {
		return nil, err
	}
	anexoIDCodificado, err := segmento(anexoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/receiving/"+idCodificado+"/attachments/"+anexoIDCodificado+"/download", nil, nil, "fluxo", false, configuracao, consumir)
}
