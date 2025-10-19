package mail

import (
	"context"

	m "github.com/easy-comerce/backend/internal/mail/model"
	"github.com/easy-comerce/backend/internal/mail/provider"
)

type MailRepository struct {
	client *provider.SmtpMailProvider
}

func NewMailRepository(client *provider.SmtpMailProvider) *MailRepository {
	return &MailRepository{
		client: client,
	}
}

func (r *MailRepository) SendEmail(ctx context.Context, req *m.SendEmailRequest) (*m.SendEmailResponse, error) {
	return r.client.Send(ctx, req)
}
