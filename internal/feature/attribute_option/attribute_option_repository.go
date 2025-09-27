package attribute_option

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AttributeOptionRepository struct {
	db *gorm.DB
}

func NewAttributeOptionRepository() *AttributeOptionRepository {
	return &AttributeOptionRepository{
		db: db.GetDB(),
	}
}

func (r *AttributeOptionRepository) GetAllAttributeOptions(include []string, showDeleted *bool, attributeTypeID *uint, sortBy, sortOrder string) ([]AttributeOption, error) {
	var attributeOptions []AttributeOption

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if attributeTypeID != nil && *attributeTypeID > 0 {
		query = query.Where(constants.AttributeOptionAttributeTypeID+" = ?", *attributeTypeID)
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
		logger.Logger.Error("Failed to fetch attribute options", "method", "GetAllAttributeOptions", "error", err, "include", include, "attributeTypeID", attributeTypeID, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return attributeOptions, err
}

func (r *AttributeOptionRepository) GetAttributeOptionByID(id uint, include []string, showDeleted *bool) (*AttributeOption, error) {
	var attributeOption AttributeOption

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if utils.ContainsString(include, "attribute_type") {
		query = query.Preload("AttributeType")
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&attributeOption).Error; err != nil {
		logger.Logger.Error("Failed to fetch attribute option by ID", "method", "GetAttributeOptionByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &attributeOption, nil
}

func (r *AttributeOptionRepository) CreateAttributeOption(attributeOption *AttributeOption) (*AttributeOption, error) {
	err := r.db.Create(attributeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to create attribute option", "method", "CreateAttributeOption", "error", err, "attributeOption", attributeOption)
		return nil, err
	}
	return attributeOption, nil
}

func (r *AttributeOptionRepository) UpdateAttributeOption(attributeOption *AttributeOption) error {
	err := r.db.Save(attributeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to update attribute option", "method", "UpdateAttributeOption", "error", err, "attributeOption", attributeOption)
	}
	return err
}

func (r *AttributeOptionRepository) DeleteAttributeOption(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&AttributeOption{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete attribute option", "method", "DeleteAttributeOption", "error", err, "id", id)
	}

	return err
}

func (r *AttributeOptionRepository) UndoDeletedAttributeOption(id uint) error {
	var attributeOption AttributeOption
	err := r.db.Unscoped().Where(constants.FieldID+" = ?", id).First(&attributeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to find deleted attribute option", "method", "UndoDeletedAttributeOption", "error", err, "id", id)
		return err
	}

	attributeOption.DeletedAt = nil
	err = r.db.Unscoped().Save(&attributeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to undo deleted attribute option", "method", "UndoDeletedAttributeOption", "error", err, "id", id)
	}

	return err
}

func (r *AttributeOptionRepository) AttributeOptionExists(id uint, showDeleted *bool) (bool, error) {
	var attributeOption AttributeOption
	query := r.db.Model(&AttributeOption{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(constants.FieldID).Take(&attributeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to check if attribute option exists", "method", "AttributeOptionExists", "error", err, "id", id)
		return false, err
	}

	return attributeOption.ID != 0, nil
}

func (r *AttributeOptionRepository) AttributeOptionExistsByNameAndType(attributeOptionName string, attributeTypeID uint, excludeID ...uint) (bool, error) {
	var attributeOption AttributeOption
	query := r.db.Model(&AttributeOption{}).Where(constants.AttributeOptionAttributeOptionName+" = ? AND "+constants.AttributeOptionAttributeTypeID+" = ?", attributeOptionName, attributeTypeID)

	if len(excludeID) > 0 {
		query = query.Where(constants.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(constants.FieldID).Take(&attributeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to check if attribute option exists by name and type", "method", "AttributeOptionExistsByNameAndType", "error", err, "attributeOptionName", attributeOptionName, "attributeTypeID", attributeTypeID)
		return false, err
	}
	return attributeOption.ID != 0, nil
}

func (r *AttributeOptionRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.AttributeOptionAttributeTypeID, constants.AttributeOptionAttributeOptionName, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
