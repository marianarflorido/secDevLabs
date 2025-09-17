package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

// SendRecoveryEmail envia o e-mail com o token de recuperação
func SendRecoveryEmail(to, token string) error {
	m := gomail.NewMessage()

	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST") // ex: smtp.gmail.com
	port := 587                    // geralmente 587, pode mudar

	recoveryLink := fmt.Sprintf("http://localhost:8080/reset?token=%s", token)

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Recuperação de Senha")
	m.SetBody("text/html", fmt.Sprintf(`
		<p>Olá,</p>
		<p>Você solicitou a recuperação de senha.</p>
		<p>Clique no link abaixo para redefinir sua senha:</p>
		<p><a href="%s">%s</a></p>
	`, recoveryLink, recoveryLink))

	d := gomail.NewDialer(host, port, from, pass)

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}
