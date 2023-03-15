package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type SMTP struct {
	Host     string `env:"SMTP_HOST,default=smtp.mailtrap.io"`
	Port     int    `env:"SMTP_PORT,default=25"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
	Sender   string `env:"SMTP_SENDER,default=Test <no-reply@testdomain.com>"`
}

type Mailer struct {
	client *mail.Client
	sender string
}

func New(smtp SMTP) (Mailer, error) {
	client, err := mail.NewClient(
		smtp.Host,
		mail.WithPort(smtp.Port),
		mail.WithUsername(smtp.Username),
		mail.WithPassword(smtp.Password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		return Mailer{}, fmt.Errorf("failed new mail client: %w", err)
	}

	mailer := Mailer{
		client: client,
		sender: smtp.Sender,
	}

	return mailer, nil
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed read template fs: %w", err)
	}

	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("failed execute template subject: %w", err)
	}

	plainBody := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return fmt.Errorf("failed execute template plain body: %w", err)
	}

	htmlBody := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return fmt.Errorf("failed execute template html body: %w", err)
	}

	msg := mail.NewMsg()

	if err := msg.To(recipient); err != nil {
		return fmt.Errorf("failed message to: %w", err)
	}

	if err := msg.From(m.sender); err != nil {
		return fmt.Errorf("failed message from: %w", err)
	}

	msg.Subject(subject.String())

	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.SetBodyString(mail.TypeTextHTML, htmlBody.String())

	err = m.client.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("failed dial and send: %w", err)
	}

	return nil
}
