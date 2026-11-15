# Mailcore para Go

[Código-fonte](https://github.com/mailcore-hq/mailcore-go) · [Problemas e sugestões](https://github.com/mailcore-hq/mailcore-go/issues)

SDK Go da API Mailcore, com **151 operações**, cliente seguro para uso concorrente,
cancelamento por contexto, downloads por callback e autenticação de webhooks.
Versão publicada **`v0.1.0-alpha.1`**. O núcleo utiliza apenas a biblioteca
padrão e requer **Go 1.26 ou superior**.

## Instalação

Na aplicação consumidora:

```sh
go get github.com/mailcore-hq/mailcore-go@v0.1.0-alpha.1
```

O módulo está disponível pelo proxy público Go, com checksums autenticados pelo
banco público `sum.golang.org`. Confira o [registro da publicação](docs/publicacao.md).
Os usuários do SDK precisam apenas do ambiente Go de sua aplicação; Docker é
utilizado no desenvolvimento e nos ensaios da biblioteca.

## Enviar um email

Disponibilize `MAILCORE_API_KEY` pelo gerenciador de segredos do servidor. Use
remetente autorizado pela conta e uma chave de idempotência estável por intenção:

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/mailcore-hq/mailcore-go"
)

func main() {
    cliente, err := mailcore.NovoCliente("")
    if err != nil { log.Fatal("Configuração Mailcore inválida") }
    defer cliente.Fechar()

    ctx, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancelar()

    resposta, err := cliente.Emails.Enviar(ctx, mailcore.PedidoEmail{
        Remetente: "Equipe <avisos@seu-dominio.com>",
        Destinatarios: []string{"pessoa@example.com"},
        Assunto: "Olá!",
        HTML: "<strong>Sua integração está pronta.</strong>",
    }, mailcore.OpcoesRequisicao{ChaveIdempotencia: "pedido-123-confirmacao"})
    if err != nil { log.Print("Não foi possível concluir a operação"); return }

    var aceita mailcore.RespostaEnvio
    if err := resposta.Decodificar(&aceita); err != nil { return }
    log.Print("Intenção de envio aceita")
}
```

A aceitação não confirma entrega. Acompanhe o estado por eventos/webhooks.
Uma tentativa explícita depois de uma resposta incerta deve conservar a chave
de idempotência e o conteúdo. Uma intenção nova recebe outra chave. A API
define deduplicação, conflitos, permissões e limites de negócio.

O [exemplo executável](exemplos/enviar/main.go) resolve também remetente,
destinatário e idempotência por configuração. Ele envia apenas quando executado
explicitamente com essas variáveis; os verificadores apenas o compilam.

## Recursos e tipos

Os recursos são `Emails`, `Lotes`, `Dominios`, `ChavesAPI`, `Contatos`,
`Segmentos`, `Topicos`, `PropriedadesContatos`, `Modelos`, `Webhooks`, `Campanhas`,
`Automacoes`, `Eventos`, `Recebimento`, `Caixas` e `CLI`. Recursos relacionados
ficam agrupados, como `cliente.Contatos.Importacoes` e `cliente.Modelos.Versoes`.
Consulte [a cobertura](docs/cobertura.md) e o GoDoc das operações.

Todos os métodos recebem contexto primeiro e retornam `(*Resposta, error)`.
Listagens recebem `Consulta` (`nil` significa ausência de filtros) e retornam
uma página; a aplicação controla os cursores. O último argumento opcional é
um único `OpcoesRequisicao`.

Pedidos aceitam structs com tags JSON ou `Documento` (`map[string]any`). Há
tipos de conveniência para envio, anexos, modelos e aceitação. Use `Documento`
para controlar os campos de outros recursos ou distinguir omissão, `nil` (JSON
`null`), `false`, objetos e listas vazias. `Documento{}` produz `{}`;
`[]string{}` produz `[]`; um slice nil num campo produz `null`, conforme Go.
Mapas de chaves numéricas representam objetos JSON. O SDK não altera os pedidos.

`Resposta.Dados` contém os bytes JSON originais, ou texto/bytes nos métodos de
download. `Status`, `Cabecalhos` e `TotalBytes` ficam separados. Em respostas
vazias permitidas, `Dados` é nil. `Decodificar(&destino)` usa `json.Number`
para números em destinos dinâmicos; structs tipadas seguem os tipos declarados.
Isso preserva números grandes, `null`, `false`, `{}` e `[]` sem coerção global.

## Configuração e concorrência

```go
cliente, err := mailcore.NovoCliente(chaveObtidaNoServidor, mailcore.Configuracao{
    URLBase: "https://mailcore.com.br/v1",
    TempoLimite: 30 * time.Second,
    TempoConexao: 10 * time.Second,
    LimiteRespostaBytes: 16 * 1024 * 1024,
})
```

O prazo total da operação usa o menor prazo entre contexto e configuração.
As opções por chamada são `TempoLimite`, `LimiteRespostaBytes`,
`ChaveIdempotencia` e `Cabecalhos`. Zero em tempo/limite utiliza o padrão;
valores negativos são recusados. Cabeçalhos de autenticação, cookies, host e
enquadramento HTTP ficam reservados ao SDK.

Sem chave explícita, `MAILCORE_API_KEY` é consultada em cada operação. Para
rotação ou várias contas, configure `FornecedorChave func(context.Context)
(string, error)`. O fornecedor substitui o ambiente e deve suportar chamadas
concorrentes e respeitar o contexto. Um fornecedor inválido nunca herda a chave
global. Chave explícita e fornecedor juntos são recusados.

Reutilize o cliente entre goroutines. Os corpos, cabeçalhos e buffers de cada
operação são independentes. Não altere os campos públicos de recursos nem os
dados de entrada enquanto estão sendo usados. O SDK não cria uma fila nem
repete operações em um laço. O transporte nativo pode recuperar leituras de
conexões interrompidas; os corpos das mutações não expõem uma função de replay
ao transporte, inclusive quando a operação não tem payload.

O transporte padrão verifica TLS (mínimo 1.2), mantém pool de conexões,
desabilita proxy de ambiente e cookies, recusa redirecionamentos e solicita
codificação `identity`. `CertificadoCA` permite acrescentar uma CA confiável ao
transporte padrão. HTTP explícito de teste fica restrito a loopback com
`PermitirHTTPLocal: true`.

É possível fornecer `ClienteHTTP`. O SDK copia a configuração do cliente,
desabilita o cookie jar e recusa redirecionamentos; o `Transport` injetado
continua pertencendo à aplicação e controla TLS, CA, proxy, buffering e
eventuais repetições próprias. Configure essas políticas no provedor;
`CertificadoCA` e `TempoConexao` configuram apenas o transporte padrão.

`Fechar()` é idempotente: recusa novas operações e fecha conexões ociosas
próprias. Operações já iniciadas terminam pelo seu contexto. Um transporte
injetado não é encerrado pelo SDK.

## Downloads em fluxo

Métodos comuns de download retornam bytes em `Dados`. As variantes terminadas
em `Fluxo` entregam partes ao consumidor síncrono:

```go
resposta, err := cliente.Recebimento.BaixarOriginalFluxo(ctx, recebidoID,
    func(parte []byte) error {
        _, err := arquivoPrivado.Write(parte)
        return err
    }, mailcore.OpcoesRequisicao{LimiteRespostaBytes: 32 * 1024 * 1024})
```

Abra e feche o destino na aplicação e descarte o arquivo parcial em caso de
erro. Para um `io.Writer` genérico, trate também `io.ErrShortWrite` se a gravação
retornar menos bytes que o recebido. A parte só é válida até o callback
retornar: faça uma cópia se precisar guardá-la. O transporte aguarda o consumo
antes da próxima leitura. Uma falha do consumidor encerra a transferência e
retorna o erro original. Nos métodos em fluxo, `TotalBytes` registra o total e
`Dados` fica vazio.

Cancelamento ou timeout interrompe a rede, mas não encerra à força código Go
bloqueado no fornecedor ou consumidor. O limite padrão de resposta é 16 MiB,
inclusive em erros, e pode ser ajustado por chamada. O JSON do pedido tem teto
de 50 MiB. Limites de negócio pertencem à API.

## Erros e webhooks

Use `errors.As(err, &erroMailcore)` com `var erroMailcore *mailcore.Erro` para
consultar `Codigo`, `Status`, `Cabecalhos` e `Mensagem`. A mensagem de negócio
pode conter dados privados. A representação textual do erro é genérica; erros
brutos do provedor não são anexados. `errors.Is` reconhece `context.Canceled`
e `context.DeadlineExceeded` quando essas causas são conhecidas.

`VerificarAssinatura(corpoBruto, cabecalhos, segredo)` autentica HMAC-SHA256
e o instante do webhook, retornando `EventoWebhook`. Use os bytes originais,
antes de decodificar ou reserializar JSON. O segredo usa prefixo `whsec_` e a
janela padrão/máxima é de 300 segundos. A aplicação controla deduplicação e
persistência do processamento. `OpcoesAssinatura` permite fixar o relógio de
teste ou reduzir a janela.

## Desenvolvimento

### Consumir o checkout local

Para desenvolver com o checkout local, adicione no `go.mod` da aplicação:

```go
require github.com/mailcore-hq/mailcore-go v0.1.0-alpha.1

replace github.com/mailcore-hq/mailcore-go => ../mailcore-go
```

O caminho é relativo ao módulo da aplicação. Execute `go mod tidy` depois de
adicionar o import. Para instalar o artefato gerado, em vez de usar `replace`, execute na aplicação
com o caminho absoluto do proxy local:

```sh
GOPROXY=file:///caminho/para/mailcore-go/dist/proxy GOSUMDB=off \
  go get github.com/mailcore-hq/mailcore-go@v0.1.0-alpha.1
```

Esse comando usa apenas o módulo local, sem rede ou dependências externas.
Confira o SHA-256 registrado em `dist/manifesto.json` antes de distribuir o ZIP.

### Verificar a biblioteca

No checkout do SDK:

```sh
sh scripts/verificar.sh
sh scripts/verificar_api.sh
```

O host precisa apenas de shell e Docker Linux. Os scripts verificam Go 1.26 e
1.27, corrida de dados, formato, `go vet`, contrato, transporte e instalação do
ZIP. O ensaio adicional utiliza API/PostgreSQL exclusivos com TLS, sem worker
ou SMTP. Consulte [verificação](docs/verificacao.md), [contribuição](CONTRIBUTING.md)
e [segurança](SECURITY.md) para resultados e requisitos.

Licença [MIT](LICENSE).
