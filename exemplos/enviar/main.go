// Exemplo de execução no servidor; utiliza MAILCORE_API_KEY e endereços autorizados.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailcore-hq/mailcore-go"
)

func main() {
	if executar() != nil {
		// A mensagem de negócio pode conter dados privados; não a publique em logs.
		fmt.Fprintln(os.Stderr, "Não foi possível enviar a mensagem de exemplo.")
		os.Exit(1)
	}
}

func executar() error {
	remetente, destinatario := os.Getenv("MAILCORE_REMETENTE"), os.Getenv("MAILCORE_DESTINATARIO")
	chave := os.Getenv("MAILCORE_IDEMPOTENCIA")
	if remetente == "" || destinatario == "" || chave == "" {
		return fmt.Errorf("configure remetente, destinatário e idempotência")
	}
	cliente, err := mailcore.NovoCliente("")
	if err != nil {
		return err
	}
	defer cliente.Fechar()
	ctx, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	resposta, err := cliente.Emails.Enviar(ctx, mailcore.PedidoEmail{
		Remetente: remetente, Destinatarios: []string{destinatario},
		Assunto: "Olá!", HTML: "<strong>Sua integração está pronta.</strong>",
	}, mailcore.OpcoesRequisicao{ChaveIdempotencia: chave})
	if err != nil {
		return err
	}
	var aceita mailcore.RespostaEnvio
	if err := resposta.Decodificar(&aceita); err != nil {
		return err
	}
	if aceita.ID == "" {
		return fmt.Errorf("aceitação sem identificador")
	}
	fmt.Println("Intenção de envio aceita.")
	return nil
}
