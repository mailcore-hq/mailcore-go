# Histórico

## v0.1.0-alpha.1 — 14/11/2026

- Primeira alpha local com 151 operações em 16 recursos principais.
- Contexto obrigatório, cliente compartilhável entre goroutines e pool HTTP.
- Corpo JSON por structs ou Documento e resposta em bytes originais, com
  decodificação opcional e json.Number em destinos dinâmicos.
- HTTPS verificado, prazos/limites, erros sanitizados, autenticação por instância
  ou fornecedor concorrente e propagação de idempotência.
- Downloads binários e por callback com consumo síncrono e fechamento em falha.
- Verificação HMAC de webhooks sobre corpo bruto e janela temporal.
- Desenvolvimento Docker, matriz Go 1.26/1.27, ensaio local de API/PostgreSQL e
  instalação do módulo via proxy de arquivos independente do checkout.

Versão ainda sem publicação. A validação local não comprova entrega SMTP,
implantação cloud ou desempenho em produção. Consulte [verificação](docs/verificacao.md).
