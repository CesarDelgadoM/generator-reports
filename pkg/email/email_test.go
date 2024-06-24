package email

import (
	"testing"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
	"github.com/stretchr/testify/assert"
)

func TestSendEmailWithAttachments(t *testing.T) {
	// Arguments
	email := "cesardelgadom2019@gmail.com"
	subject := "Send Email Test - Go"
	body := "Hi, the test works!"
	file := "file.txt"

	// Config
	config := config.GetConfig("config-dev.yml")

	// Logger
	zap.InitLogger(config)

	smtp := NewEmailSMTP(config)

	result := smtp.SendEmailWithAttachments(email, file, subject, body)

	assert.True(t, result, "Email sent!")
}
