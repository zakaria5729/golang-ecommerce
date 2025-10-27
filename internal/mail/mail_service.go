package mail

import (
	"context"
	"fmt"

	m "github.com/easy-comerce/backend/internal/mail/model"
)

type MailService interface {
	SendEmail(ctx context.Context, to []string, subject, body string, isHTML bool) error
	SendEmailWithAttachments(ctx context.Context, to []string, subject, body string, isHTML bool, attachments []string) error
}

type mailService struct {
	repo MailRepository
}

func NewMailService(repo MailRepository) MailService {
	return &mailService{
		repo: repo,
	}
}

func (s *mailService) SendEmail(ctx context.Context, to []string, subject, body string, isHTML bool) error {
	if len(to) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	if subject == "" || body == "" {
		return fmt.Errorf("subject and body are required")
	}

	req := &m.SendEmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}

	_, err := s.repo.SendEmail(ctx, req)
	return err
}

func (s *mailService) SendEmailWithAttachments(ctx context.Context, to []string, subject, body string, isHTML bool, attachments []string) error {
	if len(to) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	if subject == "" || body == "" {
		return fmt.Errorf("subject and body are required")
	}

	req := &m.SendEmailRequest{
		To:          to,
		Subject:     subject,
		Body:        body,
		IsHTML:      isHTML,
		Attachments: attachments,
	}

	_, err := s.repo.SendEmail(ctx, req)
	return err
}
