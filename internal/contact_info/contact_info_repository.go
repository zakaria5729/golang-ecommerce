package contact_info

import (
	"errors"

	e "github.com/easy-comerce/backend/pkg/app_error"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"gorm.io/gorm"
)

type ContactInfoRepository interface {
	CreateContactInfo(contactInfo *ContactInfoEntity) error
	UpdateContactInfo(id uint, contactInfo *ContactInfoEntity) error
	GetContactInfoByEmail(email string) (*ContactInfoEntity, error)
}

type contactInfoRepository struct {
	db *gorm.DB
}

func NewContactInfoRepository(db *gorm.DB) ContactInfoRepository {
	return &contactInfoRepository{
		db: db,
	}
}

func (r *contactInfoRepository) CreateContactInfo(contactInfo *ContactInfoEntity) error {
	if err := r.db.Create(contactInfo).Error; err != nil {
		l.Logger.Error("❌ Failed to create contact info", "method", "CreateContactInfo", "error", err)
		return err
	}
	return nil
}

func (r *contactInfoRepository) UpdateContactInfo(id uint, contactInfo *ContactInfoEntity) error {
	if err := r.db.Where(c.FieldID+" = ?", id).Updates(contactInfo).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("❌ Failed to update contact info", "method", "UpdateContactInfo", "error", err)
			return e.WrapServerError("Failed to update contact info", err)
		}
		return err
	}
	return nil
}

func (r *contactInfoRepository) GetContactInfoByEmail(email string) (*ContactInfoEntity, error) {
	var contactInfo ContactInfoEntity
	if err := r.db.Where(c.ContactInfoEmail+" = ?", email).First(&contactInfo).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("❌ Failed to fetch contact info", "method", "GetContactInfoByEmail", "error", err)
			return nil, e.WrapServerError("Failed to fetch contact info by email", err)
		}
		return nil, err
	}
	return &contactInfo, nil
}
