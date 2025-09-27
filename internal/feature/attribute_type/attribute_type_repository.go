package attribute_type

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AttributeTypeRepository struct {
	db *gorm.DB
}

func NewAttributeTypeRepository() *AttributeTypeRepository {
	return &AttributeTypeRepository{
		db: db.GetDB(),
	}
}

func (r *AttributeTypeRepository) GetAllAttributeTypes(include []string, showDeleted *bool, sortBy, sortOrder string) ([]AttributeType, error) {
	var attributeTypes []AttributeType

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&attributeTypes).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch attribute types", "method", "GetAllAttributeTypes", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return attributeTypes, err
}

func (r *AttributeTypeRepository) GetAttributeTypeByID(id uint, include []string, showDeleted *bool) (*AttributeType, error) {
	var attributeType AttributeType

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&attributeType).Error; err != nil {
		logger.Logger.Error("Failed to fetch attribute type by ID", "method", "GetAttributeTypeByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &attributeType, nil
}

func (r *AttributeTypeRepository) CreateAttributeType(attributeType *AttributeType) (*AttributeType, error) {
	err := r.db.Create(attributeType).Error
	if err != nil {
		logger.Logger.Error("Failed to create attribute type", "method", "CreateAttributeType", "error", err, "attributeType", attributeType)
		return nil, err
	}
	return attributeType, nil
}

func (r *AttributeTypeRepository) UpdateAttributeType(attributeType *AttributeType) error {
	err := r.db.Save(attributeType).Error
	if err != nil {
		logger.Logger.Error("Failed to update attribute type", "method", "UpdateAttributeType", "error", err, "attributeType", attributeType)
	}
	return err
}

func (r *AttributeTypeRepository) DeleteAttributeType(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&AttributeType{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete attribute type", "method", "DeleteAttributeType", "error", err, "id", id)
	}

	return err
}

func (r *AttributeTypeRepository) UndoDeletedAttributeType(id uint) error {
	var attributeType AttributeType
	err := r.db.Unscoped().Where(constants.FieldID+" = ?", id).First(&attributeType).Error
	if err != nil {
		logger.Logger.Error("Failed to find deleted attribute type", "method", "UndoDeletedAttributeType", "error", err, "id", id)
		return err
	}

	attributeType.DeletedAt = nil
	err = r.db.Unscoped().Save(&attributeType).Error
	if err != nil {
		logger.Logger.Error("Failed to undo deleted attribute type", "method", "UndoDeletedAttributeType", "error", err, "id", id)
	}

	return err
}

func (r *AttributeTypeRepository) AttributeTypeExists(id uint, showDeleted *bool) (bool, error) {
	var attributeType AttributeType
	query := r.db.Model(&AttributeType{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(constants.FieldID).Take(&attributeType).Error
	if err != nil {
		logger.Logger.Error("Failed to check if attribute type exists", "method", "AttributeTypeExists", "error", err, "id", id)
		return false, err
	}

	return attributeType.ID != 0, nil
}

func (r *AttributeTypeRepository) AttributeTypeExistsByName(name string, excludeID ...uint) (bool, error) {
	var attributeType AttributeType
	query := r.db.Model(&AttributeType{}).Where(constants.AttributeTypeName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(constants.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(constants.FieldID).Take(&attributeType).Error
	if err != nil {
		logger.Logger.Error("Failed to check if attribute type exists by name", "method", "AttributeTypeExistsByName", "error", err, "name", name)
		return false, err
	}

	return attributeType.ID != 0, nil
}

func (r *AttributeTypeRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.AttributeTypeName, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
