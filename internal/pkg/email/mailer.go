package mailer

import (
	"os"

	"github.com/resend/resend-go/v2"
)

type (
	EmailConfig struct {
		ApiKey string
		From   string
	}

	Mailer struct {
		emailConfig *EmailConfig
		client      *resend.Client
		Body        string
		Error       error
	}
)

func New() Mailer {
	apiKey := os.Getenv("RESEND_API_KEY")
	emailConfig := &EmailConfig{
		ApiKey: apiKey,
		From:   os.Getenv("RESEND_FROM"),
	}

	return Mailer{
		emailConfig: emailConfig,
		client:      resend.NewClient(apiKey),
		Body:        "",
		Error:       nil,
	}
}

func (m Mailer) Send(toEmail, subject string) Mailer {
	params := &resend.SendEmailRequest{
		From:    m.emailConfig.From,
		To:      []string{toEmail},
		Subject: subject,
		Html:    m.Body,
	}

	_, err := m.client.Emails.Send(params)
	if err != nil {
		m.Error = err
		return m
	}

	return m
}

