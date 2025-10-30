package contact_info

import (
	"errors"

	m "github.com/easy-comerce/backend/internal/contact_info/model"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
	"gorm.io/gorm"
)

type ContactInfoService interface {
	PostNewsLetter(email string) error
	CreateContactUs(req *m.CreateContactInfoRequest) error
}

type contactInfoService struct {
	repo ContactInfoRepository
}

func NewContactInfoService(repo ContactInfoRepository) ContactInfoService {
	return &contactInfoService{
		repo: repo,
	}
}

func (s *contactInfoService) PostNewsLetter(email string) error {
	contact, err := s.repo.GetContactInfoByEmail(email)
	if contact == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		contactInfo := &ContactInfoEntity{
			Email: email,
			Name:  utils.ExtractNameFromEmail(email, true),
			Type:  m.TypeNewsLetter,
		}
		return s.repo.CreateContactInfo(contactInfo)
	}
	return nil
}

func (s *contactInfoService) CreateContactUs(req *m.CreateContactInfoRequest) error {
	contactInfo := &ContactInfoEntity{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
		Type:    m.TypeContactUs,
	}

	if validationErrors := validator.ValidatePhone(req.Phone, "phone"); len(validationErrors) == 0 {
		contactInfo.Phone = req.Phone
	}

	contact, err := s.repo.GetContactInfoByEmail(req.Email)
	if contact == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return s.repo.CreateContactInfo(contactInfo)
	}

	return s.repo.UpdateContactInfo(contact.ID, contactInfo)
}
