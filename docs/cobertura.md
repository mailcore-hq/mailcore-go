# Cobertura do SDK Go

A alpha `v0.1.0-alpha.1` oferece **151 operações**, agrupadas em 16 recursos
principais e 30 estruturas de recursos. A contagem inclui variantes de download;
não representa 151 URLs distintas nem toda a API do produto.

| Recurso | Superfície | Operações |
| --- | --- | ---: |
| `Emails` | Envio, consulta, agendamento/cancelamento, eventos, idempotência, anexo | 8 |
| `Lotes` | Envio atômico e consulta por idempotência | 2 |
| `Dominios` | CRUD, capacidades, rastreamento, DNS e verificação | 7 |
| `ChavesAPI` | Criação, listagem, renomeação e revogação | 4 |
| `Contatos` | CRUD, lote, segmentos, tópicos, importações e exportações | 20 |
| `Segmentos` | Criação, listagem, consulta, edição e remoção | 5 |
| `Topicos` | Criação, listagem, consulta, edição e remoção | 5 |
| `PropriedadesContatos` | Criação, listagem, consulta, edição e remoção | 5 |
| `Modelos` | CRUD, publicação, duplicação, prévias, testes, versões e pastas | 18 |
| `Webhooks` | CRUD, eventos, tentativas e reenvio | 9 |
| `Campanhas` | CRUD, duplicação, revisão, envio/cancelamento, prévias, testes, métricas | 15 |
| `Automacoes` | CRUD, duplicação, parada e execuções | 9 |
| `Eventos` | Definição, listagem, consulta, edição, remoção e ocorrência | 6 |
| `Recebimento` | Consulta, anexos, MIME e downloads comuns/em fluxo | 8 |
| `Caixas` | CRUD, etiquetas, conversas, mensagens, respostas, encaminhamento, rascunhos | 28 |
| `CLI` | Listagem e revogação de concessões | 2 |
| **Total** | | **151** |

Os nomes são portugueses e seguem initialisms Go, como `ChavesAPI`, `RegistrosDNS`
e `CLI.Concessoes`. Recursos subordinados são campos: `Contatos.Importacoes`,
`Contatos.Exportacoes`, `Modelos.Versoes`, `Webhooks.Eventos`,
`Caixas.Conversas.Emails` e outros. Métodos públicos e parâmetros ficam visíveis
no GoDoc e nas fontes estáticas `recursos_*.go`.

## Contrato e dados

Cada operação recebe contexto e um único conjunto opcional de opções, retornando
`(*Resposta, error)`. Métodos de listagem/consulta com filtros recebem `Consulta`;
use nil para não enviar filtros. Os cursores e envelopes são os do servidor.
Não há paginação automática nem fila de envios no SDK.

Pedidos recebem structs com tags JSON ou `Documento`. Objetos raiz mantêm `{}`;
listas mantêm `[]`. Valores nil dentro dos campos são JSON null, false é
preservado e campos omitidos não são criados pelo SDK. Mapas numéricos continuam
objetos. Structs seguem `omitempty` das próprias tags; use Documento quando for
necessário explicitar uma combinação de ausência/null/zero.

`PedidoEmail`, `AnexoEmail`, `ModeloEmail` e `RespostaEnvio` oferecem conveniência
tipada para envio. Os demais DTOs permanecem no contrato da API, sem renomear
campos ou impedir campos adicionais por uma normalização local.

As respostas mantêm JSON em `json.RawMessage`. `Decodificar` usa `json.Number`
em destinos dinâmicos para evitar arredondamento de inteiros grandes. Texto
precisa de UTF-8 válido e downloads binários mantêm os bytes. As variantes em
fluxo usam callback síncrono, devolvem `TotalBytes` e não acumulam Dados.

O helper de webhooks autentica HMAC sobre o corpo bruto e não entra na contagem
das operações HTTP. A aplicação é responsável pela deduplicação dos eventos.

## Transporte

O transporte próprio usa net/http, TLS verificado, contexto com prazo total,
pool de conexões, limite durante leitura e nenhuma camada de repetição no SDK.
Não segue redirects, proxies de ambiente ou URLs temporárias dos metadados.
O corpo das mutações não oferece `GetBody` ao transporte; mutações sem payload
também mantêm um corpo opaco, evitando a condição comum de replay em conexão
reutilizada com idempotência. A pilha HTTP pode ter recuperação própria em
leituras e em condições de protocolo; preserve idempotência no servidor.

Um `ClienteHTTP` injetado mantém seu próprio Transport e políticas de rede.
O SDK não fecha esse provedor nem consegue limitar buffers que ele preencha
antes de disponibilizar a resposta. Veja [segurança](../SECURITY.md).

## Escopo desta alpha

O SDK integra o contrato público atualmente coberto pelas outras bibliotecas
Mailcore. Não inclui telas, cliente SMTP, CLI executável, renderer React,
autorização OAuth interativa, SDK específico de framework ou helpers globais
de sessão/dashboard, logs, uso ou supressões. Uma capacidade HTTP não cria
automaticamente um fluxo de interface.

Go 1.26 e 1.27 formam a matriz de verificação. Consulte [os resultados](verificacao.md)
para contagens, versões, imagens e limites efetivamente observados. Os ensaios
locais não comprovam entrega SMTP, publicação remota ou desempenho em produção.
