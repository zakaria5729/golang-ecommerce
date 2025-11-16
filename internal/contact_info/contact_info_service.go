package contact_info

import (
	"context"
	"errors"
	"fmt"

	m "github.com/easy-comerce/backend/internal/contact_info/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	op "github.com/easy-comerce/backend/pkg/option"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
	"gorm.io/gorm"
)

type ContactInfoService interface {
	PostNewsLetter(ctx context.Context, email string) error
	CreateContactUs(ctx context.Context, req *m.CreateContactInfoRequest) error
	GetContactInfoById(ctx context.Context, id uint) (*ContactInfoEntity, error)
	GetContactInfosPaginated(ctx context.Context, pageStr string, pageSizeStr string, sortBy string, sortOrder string, contactType *string, showMessage *bool, showDeleted *bool) (*response.PaginatedResponse, error)
}

type contactInfoService struct {
	repo ContactInfoRepository
}

func NewContactInfoService(repo ContactInfoRepository) ContactInfoService {
	return &contactInfoService{
		repo: repo,
	}
}

func (s *contactInfoService) PostNewsLetter(ctx context.Context, email string) error {
	option := op.QueryOptions{}
	option.AddFilter(c.ContactInfoEmail+" = ?", email)
	contact, err := s.repo.GetSingleBy(ctx, &option)

	if contact == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		contactInfo := &ContactInfoEntity{
			Email: email,
			Name:  utils.ExtractNameFromEmail(email, true),
			Type:  m.TypeNewsLetter,
		}

		return s.repo.Create(ctx, contactInfo)
	}

	return nil
}

func (s *contactInfoService) CreateContactUs(ctx context.Context, req *m.CreateContactInfoRequest) error {
	contactInfo := &ContactInfoEntity{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
		Type:    m.TypeContactUs,
	}

	if validationErrors := validator.ValidatePhone(req.Phone, "phone"); len(validationErrors) == 0 {
		contactInfo.Phone = req.Phone
	}

	option := op.QueryOptions{}
	option.AddFilter(c.ContactInfoEmail+" = ?", req.Email)

	contact, err := s.repo.GetSingleBy(ctx, &option)
	if contact == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return s.repo.Create(ctx, contactInfo)
	}

	return s.repo.Update(ctx, contact.ID, contactInfo)
}

func (s *contactInfoService) GetContactInfoById(ctx context.Context, id uint) (*ContactInfoEntity, error) {
	return s.repo.GetSingleByID(ctx, id, nil)
}

func (s *contactInfoService) GetContactInfosPaginated(ctx context.Context, pageStr string, pageSizeStr string, sortBy string, sortOrder string, contactType *string, showMessage *bool, showDeleted *bool) (*response.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	contacts, total, err := s.repo.GetContactInfosPaginated(ctx, page, pageSize, sortBy, sortOrder, contactType, showMessage, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contact infos: %w", err)
	}

	return utils.BuildPaginatedResponse(contacts, total, page, pageSize), nil
}
