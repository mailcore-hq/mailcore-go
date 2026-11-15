// Código gerado por scripts/gerar_recursos.py a partir das rotas Mailcore.
// Regere as fontes para alterar o contrato; não há geração em runtime.

package mailcore

import "context"

// ConcessoesCLI reúne operações do contrato público Mailcore.
type ConcessoesCLI struct {
	transporte executor
}

func novoConcessoesCLI(transporte executor) *ConcessoesCLI {
	return &ConcessoesCLI{
		transporte: transporte,
	}
}

// Listar executa a operação correspondente do contrato público.
func (recurso *ConcessoesCLI) Listar(ctx context.Context, filtros Consulta, opcoes ...OpcoesRequisicao) (*Resposta, error) {
	if recurso == nil || recurso.transporte == nil {
		return nil, erroConfiguracao("Inicialize o recurso por meio de um cliente Mailcore.")
	}
	configuracao, err := opcao(opcoes)
	if err != nil {
		return nil, err
	}
	return recurso.transporte.executar(ctx, "GET", "/cli/oauth/grants", nil, filtros, "json", false, configuracao, nil)
}

// Revogar executa a operação correspondente do contrato público.
func (recurso *ConcessoesCLI) Revogar(ctx context.Context, id string, opcoes ...OpcoesRequisicao) (*Resposta, error) {
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
	return recurso.transporte.executar(ctx, "POST", "/cli/oauth/grants/"+idCodificado+"/revoke", nil, nil, "json", false, configuracao, nil)
}

// CLIMailcore reúne operações do contrato público Mailcore.
type CLIMailcore struct {
	transporte executor
	Concessoes *ConcessoesCLI
}

func novoCLIMailcore(transporte executor) *CLIMailcore {
	return &CLIMailcore{
		transporte: transporte,
		Concessoes: novoConcessoesCLI(transporte),
	}
}
