package mailcore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Expectativas extraídas por AST do SDK Node a234dc2, independentemente do gerador Go.
const contratosRecursos = `[
{"classe":"ExecucoesAutomacoes","metodo":"listar","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["automacaoId","filtros"],"verbo":"GET","caminho":"/automations/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/runs","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"ExecucoesAutomacoes","metodo":"obter","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["automacaoId","execucaoId"],"verbo":"GET","caminho":"/automations/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/runs/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Automacoes","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/automations","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Automacoes","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["filtros"],"verbo":"GET","caminho":"/automations","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Automacoes","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/automations/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Automacoes","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"PATCH","caminho":"/automations/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Automacoes","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/automations/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Automacoes","metodo":"duplicar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/automations/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/duplicate","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Automacoes","metodo":"parar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/automations/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/stop","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"EmailsCaixa","metodo":"listar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["caixaId","conversaId","paginacao"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/emails","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"EmailsCaixa","metodo":"obter","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test","parte-2/ação+?@example.test"],"parametros":["caixaId","conversaId","emailId"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/emails/parte-2%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"EmailsCaixa","metodo":"responder","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test","parte-2/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","conversaId","emailId","pedido"],"verbo":"POST","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/emails/parte-2%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/reply","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"EmailsCaixa","metodo":"encaminhar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test","parte-2/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","conversaId","emailId","pedido"],"verbo":"POST","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/emails/parte-2%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/forward","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"EtiquetasConversaCaixa","metodo":"adicionar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",["etiqueta-a","etiqueta-b"]],"parametros":["caixaId","conversaId","ids"],"verbo":"POST","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels","kwargs":{"corpo":{"ids":["etiqueta-a","etiqueta-b"]}}},
{"classe":"EtiquetasConversaCaixa","metodo":"substituir","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",["etiqueta-a","etiqueta-b"]],"parametros":["caixaId","conversaId","ids"],"verbo":"PUT","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels","kwargs":{"corpo":{"ids":["etiqueta-a","etiqueta-b"]}}},
{"classe":"EtiquetasConversaCaixa","metodo":"remover","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test","parte-2/ação+?@example.test"],"parametros":["caixaId","conversaId","etiquetaId"],"verbo":"DELETE","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels/parte-2%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"ConversasCaixa","metodo":"listar","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["caixaId","filtros"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"ConversasCaixa","metodo":"buscar","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["caixaId","filtros"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/search","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"ConversasCaixa","metodo":"obter","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["caixaId","conversaId"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"ConversasCaixa","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","conversaId","pedido"],"verbo":"PATCH","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"ConversasCaixa","metodo":"remover","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["caixaId","conversaId"],"verbo":"DELETE","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/threads/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"EtiquetasCaixa","metodo":"listar","argumentos":["parte-0/ação+?@example.test"],"parametros":["caixaId"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels","kwargs":{}},
{"classe":"EtiquetasCaixa","metodo":"obter","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["caixaId","etiquetaId"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"EtiquetasCaixa","metodo":"criar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","pedido"],"verbo":"POST","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"EtiquetasCaixa","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","etiquetaId","pedido"],"verbo":"PATCH","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"EtiquetasCaixa","metodo":"remover","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["caixaId","etiquetaId"],"verbo":"DELETE","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/labels/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"RascunhosCaixa","metodo":"listar","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["caixaId","paginacao"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/drafts","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"RascunhosCaixa","metodo":"obter","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["caixaId","rascunhoId"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/drafts/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"RascunhosCaixa","metodo":"criar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","pedido"],"verbo":"POST","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/drafts","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"RascunhosCaixa","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","rascunhoId","pedido"],"verbo":"PATCH","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/drafts/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"RascunhosCaixa","metodo":"remover","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["caixaId","rascunhoId"],"verbo":"DELETE","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/drafts/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"RascunhosCaixa","metodo":"enviar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","rascunhoId","pedido"],"verbo":"POST","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/drafts/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/send","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Caixas","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["filtros"],"verbo":"GET","caminho":"/inboxes","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Caixas","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["caixaId"],"verbo":"GET","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Caixas","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/inboxes","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Caixas","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["caixaId","pedido"],"verbo":"PATCH","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Caixas","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["caixaId"],"verbo":"DELETE","caminho":"/inboxes/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Campanhas","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/broadcasts","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["filtros"],"verbo":"GET","caminho":"/broadcasts","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Campanhas","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Campanhas","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"PATCH","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Campanhas","metodo":"duplicar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/duplicate","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"revisar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/review","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"enviar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/send","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"cancelar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/cancel","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"prever","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/preview","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"ler_previa","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["organizacaoId","previaId"],"verbo":"GET","caminho":"/broadcast-previews/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"formato":"texto"}},
{"classe":"Campanhas","metodo":"testar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Campanhas","metodo":"obter_metricas","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/metrics","kwargs":{}},
{"classe":"Campanhas","metodo":"listar_destinatarios","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["id","filtros"],"verbo":"GET","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/recipients","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Campanhas","metodo":"listar_links_clicados","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["id","paginacao"],"verbo":"GET","caminho":"/broadcasts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/clicked-links","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"ChavesApi","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["dados"],"verbo":"POST","caminho":"/api-keys","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"ChavesApi","metodo":"listar","argumentos":[],"parametros":[],"verbo":"GET","caminho":"/api-keys","kwargs":{}},
{"classe":"ChavesApi","metodo":"renomear","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","dados"],"verbo":"PATCH","caminho":"/api-keys/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"ChavesApi","metodo":"revogar","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/api-keys/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"ConcessoesCli","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["paginacao"],"verbo":"GET","caminho":"/cli/oauth/grants","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"ConcessoesCli","metodo":"revogar","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"POST","caminho":"/cli/oauth/grants/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/revoke","kwargs":{}},
{"classe":"SegmentosContato","metodo":"listar","argumentos":["parte-0/ação+?@example.test"],"parametros":["contatoId"],"verbo":"GET","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/segments","kwargs":{}},
{"classe":"SegmentosContato","metodo":"adicionar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["contatoId","segmentoId"],"verbo":"POST","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/segments","kwargs":{"aceitar_vazio":true,"corpo":{"segmento_id":"parte-1/ação+?@example.test"}}},
{"classe":"SegmentosContato","metodo":"remover","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["contatoId","segmentoId"],"verbo":"DELETE","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/segments/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"TopicosContato","metodo":"listar","argumentos":["parte-0/ação+?@example.test"],"parametros":["contatoId"],"verbo":"GET","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/topics","kwargs":{}},
{"classe":"TopicosContato","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",[{"topico_id":"topico-a","assinatura":"opt_out"}]],"parametros":["contatoId","topicos"],"verbo":"PATCH","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/topics","kwargs":{"corpo":{"topicos":[{"topico_id":"topico-a","assinatura":"opt_out"}]}}},
{"classe":"Contatos","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/contacts","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Contatos","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["filtros"],"verbo":"GET","caminho":"/contacts","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Contatos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["idOuEmail"],"verbo":"GET","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Contatos","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["idOuEmail","pedido"],"verbo":"PATCH","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Contatos","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["idOuEmail"],"verbo":"DELETE","caminho":"/contacts/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Contatos","metodo":"adicionar_em_lote","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/contacts/batch","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Contatos","metodo":"alterar_em_lote","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/contacts/bulk","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"ExportacoesContatos","metodo":"criar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["filtros"],"verbo":"POST","caminho":"/contacts/exports","kwargs":{"corpo":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"ExportacoesContatos","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["consulta"],"verbo":"GET","caminho":"/contacts/exports","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"ExportacoesContatos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/contacts/exports/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"ExportacoesContatos","metodo":"baixar","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/contacts/exports/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/download","kwargs":{"formato":"bytes"}},
{"classe":"ExportacoesContatos","metodo":"baixar_fluxo","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/contacts/exports/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/download","kwargs":{"formato":"fluxo"}},
{"classe":"ImportacoesContatos","metodo":"prever","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/contacts/imports/preview","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"ImportacoesContatos","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/contacts/imports","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"ImportacoesContatos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/contacts/imports/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Segmentos","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/segments","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Segmentos","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["paginacao"],"verbo":"GET","caminho":"/segments","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Segmentos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/segments/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Segmentos","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"PATCH","caminho":"/segments/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Segmentos","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/segments/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Topicos","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/topics","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Topicos","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["paginacao"],"verbo":"GET","caminho":"/topics","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Topicos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/topics/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Topicos","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"PATCH","caminho":"/topics/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Topicos","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/topics/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"PropriedadesContatos","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/contact-properties","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"PropriedadesContatos","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["paginacao"],"verbo":"GET","caminho":"/contact-properties","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"PropriedadesContatos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/contact-properties/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"PropriedadesContatos","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"PATCH","caminho":"/contact-properties/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"PropriedadesContatos","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/contact-properties/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Dominios","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["dados"],"verbo":"POST","caminho":"/domains","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Dominios","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["consulta"],"verbo":"GET","caminho":"/domains","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Dominios","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/domains/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Dominios","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","dados"],"verbo":"PATCH","caminho":"/domains/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Dominios","metodo":"registros_dns","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/domains/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/dns","kwargs":{}},
{"classe":"Dominios","metodo":"verificar","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"POST","caminho":"/domains/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/verify","kwargs":{}},
{"classe":"Dominios","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/domains/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Emails","metodo":"enviar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["dados"],"verbo":"POST","caminho":"/emails","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Emails","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["consulta"],"verbo":"GET","caminho":"/emails","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Emails","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/emails/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Emails","metodo":"reagendar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","dados"],"verbo":"PATCH","caminho":"/emails/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Emails","metodo":"cancelar","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"POST","caminho":"/emails/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/cancel","kwargs":{}},
{"classe":"Emails","metodo":"obter_por_idempotencia","argumentos":["parte-0/ação+?@example.test"],"parametros":["chave"],"verbo":"GET","caminho":"/emails/idempotency/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Emails","metodo":"eventos","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/emails/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/events","kwargs":{}},
{"classe":"Emails","metodo":"baixar_anexo","argumentos":["parte-0/ação+?@example.test",2],"parametros":["id","indice"],"verbo":"GET","caminho":"/emails/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/attachments/2","kwargs":{"formato":"bytes"}},
{"classe":"Lotes","metodo":"enviar","argumentos":[[{"texto":"ação","opcional":null,"ativo":false}]],"parametros":["dados"],"verbo":"POST","caminho":"/emails/batch","kwargs":{"corpo":[{"texto":"ação","opcional":null,"ativo":false}]}},
{"classe":"Lotes","metodo":"obter_por_idempotencia","argumentos":["parte-0/ação+?@example.test"],"parametros":["chave"],"verbo":"GET","caminho":"/emails/batch/idempotency/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Eventos","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/events","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Eventos","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["paginacao"],"verbo":"GET","caminho":"/events","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Eventos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["idOuNome"],"verbo":"GET","caminho":"/events/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Eventos","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["idOuNome","pedido"],"verbo":"PATCH","caminho":"/events/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Eventos","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["idOuNome"],"verbo":"DELETE","caminho":"/events/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Eventos","metodo":"enviar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/events/send","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"VersoesModelo","metodo":"listar","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["modeloId","consulta"],"verbo":"GET","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/versions","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"VersoesModelo","metodo":"obter","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["modeloId","versaoId"],"verbo":"GET","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/versions/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"VersoesModelo","metodo":"restaurar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["modeloId","versaoId","pedido"],"verbo":"POST","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/versions/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/restore","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"VersoesModelo","metodo":"prever","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["modeloId","versaoId","pedido"],"verbo":"POST","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/versions/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/preview","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"PastasModelo","metodo":"listar","argumentos":[],"parametros":[],"verbo":"GET","caminho":"/template-folders","kwargs":{}},
{"classe":"PastasModelo","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/template-folders","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"PastasModelo","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"PATCH","caminho":"/template-folders/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"PastasModelo","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/template-folders/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Modelos","metodo":"prever","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/preview","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Modelos","metodo":"obter_previa","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["organizacaoId","previaId"],"verbo":"GET","caminho":"/template-previews/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"formato":"texto"}},
{"classe":"Modelos","metodo":"testar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"POST","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Modelos","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/templates","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Modelos","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["filtros"],"verbo":"GET","caminho":"/templates","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Modelos","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["idOuAlias"],"verbo":"GET","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Modelos","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["idOuAlias","pedido"],"verbo":"PATCH","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Modelos","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["idOuAlias"],"verbo":"DELETE","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Modelos","metodo":"publicar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["idOuAlias","pedido"],"verbo":"POST","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/publish","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Modelos","metodo":"duplicar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["idOuAlias","pedido"],"verbo":"POST","caminho":"/templates/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/duplicate","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Recebimento","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["filtros"],"verbo":"GET","caminho":"/receiving","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Recebimento","metodo":"obter","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["id","formato"],"verbo":"GET","caminho":"/receiving/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Recebimento","metodo":"listar_anexos","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/receiving/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/attachments","kwargs":{}},
{"classe":"Recebimento","metodo":"obter_anexo","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["id","anexoId"],"verbo":"GET","caminho":"/receiving/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/attachments/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Recebimento","metodo":"baixar_original","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/receiving/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/raw","kwargs":{"formato":"bytes"}},
{"classe":"Recebimento","metodo":"baixar_anexo","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["id","anexoId"],"verbo":"GET","caminho":"/receiving/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/attachments/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/download","kwargs":{"formato":"bytes"}},
{"classe":"Recebimento","metodo":"baixar_original_fluxo","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/receiving/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/raw","kwargs":{"formato":"fluxo"}},
{"classe":"Recebimento","metodo":"baixar_anexo_fluxo","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["id","anexoId"],"verbo":"GET","caminho":"/receiving/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/attachments/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/download","kwargs":{"formato":"fluxo"}},
{"classe":"EventosWebhook","metodo":"listar","argumentos":["parte-0/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["webhookId","paginacao"],"verbo":"GET","caminho":"/webhooks/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/events","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"EventosWebhook","metodo":"obter","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["webhookId","eventoId"],"verbo":"GET","caminho":"/webhooks/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/events/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"EventosWebhook","metodo":"listar_tentativas","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test",{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["webhookId","eventoId","paginacao"],"verbo":"GET","caminho":"/webhooks/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/events/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/attempts","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"EventosWebhook","metodo":"reenviar","argumentos":["parte-0/ação+?@example.test","parte-1/ação+?@example.test"],"parametros":["webhookId","eventoId"],"verbo":"POST","caminho":"/webhooks/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/events/parte-1%2Fa%C3%A7%C3%A3o%2B%3F%40example.test/replay","kwargs":{}},
{"classe":"Webhooks","metodo":"criar","argumentos":[{"texto":"ação","opcional":null,"ativo":false}],"parametros":["pedido"],"verbo":"POST","caminho":"/webhooks","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Webhooks","metodo":"listar","argumentos":[{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}],"parametros":["paginacao"],"verbo":"GET","caminho":"/webhooks","kwargs":{"consulta":{"limite":2,"antes":"cursor/ação","ativo":false,"opcional":null}}},
{"classe":"Webhooks","metodo":"obter","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"GET","caminho":"/webhooks/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}},
{"classe":"Webhooks","metodo":"atualizar","argumentos":["parte-0/ação+?@example.test",{"texto":"ação","opcional":null,"ativo":false}],"parametros":["id","pedido"],"verbo":"PATCH","caminho":"/webhooks/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{"corpo":{"texto":"ação","opcional":null,"ativo":false}}},
{"classe":"Webhooks","metodo":"remover","argumentos":["parte-0/ação+?@example.test"],"parametros":["id"],"verbo":"DELETE","caminho":"/webhooks/parte-0%2Fa%C3%A7%C3%A3o%2B%3F%40example.test","kwargs":{}}
]`

type casoRecurso struct {
	Classe       string            `json:"classe"`
	Metodo       string            `json:"metodo"`
	Argumentos   []json.RawMessage `json:"argumentos"`
	Verbo        string            `json:"verbo"`
	Caminho      string            `json:"caminho"`
	Configuracao struct {
		Corpo        json.RawMessage `json:"corpo"`
		Consulta     Consulta        `json:"consulta"`
		Formato      string          `json:"formato"`
		AceitarVazio bool            `json:"aceitar_vazio"`
	} `json:"kwargs"`
}

type chamadaRecurso struct {
	contexto     context.Context
	metodo       string
	caminho      string
	corpo        []byte
	consulta     Consulta
	formato      string
	aceitarVazio bool
	opcoes       OpcoesRequisicao
	consumidor   bool
}

type transporteColetorRecursos struct {
	chamadas     []chamadaRecurso
	retorno      *Resposta
	retornoFluxo *Resposta
}

func novoColetorRecursos() *transporteColetorRecursos {
	return &transporteColetorRecursos{
		retorno:      &Resposta{Dados: json.RawMessage(`{"id":"sintetico"}`), Status: 200, Cabecalhos: http.Header{}},
		retornoFluxo: &Resposta{Status: 200, Cabecalhos: http.Header{}, TotalBytes: 6},
	}
}

func (t *transporteColetorRecursos) executar(ctx context.Context, metodo, caminho string, corpo any, consulta Consulta, formato string, aceitarVazio bool, opcoes OpcoesRequisicao, consumidor func([]byte) error) (*Resposta, error) {
	var dados []byte
	var err error
	if corpo != nil {
		// O core mantém wrappers de shape até a serialização. O coletor registra
		// somente pedidos que estariam prontos para a camada de rede.
		dados, err = codificarCorpo(corpo)
		if err != nil {
			return nil, err
		}
	}
	opcoes.Cabecalhos = opcoes.Cabecalhos.Clone()
	t.chamadas = append(t.chamadas, chamadaRecurso{ctx, metodo, caminho, dados, consulta, formato, aceitarVazio, opcoes, consumidor != nil})
	if formato == "fluxo" {
		for _, parte := range [][]byte{[]byte("abc"), []byte("def")} {
			if err := consumidor(parte); err != nil {
				return nil, err
			}
		}
		return t.retornoFluxo, nil
	}
	return t.retorno, nil
}

var construtoresRecursos = map[string]func(executor) any{
	"Emails":                 func(t executor) any { return novoEmails(t) },
	"Lotes":                  func(t executor) any { return novoLotes(t) },
	"Dominios":               func(t executor) any { return novoDominios(t) },
	"ChavesApi":              func(t executor) any { return novoChavesAPI(t) },
	"Contatos":               func(t executor) any { return novoContatos(t) },
	"SegmentosContato":       func(t executor) any { return novoSegmentosContato(t) },
	"TopicosContato":         func(t executor) any { return novoTopicosContato(t) },
	"ImportacoesContatos":    func(t executor) any { return novoImportacoesContatos(t) },
	"ExportacoesContatos":    func(t executor) any { return novoExportacoesContatos(t) },
	"Segmentos":              func(t executor) any { return novoSegmentos(t) },
	"Topicos":                func(t executor) any { return novoTopicos(t) },
	"PropriedadesContatos":   func(t executor) any { return novoPropriedadesContatos(t) },
	"Modelos":                func(t executor) any { return novoModelos(t) },
	"VersoesModelo":          func(t executor) any { return novoVersoesModelo(t) },
	"PastasModelo":           func(t executor) any { return novoPastasModelo(t) },
	"Webhooks":               func(t executor) any { return novoWebhooks(t) },
	"EventosWebhook":         func(t executor) any { return novoEventosWebhook(t) },
	"Campanhas":              func(t executor) any { return novoCampanhas(t) },
	"Automacoes":             func(t executor) any { return novoAutomacoes(t) },
	"ExecucoesAutomacoes":    func(t executor) any { return novoExecucoesAutomacoes(t) },
	"Eventos":                func(t executor) any { return novoEventos(t) },
	"Recebimento":            func(t executor) any { return novoRecebimento(t) },
	"Caixas":                 func(t executor) any { return novoCaixas(t) },
	"ConversasCaixa":         func(t executor) any { return novoConversasCaixa(t) },
	"EtiquetasCaixa":         func(t executor) any { return novoEtiquetasCaixa(t) },
	"EtiquetasConversaCaixa": func(t executor) any { return novoEtiquetasConversaCaixa(t) },
	"EmailsCaixa":            func(t executor) any { return novoEmailsCaixa(t) },
	"RascunhosCaixa":         func(t executor) any { return novoRascunhosCaixa(t) },
	"CliMailcore":            func(t executor) any { return novoCLIMailcore(t) },
	"ConcessoesCli":          func(t executor) any { return novoConcessoesCLI(t) },
}

func casosRecursos(t *testing.T) []casoRecurso {
	t.Helper()
	var casos []casoRecurso
	decodificador := json.NewDecoder(strings.NewReader(contratosRecursos))
	decodificador.UseNumber()
	if err := decodificador.Decode(&casos); err != nil {
		t.Fatal(err)
	}
	if len(casos) != 151 {
		t.Fatalf("Contrato deveria ter 151 métodos, obteve %d", len(casos))
	}
	return casos
}

func metodoRecurso(nome string) string {
	if nome == "registros_dns" {
		return "RegistrosDNS"
	}
	partes := strings.Split(nome, "_")
	for i := range partes {
		partes[i] = strings.ToUpper(partes[i][:1]) + partes[i][1:]
	}
	return strings.Join(partes, "")
}

func argumentoRecurso(t *testing.T, dados json.RawMessage, tipo reflect.Type) reflect.Value {
	t.Helper()
	destino := reflect.New(tipo)
	decodificador := json.NewDecoder(bytes.NewReader(dados))
	decodificador.UseNumber()
	if err := decodificador.Decode(destino.Interface()); err != nil {
		t.Fatal(err)
	}
	return destino.Elem()
}

func jsonRecurso(t *testing.T, dados []byte) any {
	t.Helper()
	if len(dados) == 0 {
		return nil
	}
	var valor any
	decodificador := json.NewDecoder(bytes.NewReader(dados))
	decodificador.UseNumber()
	if err := decodificador.Decode(&valor); err != nil {
		t.Fatal(err)
	}
	return valor
}

func segmentosRecurso(t *testing.T, caminho string) []string {
	t.Helper()
	partes := strings.Split(caminho, "/")
	for i := range partes {
		var err error
		partes[i], err = url.PathUnescape(partes[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	return partes
}

type chaveContextoRecurso struct{}

func TestRecursosContrato151(t *testing.T) {
	for _, caso := range casosRecursos(t) {
		t.Run(caso.Classe+"."+metodoRecurso(caso.Metodo), func(t *testing.T) {
			transporte := novoColetorRecursos()
			recurso := construtoresRecursos[caso.Classe](transporte)
			metodo := reflect.ValueOf(recurso).MethodByName(metodoRecurso(caso.Metodo))
			if !metodo.IsValid() {
				t.Fatal("Método do contrato não encontrado")
			}
			ctx := context.WithValue(context.Background(), chaveContextoRecurso{}, caso.Classe+caso.Metodo)
			argumentos := []reflect.Value{reflect.ValueOf(ctx)}
			for i, dados := range caso.Argumentos {
				argumentos = append(argumentos, argumentoRecurso(t, dados, metodo.Type().In(i+1)))
			}
			var partes []string
			fluxo := caso.Configuracao.Formato == "fluxo"
			if fluxo {
				argumentos = append(argumentos, reflect.ValueOf(func(parte []byte) error { partes = append(partes, string(parte)); return nil }))
			}
			opcoes := OpcoesRequisicao{TempoLimite: 2500 * time.Millisecond, LimiteRespostaBytes: 4096, ChaveIdempotencia: "negocio-123", Cabecalhos: http.Header{"X-Correlacao": []string{"sintetica"}}}
			argumentos = append(argumentos, reflect.ValueOf(opcoes))
			resultados := metodo.Call(argumentos)
			if !resultados[1].IsNil() {
				t.Fatal(resultados[1].Interface().(error))
			}
			esperadoRetorno := transporte.retorno
			if fluxo {
				esperadoRetorno = transporte.retornoFluxo
			}
			if resultados[0].Interface().(*Resposta) != esperadoRetorno {
				t.Fatal("Recurso alterou a resposta do transporte")
			}
			if len(transporte.chamadas) != 1 {
				t.Fatalf("Esperava uma chamada, obteve %d", len(transporte.chamadas))
			}
			chamada := transporte.chamadas[0]
			if chamada.contexto != ctx {
				t.Fatal("O contexto do chamador não foi preservado")
			}
			if chamada.metodo != caso.Verbo {
				t.Fatalf("Método HTTP: %s != %s", chamada.metodo, caso.Verbo)
			}
			if !reflect.DeepEqual(segmentosRecurso(t, chamada.caminho), segmentosRecurso(t, caso.Caminho)) {
				t.Fatalf("Caminho incorreto: %q", chamada.caminho)
			}
			if !reflect.DeepEqual(jsonRecurso(t, chamada.corpo), jsonRecurso(t, caso.Configuracao.Corpo)) {
				t.Fatalf("Corpo divergente: %s != %s", chamada.corpo, caso.Configuracao.Corpo)
			}
			if len(chamada.consulta) != 0 || len(caso.Configuracao.Consulta) != 0 {
				if !reflect.DeepEqual(chamada.consulta, caso.Configuracao.Consulta) {
					t.Fatalf("Consulta divergente: %#v", chamada.consulta)
				}
			}
			formato := caso.Configuracao.Formato
			if formato == "" {
				formato = "json"
			}
			if chamada.formato != formato || chamada.aceitarVazio != caso.Configuracao.AceitarVazio || chamada.consumidor != fluxo {
				t.Fatal("Formato, resposta vazia ou callback divergentes")
			}
			if !reflect.DeepEqual(chamada.opcoes, opcoes) {
				t.Fatal("Opções da requisição foram alteradas")
			}
			if fluxo && !reflect.DeepEqual(partes, []string{"abc", "def"}) {
				t.Fatal("O consumidor não recebeu os blocos")
			}
		})
	}
}

func TestRecursosSuperficiePublica(t *testing.T) {
	porClasse := map[string]map[string]bool{}
	for _, caso := range casosRecursos(t) {
		if porClasse[caso.Classe] == nil {
			porClasse[caso.Classe] = map[string]bool{}
		}
		porClasse[caso.Classe][metodoRecurso(caso.Metodo)] = true
	}
	if len(construtoresRecursos) != 30 {
		t.Fatal("São esperados 30 recursos construídos")
	}
	for nome, construir := range construtoresRecursos {
		tipo := reflect.TypeOf(construir(novoColetorRecursos()))
		if tipo.NumMethod() != len(porClasse[nome]) {
			t.Errorf("%s possui %d métodos, contrato contém %d", nome, tipo.NumMethod(), len(porClasse[nome]))
		}
		for i := 0; i < tipo.NumMethod(); i++ {
			if !porClasse[nome][tipo.Method(i).Name] {
				t.Errorf("Método público fora do contrato: %s.%s", nome, tipo.Method(i).Name)
			}
		}
	}
}

func TestRecursosFilhosCompartilhamTransporte(t *testing.T) {
	transporte := novoColetorRecursos()
	contatos := novoContatos(transporte)
	contatos.Segmentos.Listar(context.Background(), "contato")
	contatos.Topicos.Listar(context.Background(), "contato")
	contatos.Importacoes.Obter(context.Background(), "importacao")
	contatos.Exportacoes.Obter(context.Background(), "exportacao")
	modelos := novoModelos(transporte)
	modelos.Versoes.Obter(context.Background(), "modelo", "versao")
	modelos.Pastas.Listar(context.Background())
	novoWebhooks(transporte).Eventos.Obter(context.Background(), "webhook", "evento")
	novoAutomacoes(transporte).Execucoes.Obter(context.Background(), "automacao", "execucao")
	caixas := novoCaixas(transporte)
	caixas.Conversas.Listar(context.Background(), "caixa", nil)
	caixas.Conversas.Etiquetas.Substituir(context.Background(), "caixa", "conversa", nil)
	caixas.Conversas.Emails.Obter(context.Background(), "caixa", "conversa", "email")
	caixas.Etiquetas.Listar(context.Background(), "caixa")
	caixas.Rascunhos.Obter(context.Background(), "caixa", "rascunho")
	novoCLIMailcore(transporte).Concessoes.Listar(context.Background(), nil)
	if len(transporte.chamadas) != 14 {
		t.Fatalf("Filhos não compartilharam o transporte: %d chamadas", len(transporte.chamadas))
	}
	if string(transporte.chamadas[9].corpo) != `{"ids":[]}` {
		t.Fatalf("Lista vazia de etiquetas não foi preservada: %s", transporte.chamadas[9].corpo)
	}
}

func TestRecursosNilCorposObjetosEListas(t *testing.T) {
	transporte := novoColetorRecursos()
	ctx := context.Background()
	var pedido *PedidoEmail
	var documento Documento
	operacoes := []func() (*Resposta, error){
		func() (*Resposta, error) { return novoEmails(transporte).Enviar(ctx, nil) },
		func() (*Resposta, error) { return novoEmails(transporte).Enviar(ctx, pedido) },
		func() (*Resposta, error) { return novoEmails(transporte).Enviar(ctx, documento) },
		func() (*Resposta, error) { return novoModelos(transporte).Duplicar(ctx, "modelo", nil) },
		func() (*Resposta, error) { return novoVersoesModelo(transporte).Prever(ctx, "modelo", "versao", nil) },
		func() (*Resposta, error) { return novoExportacoesContatos(transporte).Criar(ctx, nil) },
		func() (*Resposta, error) { return novoRascunhosCaixa(transporte).Criar(ctx, "caixa", nil) },
	}
	for _, operacao := range operacoes {
		if _, err := operacao(); err != nil {
			t.Fatal(err)
		}
		if string(transporte.chamadas[len(transporte.chamadas)-1].corpo) != "{}" {
			t.Fatal("Corpo opcional deveria ser objeto vazio")
		}
	}
	if _, err := novoTopicosContato(transporte).Atualizar(ctx, "contato", nil); err != nil {
		t.Fatal(err)
	}
	if string(transporte.chamadas[len(transporte.chamadas)-1].corpo) != `{"topicos":[]}` {
		t.Fatal("Tópicos nil devem manter lista vazia explícita")
	}
}

func TestRecursosPreservamStructsMapasNullFalseEListas(t *testing.T) {
	transporte := novoColetorRecursos()
	dados := Documento{"evento": "compra", "email": "usuario@example.test", "dados": Documento{
		"nulo": nil, "ativo": false, "objetoVazio": Documento{}, "listaVazia": []string{},
		"mapaNumerico": Documento{"0": "email", "1": "primeiro_nome"}, "campoNulo": (Documento)(nil),
	}}
	antes, _ := json.Marshal(dados)
	if _, err := novoEventos(transporte).Enviar(context.Background(), dados); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(antes, transporte.chamadas[0].corpo) {
		t.Fatal("Valores JSON foram modificados")
	}
	depois, _ := json.Marshal(dados)
	if !bytes.Equal(antes, depois) {
		t.Fatal("Documento do chamador foi alterado")
	}
	tipado := &PedidoEmail{Remetente: "origem@example.test", Destinatarios: []string{"destino@example.test"}, Assunto: "Ação", HTML: "<p>Olá</p>"}
	if _, err := novoEmails(transporte).Enviar(context.Background(), tipado); err != nil {
		t.Fatal(err)
	}
	esperado, _ := json.Marshal(tipado)
	if !bytes.Equal(esperado, transporte.chamadas[1].corpo) {
		t.Fatal("Pedido tipado não manteve suas tags JSON")
	}
}

func conferirErroRecurso(t *testing.T, err error) {
	t.Helper()
	var erro *Erro
	if !errors.As(err, &erro) || erro.Codigo != "configuracao_invalida" {
		t.Fatalf("Esperava erro de configuração, recebeu %v", err)
	}
}

type listaJSONObjeto []int

func (listaJSONObjeto) MarshalJSON() ([]byte, error) { return []byte(`{}`), nil }

func TestRecursosValidamConfiguracaoAntesDaRede(t *testing.T) {
	transporte := novoColetorRecursos()
	ctx := context.Background()
	emails := novoEmails(transporte)
	for _, id := range []string{"", ".", "..", "linha\nquebra", "controle\x00", string([]byte{0xff})} {
		_, err := emails.Obter(ctx, id)
		conferirErroRecurso(t, err)
	}
	_, err := emails.BaixarAnexo(ctx, "email", -1)
	conferirErroRecurso(t, err)
	_, err = emails.Obter(ctx, "email", OpcoesRequisicao{}, OpcoesRequisicao{})
	conferirErroRecurso(t, err)
	var nulo *Emails
	_, err = nulo.Obter(ctx, "email")
	conferirErroRecurso(t, err)
	_, err = novoEmails(nil).Obter(ctx, "email")
	conferirErroRecurso(t, err)
	_, err = novoRecebimento(transporte).BaixarOriginalFluxo(ctx, "email", nil)
	conferirErroRecurso(t, err)
	_, err = novoEmails(transporte).Enviar(ctx, []string{"objeto obrigatório"})
	conferirErroRecurso(t, err)
	for _, dados := range []any{nil, ([]Documento)(nil), Documento{"nao": "lista"}, []byte{1, 2}, listaJSONObjeto{1}, json.RawMessage(`{}`)} {
		_, err = novoLotes(transporte).Enviar(ctx, dados)
		conferirErroRecurso(t, err)
	}
	if len(transporte.chamadas) != 0 {
		t.Fatal("Entrada inválida iniciou pedido")
	}
}

func TestRecursosFluxoPreservaErroOriginalDoConsumidor(t *testing.T) {
	transporte := novoColetorRecursos()
	original := errors.New("cancelamento sintético")
	vezes := 0
	_, err := novoRecebimento(transporte).BaixarOriginalFluxo(context.Background(), "email", func([]byte) error { vezes++; return original })
	if err != original || vezes != 1 || len(transporte.chamadas) != 1 {
		t.Fatal("Recurso repetiu ou mascarou o erro do consumidor")
	}
}
