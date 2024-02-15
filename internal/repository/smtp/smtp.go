package smtp

import (
	"bytes"
	"fmt"
	"github.com/mbocek/meet-go/internal/config"
	"html/template"
	"net/smtp"
)

type SMTPRepository struct {
	host     string
	port     string
	from     string
	user     string
	password string
}

func NewSMTP(smtp config.SMTP) *SMTPRepository {
	return &SMTPRepository{
		host:     smtp.Host,
		port:     smtp.Port,
		from:     smtp.From,
		user:     smtp.User,
		password: smtp.Password,
	}
}

func (s *SMTPRepository) Send(to, tmpl string, params any) error {
	// authentication
	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	t, err := template.New("test").Parse(tmpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, params)
	if err != nil {
		return err
	}

	// send
	err = smtp.SendMail(fmt.Sprintf("%s:%s", s.host, s.port), auth, s.from, []string{to}, body.Bytes())
	if err != nil {
		return err
	}

	return nil
}
