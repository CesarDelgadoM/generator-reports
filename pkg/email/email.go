package email

import (
	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"gopkg.in/gomail.v2"
)

type Email struct {
	config *config.Config
	dialer *gomail.Dialer
}

func NewEmailSMTP(config *config.Config) *Email {
	// Config server smtp
	host := config.SMTP.Gmail.Client
	port := config.SMTP.Gmail.Port
	from := config.SMTP.Gmail.Email
	password := config.SMTP.Gmail.Password

	return &Email{
		config: config,
		dialer: gomail.NewDialer(host, port, from, password),
	}
}

func (e *Email) SendEmailWithAttachments(email string, attachment string, subject string, body string) bool {

	// Message
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.SMTP.Gmail.Email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.Attach(attachment)

	// Send email
	err := e.dialer.DialAndSend(m)
	if err != nil {
		zap.Log.Error("Error: ", err)
		return false
	}

	return true
}

func (e *Email) SendEmailNotification(email string, subject string, body string) bool {
	// Message
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.SMTP.Gmail.Email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Send email
	err := e.dialer.DialAndSend(m)
	if err != nil {
		zap.Log.Error("Error: ", err)
		return false
	}

	return true
}
