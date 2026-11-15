# Verificação do SDK Go

O laboratório usa dados sintéticos e ferramentas em Docker. Não publica módulo,
não faz deploy e não entrega emails por SMTP. Os relatórios JSON em `docs/`
registram versões, imagens, cenários e resultados observados.

Em 14/11/2026 passaram **201 casos de teste por versão**, em Go 1.26.11 e
1.27.4, com `-race`, sem falhas ou testes ignorados. A contagem representa
folhas; incluindo os grupos de subtestes, são 211 resultados aprovados em
cada versão.
`go vet`, formato, gerador e instalação offline nos dois perfis também passaram.

## Matriz e qualidade

```sh
sh scripts/verificar.sh
```

Execute no checkout. O host precisa de shell POSIX e Docker Linux acessível;
Go e Python rodam nas imagens. A construção precisa de acesso aos registros de
imagens/pacotes de sistema. Testes e instalação do módulo usam `--network none`.

As bases oficiais Go **1.26.11** e **1.27.4** estão fixadas por digest, com
`GOTOOLCHAIN=local` para impedir atualização automática do compilador. Não há
dependências externas no módulo. `go vet`, `gofmt` e o gerador conferem a árvore.
As fontes de recursos são estáticas; nada é gerado ao importar a biblioteca.

Cada `go-1.*.json` precisa indicar resultado aprovado, nenhuma falha e nenhum
teste ignorado. O runner conta folhas separadamente dos grupos de subtestes.
Todas as suítes rodam com `-race` e execução forçada por `-count=1`. A matriz
utiliza dois processadores e limites de memória/processos no container.

O contrato verifica as 151 operações com casos independentes da especificação
de geração. Os testes de transporte usam HTTPS real em loopback para conferir
certificados/hostname, pool, queda de conexão reutilizada, mutações com/sem corpo,
idempotência, concorrência, contexto, prazos, cookies e redirecionamentos.
São exercitados também JSON/UTF-8/números, limites, bytes e streaming, truncamento
por Content-Length e chunked, erros do consumidor e fechamento em falha.
Assinaturas incluem vetor independente, janela temporal, rotação, ambiguidade,
versão e Base64 canônico.

Os testes não precisam de credenciais reais. O programa em `exemplos/enviar`
é compilado, sem executar seu envio. O cliente do ensaio de API também é
compilado antes da execução explícita do laboratório adicional.

## Compilação em outras plataformas

`scripts/verificar_plataformas.py` compila o núcleo com `CGO_ENABLED=0` para
Linux, macOS e Windows, em amd64 e arm64. `plataformas.json` registra o
compilador, a imagem e o resultado por alvo. Isso verifica compilação cruzada;
não executa o SDK nesses sistemas. A execução dos testes permanece Linux em
Docker. O script principal inclui essa verificação no perfil Go 1.27.

## Instalação independente

O estágio de distribuição cria um ZIP no formato de módulo Go com prefixo
`github.com/mailcore-hq/mailcore-go@v0.1.0-alpha.1/`, arquivos `.mod`, `.info` e lista de
versões. `dist/manifesto.json` registra SHA-256 e arquivos incluídos.

O ZIP usa uma seleção de fontes, licença, documentos e exemplos, sem `.git`,
`.env`, testes, scripts ou relatórios JSON. O consumidor temporário usa somente
esse proxy de arquivos, sem checkout, `replace`, sumdb, rede ou dependências.
Ele verifica a versão realmente carregada, 30 recursos, 151 operações, três
callbacks de fluxo, DTO de envio e preservação de JSON/números.

`instalacao-go-1.26.json`, `instalacao-go-1.27.json` e `instalacao-local.json`
registram os dois perfis e o mesmo SHA-256 do módulo instalado. Esses arquivos
precisam indicar resultado aprovado. A prova não publica descoberta remota do
domínio nem habilita `go get` público.

## API e PostgreSQL locais

```sh
sh scripts/verificar_api.sh
```

O ensaio adicional precisa do checkout backend, imagem runtime da API, imagem
Rust e cache público Cargo previamente preparados. Compila a API com
`--offline --locked` em Docker e registra revisão e SHA-256 dos fontes/binário.
Use `MAILCORE_BACKEND` para outro checkout e `MAILCORE_IMAGEM_CLIENTE` para outra
imagem do SDK. As opções adicionais são listadas por `--help` no script.

O controlador utiliza o socket Docker para criar/remover serviços exclusivos,
sem reutilizar o banco da aplicação em execução. A rede é interna, sem portas
publicadas, worker, Postal ou SMTP. Segredos/certificados/configuração sintéticos
ficam em volume temporário em memória. HTTPS e a identidade da API são
verificados; PostgreSQL utiliza TLS com `verify-full`.

Os dez grupos cobrem permissões Full/Sending, composição/CID/Cc/Cco, replay e
conflito por idempotência, lote atômico, domínios, contatos/catálogos/importação
aceita/exportação CSV, modelos/versões/publicação com revisão, marketing sem
disparo, caixas/etiquetas/rascunhos, recebimento MIME e revogação CLI.
Quatro callbacks de CSV/MIME usam o transporte próprio e um cliente HTTP
injetado; leituras simultâneas usam o mesmo cliente.

Uma execução aprovada exige `resultado: "aprovado"` e limpeza concluída. As
contagens esperadas são três mensagens, três tarefas e um recebido pendente,
sem processamento de importação ou recebimento por worker. Elas comprovam a
aceitação/persistência exercitada e não a entrega de emails. Os relatórios
`integracao-api.json` e `integracao-api-interrupcao.json` distinguem a execução
normal da interrupção SIGTERM esperada, com limpeza em ambos os casos.

## Limites da evidência

Esses ensaios não medem desempenho em produção nem comprovam instalação em toda
combinação de sistema/provedor HTTP, entrega SMTP, implantação cloud ou
publicação. Um provedor injetado controla sua política de rede e seu buffer.
O núcleo Go não exige framework ou cgo; o detector de corridas usado no
desenvolvimento precisa do compilador C presente nas imagens de teste.
