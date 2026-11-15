# Contribuir

Esta alpha é mantida na branch `codex/sdk-go`. Escreva código, testes e
documentação em português do Brasil, preservando nomes definidos pelo protocolo
e as convenções de exportação/initialisms do Go.

Confirme cada mudança no contrato atual da API: método, caminho, DTO, escopo,
revisão e resposta. Uma rota de API não comprova a existência de uma tela.
Os recursos são fontes estáticas; o gerador em `scripts/gerar_recursos.py`
precisa permanecer consistente com as classes e os testes de contrato.

Use contexto primeiro e preserve objetos/listas, campos omitidos, null, false,
números grandes e bytes binários. Não reescreva nomes ou envelopes do servidor.
Mantenha o transporte centralizado e as operações concorrentes independentes.
Não introduza estado global de credenciais ou configuração mutável por requisição.

Execute em Docker, no checkout:

```sh
sh scripts/verificar.sh
sh scripts/verificar_api.sh
```

Go, Python e ferramentas rodam nas imagens. Os ensaios de protocolo precisam
usar HTTP real para comprovar TLS, limites, cancelamento, pooling e comportamento
de replay; um coletor de pedidos comprova apenas o contrato das operações.
Alterações de concorrência devem passar em `go test -race`.

Use dados sintéticos, serviços Docker exclusivos, rede interna e nenhuma porta
publicada. Não reutilize banco, token ou conta de produção. Registre os recursos
criados e remova-os também em falha/SIGTERM. A integração local compila a API
com dependências públicas previamente cacheadas e verifica seus hashes.

Não publique módulo, tag remota ou aplicação por meio dos verificadores.
Atualize cobertura, histórico e os resultados pertinentes à mudança.
Relate vulnerabilidades conforme [SECURITY.md](SECURITY.md).
