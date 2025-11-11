package contact_info

import (
	"errors"

	e "github.com/easy-comerce/backend/pkg/app_error"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type ContactInfoRepository interface {
	CreateContactInfo(contactInfo *ContactInfoEntity) error
	UpdateContactInfo(id uint, contactInfo *ContactInfoEntity) error
	GetContactInfoByEmail(email string) (*ContactInfoEntity, error)
	GetContactInfoById(id uint) (*ContactInfoEntity, error)
	GetContactInfosPaginated(page int, pageSize int, sortBy string, sortOrder string, contactType *string, showMessage *bool, showDeleted *bool) (*[]ContactInfoEntity, int64, error)
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

func (r *contactInfoRepository) GetContactInfoById(id uint) (*ContactInfoEntity, error) {
	var contactInfo ContactInfoEntity
	if err := r.db.Where(c.FieldID+" = ?", id).First(&contactInfo).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("❌ Failed to fetch contact info", "method", "GetContactInfoById", "error", err)
			return nil, e.WrapServerError("Failed to fetch contact info by id", err)
		}
		return nil, err
	}

	return &contactInfo, nil
}

func (r *contactInfoRepository) GetContactInfosPaginated(page int, pageSize int, sortBy string, sortOrder string, contactType *string, showMessage *bool, showDeleted *bool) (*[]ContactInfoEntity, int64, error) {
	var contacts []ContactInfoEntity
	var total int64

	selectFields := []string{
		c.FieldID, c.FieldCreatedAt, c.FieldUpdatedAt, c.ContactInfoName,
		c.ContactInfoEmail, c.ContactInfoType, c.ContactInfoPhone,
	}
	if showMessage != nil && *showMessage {
		selectFields = append(selectFields, c.ContactInfoMessage)
	}

	query := r.db.Model(&ContactInfoEntity{}).Select(selectFields)
	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if contactType != nil && *contactType != "" {
		query = query.Where(c.ContactInfoType+" = ?", *contactType)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Count(&total).Error; err != nil {
		l.Logger.Error("❌ Failed to count pages", "method", "GetContactInfosPaginated", "error", err, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder, "contactType", contactType, "showDeleted", showDeleted)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&contacts).Error
	if err != nil {
		l.Logger.Error("❌ Failed to fetch pages paginated", "method", "GetContactInfosPaginated", "error", err, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder, "contactType", contactType, "showDeleted", showDeleted)
		return nil, 0, err
	}

	return &contacts, total, nil
}
