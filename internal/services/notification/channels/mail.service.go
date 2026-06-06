package channels

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/bigdann09/notifications/internal/config"
	"gopkg.in/gomail.v2"
)

//go:embed templates/*.html
var templates embed.FS

type MailChannel struct {
	cfg *config.MailConfig
}

func NewMailChannel(cfg *config.MailConfig) *MailChannel {
	return &MailChannel{cfg: cfg}
}

func (channel *MailChannel) Type() string {
	return string(Email)
}

func (channel *MailChannel) ComposeTemplate(templateName string, data map[string]any) (string, error) {
	tmpl, err := template.ParseFS(templates, templateName)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}
	return result.String(), nil
}

func (channel *MailChannel) Send(payload SendNotificationPayload) error {
	if channel.cfg.Host == "" {
		return fmt.Errorf("mail host not configured")
	}
	if payload.Email == "" {
		return fmt.Errorf("recipient email is required")
	}

	body, err := channel.ComposeTemplate("templates/random.html", map[string]any{
		"Title":   payload.Title,
		"Message": payload.Message,
		"Type":    payload.Type,
	})
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", channel.cfg.From)
	mail.SetHeader("To", payload.Email)
	mail.SetHeader("Subject", payload.Title)
	mail.SetBody("text/html", body)

	d := gomail.NewDialer(channel.cfg.Host, channel.cfg.Port, channel.cfg.User, channel.cfg.Pass)
	return d.DialAndSend(mail)
}
