package branch

import (
	"strings"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/pkg/email"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

type IEmailBranch interface {
	SendEmail(queuename string, file string, email string)
}

type EmailBranch struct {
	config *config.Config
	email  *email.Email
}

func NewEmailBranch(config *config.Config, email *email.Email) IEmailBranch {
	return &EmailBranch{
		config: config,
		email:  email,
	}
}

func (e *EmailBranch) SendEmail(queuename string, file string, email string) {
	path := e.config.Branch.Pdf.Path + file
	subject := e.config.Branch.Notification.Success.Subject + strings.Split(queuename, "-")[0]
	body := e.config.Branch.Notification.Success.Body

	if e.email.SendEmailWithAttachments(email, path, subject, body) {
		zap.Log.Info(queuename, " email sent!")

	} else {
		zap.Log.Info(queuename, " The sent notification success email failed")

		subject = e.config.Branch.Notification.Failed.Subject
		body = e.config.Branch.Notification.Failed.Body

		if !e.email.SendEmailNotification(email, subject, body) {
			zap.Log.Info(queuename, " The sent notification failed email failed")
		}
	}
}
