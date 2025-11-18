package attribute_option

import (
	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AttributeOptionRepository interface {
	GetAllAttributeOptions(include []string, showDeleted *bool, attributeTypeID *uint, sortBy, sortOrder string) ([]AttributeOptionEntity, error)
	GetAttributeOptionByID(id uint, include []string, showDeleted *bool) (*AttributeOptionEntity, error)
	CreateAttributeOption(attributeOption *AttributeOptionEntity) (*AttributeOptionEntity, error)
	UpdateAttributeOption(attributeOption *AttributeOptionEntity) error
	DeleteAttributeOption(id uint) error
	UndoDeletedAttributeOption(id uint) error
	AttributeOptionExists(id uint, showDeleted *bool) (bool, error)
	AttributeOptionExistsByNameAndType(attributeOptionName string, attributeTypeID uint, excludeID ...uint) (bool, error)
}

type attributeOptionRepository struct {
	db *gorm.DB
}

func NewAttributeOptionRepository() AttributeOptionRepository {
	return &attributeOptionRepository{
		db: db.GetDB(),
	}
}

func (r *attributeOptionRepository) GetAllAttributeOptions(include []string, showDeleted *bool, attributeTypeID *uint, sortBy, sortOrder string) ([]AttributeOptionEntity, error) {
	var attributeOptions []AttributeOptionEntity
	query := r.db.Model(&AttributeOptionEntity{})

	if attributeTypeID != nil && *attributeTypeID > 0 {
		query = query.Where(c.AttributeOptionAttributeTypeID+" = ?", *attributeTypeID)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if utils.ContainsString(include, "attribute_type") {
		query = query.Preload("AttributeType")
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&attributeOptions).Error
	if err != nil {
		l.Error("❌ Failed to fetch attribute options", "method", "GetAllAttributeOptions", "error", err, "include", include, "attributeTypeID", attributeTypeID, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return attributeOptions, err
}

func (r *attributeOptionRepository) GetAttributeOptionByID(id uint, include []string, showDeleted *bool) (*AttributeOptionEntity, error) {
	var attributeOption AttributeOptionEntity
	query := r.db.Model(&AttributeOptionEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if utils.ContainsString(include, "attribute_type") {
		query = query.Preload("AttributeType")
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&attributeOption).Error; err != nil {
		l.Error("❌ Failed to fetch attribute option by ID", "method", "GetAttributeOptionByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &attributeOption, nil
}

func (r *attributeOptionRepository) CreateAttributeOption(attributeOption *AttributeOptionEntity) (*AttributeOptionEntity, error) {
	err := r.db.Create(attributeOption).Error
	if err != nil {
		l.Error("❌ Failed to create attribute option", "method", "CreateAttributeOption", "error", err, "attributeOption", attributeOption)
		return nil, err
	}

	return attributeOption, nil
}

func (r *attributeOptionRepository) UpdateAttributeOption(attributeOption *AttributeOptionEntity) error {
	err := r.db.Save(attributeOption).Error
	if err != nil {
		l.Error("❌ Failed to update attribute option", "method", "UpdateAttributeOption", "error", err, "attributeOption", attributeOption)
	}

	return err
}

func (r *attributeOptionRepository) DeleteAttributeOption(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&AttributeOptionEntity{}).Error
	if err != nil {
		l.Error("❌ Failed to delete attribute option", "method", "DeleteAttributeOption", "error", err, "id", id)
	}

	return err
}

func (r *attributeOptionRepository) UndoDeletedAttributeOption(id uint) error {
	var attributeOption AttributeOptionEntity
	err := r.db.Unscoped().Where(c.FieldID+" = ?", id).First(&attributeOption).Error
	if err != nil {
		l.Error("❌ Failed to find deleted attribute option", "method", "UndoDeletedAttributeOption", "error", err, "id", id)
		return err
	}

	attributeOption.DeletedAt = nil
	err = r.db.Unscoped().Save(&attributeOption).Error
	if err != nil {
		l.Error("❌ Failed to undo deleted attribute option", "method", "UndoDeletedAttributeOption", "error", err, "id", id)
	}

	return err
}

func (r *attributeOptionRepository) AttributeOptionExists(id uint, showDeleted *bool) (bool, error) {
	var attributeOption AttributeOptionEntity
	query := r.db.Model(&AttributeOptionEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&attributeOption).Error
	if err != nil {
		l.Error("❌ Failed to check if attribute option exists", "method", "AttributeOptionExists", "error", err, "id", id)
		return false, err
	}

	return attributeOption.ID != 0, nil
}

func (r *attributeOptionRepository) AttributeOptionExistsByNameAndType(attributeOptionName string, attributeTypeID uint, excludeID ...uint) (bool, error) {
	var attributeOption AttributeOptionEntity
	query := r.db.Model(&AttributeOptionEntity{}).Where(c.AttributeOptionAttributeOptionName+" = ? AND "+c.AttributeOptionAttributeTypeID+" = ?", attributeOptionName, attributeTypeID)

	if len(excludeID) > 0 {
		query = query.Where(c.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(c.FieldID).Take(&attributeOption).Error
	if err != nil {
		l.Error("❌ Failed to check if attribute option exists by name and type", "method", "AttributeOptionExistsByNameAndType", "error", err, "attributeOptionName", attributeOptionName, "attributeTypeID", attributeTypeID)
		return false, err
	}

	return attributeOption.ID != 0, nil
}
