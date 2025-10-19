package mail

import (
	"os"
	"strconv"

	"github.com/easy-comerce/backend/internal/mail/provider"
)

func NewSmtpMailProviderFromEnv() (*provider.SmtpMailProvider, error) {
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		smtpServer = "smtp.gmail.com"
	}

	smtpPort := 587
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			smtpPort = p
		}
	}

	username := os.Getenv("SMTP_USERNAME")
	if username == "" {
		username = os.Getenv("EMAIL_FROM")
	}

	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("EMAIL_FROM")
	fromName := os.Getenv("EMAIL_FROM_NAME")

	if fromEmail == "" {
		fromEmail = username
	}

	if fromName == "" {
		fromName = "Easy-Commerce"
	}

	return provider.NewSmtpMailProvider(smtpServer, smtpPort, username, password, fromEmail, fromName), nil
}
