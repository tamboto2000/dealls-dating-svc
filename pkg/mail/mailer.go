package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Email struct {
	From    string
	To      []string
	Cc      []string
	Subject string
	Body    string
}

func (e Email) build() []byte {
	from := fmt.Sprintf("From: %s", e.From)
	to := fmt.Sprintf("To: %s", strings.Join(e.To, ","))
	cc := fmt.Sprintf("Cc: %s", strings.Join(e.Cc, ","))
	subject := fmt.Sprintf("Subject: %s", e.Subject)

	body := strings.Join(
		[]string{
			from,
			to,
			cc,
			subject,
			"MIME-version: 1.0;",
			"Content-Type: text/html; charset=\"UTF-8\";\n",
			e.Body,
		},
		"\n",
	)

	return []byte(body)
}

type Mailer struct {
	senderName string
	username   string
	pwd        string
	host       string
	port       string
}

func NewMailer(senderName, username, pwd, host, port string) *Mailer {
	return &Mailer{
		senderName: senderName,
		username:   username,
		pwd:        pwd,
		host:       host,
		port:       port,
	}
}

func (m *Mailer) SendMail(email Email) error {
	email.From = m.senderName

	auth := smtp.PlainAuth("", m.username, m.pwd, m.host)
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	err := smtp.SendMail(addr, auth, m.username, email.To, email.build())

	return err
}

func SendMail(host, username, pass string, email Email) error {
	smtp.PlainAuth("", username, pass, host)

	return nil
}
