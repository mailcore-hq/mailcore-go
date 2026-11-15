// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// VersoesModelo reúne operações do contrato público Mailcore.
type VersoesModelo struct {
	transporte executor
}

func novoVersoesModelo(transporte executor) *VersoesModelo {
	return &VersoesModelo{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *VersoesModelo) Listar(ctx context.Context, modeloID string, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	modeloIDCodificado, err := segmento(modeloID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/templates/"+modeloIDCodificado+"/versions", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *VersoesModelo) Obter(ctx context.Context, modeloID string, versaoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	modeloIDCodificado, err := segmento(modeloID)
	if err != nil {
		return nil, err
	}
	versaoIDCodificado, err := segmento(versaoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/templates/"+modeloIDCodificado+"/versions/"+versaoIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// Restaurar executa a operação correspondente do contrato público.
func (recurso *VersoesModelo) Restaurar(ctx context.Context, modeloID string, versaoID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	modeloIDCodificado, err := segmento(modeloID)
	if err != nil {
		return nil, err
	}
	versaoIDCodificado, err := segmento(versaoID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/templates/"+modeloIDCodificado+"/versions/"+versaoIDCodificado+"/restore", corpo, nil, "json", false, configuracao, nil)
}

// Prever executa a operação correspondente do contrato público.
func (recurso *VersoesModelo) Prever(ctx context.Context, modeloID string, versaoID string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	modeloIDCodificado, err := segmento(modeloID)
	if err != nil {
		return nil, err
	}
	versaoIDCodificado, err := segmento(versaoID)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(pedido)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/templates/"+modeloIDCodificado+"/versions/"+versaoIDCodificado+"/preview", corpo, nil, "json", false, configuracao, nil)
}

// PastasModelo reúne operações do contrato público Mailcore.
type PastasModelo struct {
	transporte executor
}

func novoPastasModelo(transporte executor) *PastasModelo {
	return &PastasModelo{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *PastasModelo) Listar(ctx context.Context, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/template-folders", nil, nil, "json", false, configuracao, nil)
}

// Criar executa a operação correspondente do contrato público.
func (recurso *PastasModelo) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/template-folders", corpo, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *PastasModelo) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/template-folders/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *PastasModelo) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/template-folders/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Modelos reúne operações do contrato público Mailcore.
type Modelos struct {
	transporte executor
	Versoes    *VersoesModelo
	Pastas     *PastasModelo
}

func novoModelos(transporte executor) *Modelos {
	return &Modelos{
		transporte: transporte,
		Versoes:    novoVersoesModelo(transporte),
		Pastas:     novoPastasModelo(transporte),
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *Modelos) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/templates", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Modelos) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/templates", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Modelos) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/templates/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *Modelos) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/templates/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *Modelos) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/templates/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Prever executa a operação correspondente do contrato público.
func (recurso *Modelos) Prever(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/templates/"+idCodificado+"/preview", corpo, nil, "json", false, configuracao, nil)
}

// ObterPrevia executa a operação correspondente do contrato público.
func (recurso *Modelos) ObterPrevia(ctx context.Context, organizacaoID string, previaID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	organizacaoIDCodificado, err := segmento(organizacaoID)
	if err != nil {
		return nil, err
	}
	previaIDCodificado, err := segmento(previaID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/template-previews/"+organizacaoIDCodificado+"/"+previaIDCodificado, nil, nil, "texto", false, configuracao, nil)
}

// Testar executa a operação correspondente do contrato público.
func (recurso *Modelos) Testar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/templates/"+idCodificado+"/test", corpo, nil, "json", false, configuracao, nil)
}

// Publicar executa a operação correspondente do contrato público.
func (recurso *Modelos) Publicar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/templates/"+idCodificado+"/publish", corpo, nil, "json", false, configuracao, nil)
}

// Duplicar executa a operação correspondente do contrato público.
func (recurso *Modelos) Duplicar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/templates/"+idCodificado+"/duplicate", corpo, nil, "json", false, configuracao, nil)
}
