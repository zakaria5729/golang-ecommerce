package provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	m "github.com/easy-comerce/backend/internal/mail/model"
	"github.com/wneessen/go-mail"
)

type SmtpMailProvider struct {
	smtpServer string
	smtpPort   int
	username   string
	password   string
	fromEmail  string
	fromName   string
}

func NewSmtpMailProvider(smtpServer string, smtpPort int, username, password, fromEmail, fromName string) *SmtpMailProvider {
	return &SmtpMailProvider{
		smtpServer: smtpServer,
		smtpPort:   smtpPort,
		username:   username,
		password:   password,
		fromEmail:  fromEmail,
		fromName:   fromName,
	}
}

func (c *SmtpMailProvider) Send(ctx context.Context, req *m.SendEmailRequest) (*m.SendEmailResponse, error) {
	msg := mail.NewMsg()
	if err := msg.FromFormat(c.fromName, c.fromEmail); err != nil {
		return nil, err
	}

	for _, to := range req.To {
		if err := msg.To(to); err != nil {
			return nil, err
		}
	}

	for _, cc := range req.Cc {
		if err := msg.AddCc(cc); err != nil {
			return nil, err
		}
	}

	for _, bcc := range req.Bcc {
		if err := msg.AddBcc(bcc); err != nil {
			return nil, err
		}
	}

	msg.Subject(req.Subject)

	if req.IsHTML {
		msg.SetBodyString(mail.TypeTextHTML, req.Body)
	} else {
		msg.SetBodyString(mail.TypeTextPlain, req.Body)
	}

	for _, attachment := range req.Attachments {
		msg.AttachFile(attachment)
	}

	client, err := mail.NewClient(c.smtpServer, mail.WithPort(c.smtpPort), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(c.username), mail.WithPassword(c.password))
	if err != nil {
		return nil, err
	}

	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	messageID := fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), hex.EncodeToString(randomBytes), c.fromEmail)
	msg.SetGenHeader(mail.HeaderMessageID, messageID)

	if err := client.DialAndSend(msg); err != nil {
		return &m.SendEmailResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &m.SendEmailResponse{
		Success:   true,
		MessageID: messageID,
	}, nil
}
