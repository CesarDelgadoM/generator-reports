package email

import (
	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"gopkg.in/gomail.v2"
)

type Email struct {
	config *config.Config
}

func NewEmailSMTP(config *config.Config) *Email {
	return &Email{
		config: config,
	}
}

func (e *Email) SendEmailWithAttachments(email string, attachment string, subject string, body string) bool {

	// Config server smtp
	host := e.config.SMTP.Gmail.Client
	port := e.config.SMTP.Gmail.Port
	from := e.config.SMTP.Gmail.Email
	password := e.config.SMTP.Gmail.Password

	dealer := gomail.NewDialer(host, port, from, password)

	// Message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.Attach(attachment)

	// Send email
	err := dealer.DialAndSend(m)
	if err != nil {
		zap.Log.Error("Error: ", err)
		return false
	}

	return true
}

func (e *Email) SendEmailNotification(email string, subject string, body string) bool {
	return true
}
