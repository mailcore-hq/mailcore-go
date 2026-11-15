# Publicação do módulo Go

A versão **`v0.1.0-alpha.1`** foi publicada em **15/11/2026** no repositório
[mailcore-hq/mailcore-go](https://github.com/mailcore-hq/mailcore-go).

- Módulo: `github.com/mailcore-hq/mailcore-go`.
- Tag anotada: [`v0.1.0-alpha.1`](https://github.com/mailcore-hq/mailcore-go/tree/v0.1.0-alpha.1).
- Objeto da tag: `cc4b06b8ab9aa5839d2d8bd1f8882b42eda818c7`.
- Commit da versão: [`4d68db0942e189b2aed5fa60c303b64c0b35aa8e`](https://github.com/mailcore-hq/mailcore-go/commit/4d68db0942e189b2aed5fa60c303b64c0b35aa8e).

## Instalação pública

```sh
go get github.com/mailcore-hq/mailcore-go@v0.1.0-alpha.1
```

O download foi verificado com `GOPROXY=https://proxy.golang.org` e
`GOSUMDB=sum.golang.org`, partindo de cache vazio. A origem devolvida pelo
proxy identificou a tag e o commit acima. O comando `go mod verify` confirmou
`all modules verified`.

O [registro público de checksums](https://sum.golang.org/lookup/github.com/mailcore-hq/mailcore-go@v0.1.0-alpha.1)
foi consultado e conferido com os valores autenticados pelo cliente Go:

| Conteúdo | Checksum |
|---|---|
| Módulo | `h1:kvgwjJUH2NacWMqisDx5hhE9OTKvIXmkr3dlUTlRigA=` |
| `go.mod` | `h1:tjY19Xr1nPmiFg3CyuFkmsTIGZ4b3cblQ2pH4WddZHY=` |

## Consumidor verificado

Um projeto externo ao checkout instalou o módulo público, compilou e executou
com **Go 1.27.4**, encerrando com código **0**. O consumidor conferiu 30 recursos,
151 métodos e três callbacks de download, além de pedidos tipados, preservação
de JSON e proteção das representações de dados sensíveis.

As operações do consumidor usaram transporte HTTP sintético. Este ensaio não
enviou e-mails reais e não mede entrega SMTP nem desempenho de produção.
A versão alpha não oferece garantia de estabilidade da API pública. O módulo
requer Go 1.26 ou superior; a instalação pública descrita aqui foi executada
com Go 1.27.4.

A tag publicada permanece imutável. Novos ajustes documentais na branch `main`
não substituem o conteúdo autenticado dessa versão.
