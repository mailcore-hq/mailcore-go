# Segurança e privacidade

Use o SDK no servidor e guarde credenciais no gerenciador de segredos da
aplicação. Limite escopos e domínios. Não publique chaves, corpos de email,
anexos, idempotência ou dados de clientes em logs, issues ou URLs.

O transporte padrão mantém TLS verificado, não segue redirecionamentos,
desabilita proxies de ambiente e não mantém cookies. Autoridades privadas
devem ser fornecidas por arquivo confiável. Não desabilite verificação de
certificado ou identidade. Ao injetar um cliente HTTP, configure também essas
políticas em seu Transport, inclusive prazos, buffering e tentativas próprias.

O limite de leitura vale para os bytes consumidos pelo SDK. Um provedor pode
armazenar a resposta antes desse consumo; a política dele é da aplicação.
O cliente não busca URLs de anexos/metadados automaticamente: utilize os
downloads autenticados dos recursos.

Não transforme entradas de usuários em URL base, chave, cliente HTTP ou
destinatários livres de um endpoint aberto. A autorização dos recursos
continua pertencendo à API. Em múltiplas contas, resolva a credencial pelo
contexto autenticado da aplicação e use um fornecedor concorrente. O fornecedor
inválido falha sem herdar a variável de ambiente global.

O SDK não registra pedidos. Erros de negócio podem preservar mensagens privadas
em `Erro.Mensagem`; evite publicá-las automaticamente. Revise também a política
de logs/profiler dos provedores e da aplicação. As representações de cliente,
erro, resposta e evento omitem seu conteúdo, mas isso não altera serializações
explícitas realizadas pela aplicação.

Preserve idempotência e conteúdo ao decidir repetir uma intenção depois de uma
resposta incerta. Cancelar a espera não desfaz uma operação já aceita. O SDK não
possui laço de repetição, mas a pilha HTTP ou a fila da aplicação pode ter políticas
próprias. A garantia de deduplicação pertence ao contrato do servidor.

Nos fluxos, uma parte só permanece válida durante o callback. Faça uma cópia
antes de armazená-la em outra goroutine. Proteja o destino, use permissões e
criptografia adequadas, controle retenção e descarte arquivos incompletos.

Valide webhooks sobre os bytes brutos antes de executar efeitos. Limite a
entrada também no servidor HTTP. HMAC e janela temporal não substituem
deduplicação durável; preserve o identificador autenticado do evento.

Relate vulnerabilidades pelo canal privado de manutenção do projeto, com versão
afetada e reprodução sintética. Não inclua credenciais ou dados reais.
