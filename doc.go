// Package mailcore integra aplicações Go ao contrato público da API Mailcore.
// O cliente pode ser compartilhado entre goroutines e aceita context.Context
// em todas as operações. O chamador mantém a chave de idempotência por intenção.
package mailcore
