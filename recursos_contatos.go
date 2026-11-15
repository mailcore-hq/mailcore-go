// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// SegmentosContato reúne operações do contrato público Mailcore.
type SegmentosContato struct {
	transporte executor
}

func novoSegmentosContato(transporte executor) *SegmentosContato {
	return &SegmentosContato{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *SegmentosContato) Listar(ctx context.Context, contatoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	contatoIDCodificado, err := segmento(contatoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/contacts/"+contatoIDCodificado+"/segments", nil, nil, "json", false, configuracao, nil)
}

// Adicionar executa a operação correspondente do contrato público.
func (recurso *SegmentosContato) Adicionar(ctx context.Context, contatoID string, segmentoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	contatoIDCodificado, err := segmento(contatoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/contacts/"+contatoIDCodificado+"/segments", Documento{"segmento_id": segmentoID}, nil, "json", true, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *SegmentosContato) Remover(ctx context.Context, contatoID string, segmentoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	contatoIDCodificado, err := segmento(contatoID)
	if err != nil {
		return nil, err
	}
	segmentoIDCodificado, err := segmento(segmentoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "DELETE", "/contacts/"+contatoIDCodificado+"/segments/"+segmentoIDCodificado, nil, nil, "json", false, configuracao, nil)
}

// TopicosContato reúne operações do contrato público Mailcore.
type TopicosContato struct {
	transporte executor
}

func novoTopicosContato(transporte executor) *TopicosContato {
	return &TopicosContato{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *TopicosContato) Listar(ctx context.Context, contatoID string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	contatoIDCodificado, err := segmento(contatoID)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/contacts/"+contatoIDCodificado+"/topics", nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *TopicosContato) Atualizar(ctx context.Context, contatoID string, topicos []Documento, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	contatoIDCodificado, err := segmento(contatoID)
	if err != nil {
		return nil, err
	}
	if topicos == nil {
		topicos = make([]Documento, 0)
	}
	return recurso.transporte.executar(ctx, "PATCH", "/contacts/"+contatoIDCodificado+"/topics", Documento{"topicos": topicos}, nil, "json", false, configuracao, nil)
}

// ImportacoesContatos reúne operações do contrato público Mailcore.
type ImportacoesContatos struct {
	transporte executor
}

func novoImportacoesContatos(transporte executor) *ImportacoesContatos {
	return &ImportacoesContatos{
		transporte: transporte,
	}
}

// Prever executa a operação correspondente do contrato público.
func (recurso *ImportacoesContatos) Prever(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/contacts/imports/preview", corpo, nil, "json", false, configuracao, nil)
}

// Criar executa a operação correspondente do contrato público.
func (recurso *ImportacoesContatos) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/contacts/imports", corpo, nil, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *ImportacoesContatos) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/contacts/imports/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// ExportacoesContatos reúne operações do contrato público Mailcore.
type ExportacoesContatos struct {
	transporte executor
}

func novoExportacoesContatos(transporte executor) *ExportacoesContatos {
	return &ExportacoesContatos{
		transporte: transporte,
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *ExportacoesContatos) Criar(ctx context.Context, filtros any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	corpo, err := objeto(filtros)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "POST", "/contacts/exports", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *ExportacoesContatos) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/contacts/exports", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *ExportacoesContatos) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/contacts/exports/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Baixar executa a operação correspondente do contrato público.
func (recurso *ExportacoesContatos) Baixar(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/contacts/exports/"+idCodificado+"/download", nil, nil, "bytes", false, configuracao, nil)
}

// BaixarFluxo entrega blocos sem acumular todo o conteúdo.
// Copie o bloco no consumidor se precisar mantê-lo após o retorno.
func (recurso *ExportacoesContatos) BaixarFluxo(ctx context.Context, id string, consumir func([]byte) error, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/contacts/exports/"+idCodificado+"/download", nil, nil, "fluxo", false, configuracao, consumir)
}

// Contatos reúne operações do contrato público Mailcore.
type Contatos struct {
	transporte  executor
	Segmentos   *SegmentosContato
	Topicos     *TopicosContato
	Importacoes *ImportacoesContatos
	Exportacoes *ExportacoesContatos
}

func novoContatos(transporte executor) *Contatos {
	return &Contatos{
		transporte:  transporte,
		Segmentos:   novoSegmentosContato(transporte),
		Topicos:     novoTopicosContato(transporte),
		Importacoes: novoImportacoesContatos(transporte),
		Exportacoes: novoExportacoesContatos(transporte),
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *Contatos) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/contacts", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Contatos) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/contacts", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Contatos) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/contacts/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *Contatos) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/contacts/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *Contatos) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/contacts/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// AdicionarEmLote executa a operação correspondente do contrato público.
func (recurso *Contatos) AdicionarEmLote(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/contacts/batch", corpo, nil, "json", false, configuracao, nil)
}

// AlterarEmLote executa a operação correspondente do contrato público.
func (recurso *Contatos) AlterarEmLote(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/contacts/bulk", corpo, nil, "json", false, configuracao, nil)
}

// Segmentos reúne operações do contrato público Mailcore.
type Segmentos struct {
	transporte executor
}

func novoSegmentos(transporte executor) *Segmentos {
	return &Segmentos{
		transporte: transporte,
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *Segmentos) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/segments", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Segmentos) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/segments", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Segmentos) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/segments/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *Segmentos) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/segments/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *Segmentos) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/segments/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Topicos reúne operações do contrato público Mailcore.
type Topicos struct {
	transporte executor
}

func novoTopicos(transporte executor) *Topicos {
	return &Topicos{
		transporte: transporte,
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *Topicos) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/topics", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *Topicos) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/topics", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *Topicos) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/topics/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *Topicos) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/topics/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *Topicos) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/topics/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// PropriedadesContatos reúne operações do contrato público Mailcore.
type PropriedadesContatos struct {
	transporte executor
}

func novoPropriedadesContatos(transporte executor) *PropriedadesContatos {
	return &PropriedadesContatos{
		transporte: transporte,
	}
}

// Criar executa a operação correspondente do contrato público.
func (recurso *PropriedadesContatos) Criar(ctx context.Context, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/contact-properties", corpo, nil, "json", false, configuracao, nil)
}

// Listar executa a operação correspondente do contrato público.
func (recurso *PropriedadesContatos) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/contact-properties", nil, filtros, "json", false, configuracao, nil)
}

// Obter executa a operação correspondente do contrato público.
func (recurso *PropriedadesContatos) Obter(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "GET", "/contact-properties/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}

// Atualizar executa a operação correspondente do contrato público.
func (recurso *PropriedadesContatos) Atualizar(ctx context.Context, id string, pedido any, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "PATCH", "/contact-properties/"+idCodificado, corpo, nil, "json", false, configuracao, nil)
}

// Remover executa a operação correspondente do contrato público.
func (recurso *PropriedadesContatos) Remover(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "DELETE", "/contact-properties/"+idCodificado, nil, nil, "json", false, configuracao, nil)
}
