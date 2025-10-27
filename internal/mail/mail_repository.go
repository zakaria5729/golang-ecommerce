package mail

import (
	"context"

	m "github.com/easy-comerce/backend/internal/mail/model"
	"github.com/easy-comerce/backend/internal/mail/provider"
)

type MailRepository interface {
	SendEmail(ctx context.Context, req *m.SendEmailRequest) (*m.SendEmailResponse, error)
}

type mailRepository struct {
	client *provider.SmtpMailProvider
}

func NewMailRepository(client *provider.SmtpMailProvider) MailRepository {
	return &mailRepository{
		client: client,
	}
}

func (r *mailRepository) SendEmail(ctx context.Context, req *m.SendEmailRequest) (*m.SendEmailResponse, error) {
	return r.client.Send(ctx, req)
}
