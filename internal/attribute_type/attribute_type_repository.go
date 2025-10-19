package attribute_type

import (
	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
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

func (r *AttributeTypeRepository) GetAllAttributeTypes(showDeleted *bool, sortBy, sortOrder string) ([]AttributeTypeEntity, error) {
	var attributeTypes []AttributeTypeEntity
	query := r.db.Model(&AttributeTypeEntity{})

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&attributeTypes).Error
	if err != nil {
		l.Logger.Error("❌ Failed to fetch attribute types", "method", "GetAllAttributeTypes", "error", err, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return attributeTypes, err
}

func (r *AttributeTypeRepository) GetAttributeTypeByID(id uint, showDeleted *bool) (*AttributeTypeEntity, error) {
	var attributeType AttributeTypeEntity
	query := r.db.Model(&AttributeTypeEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&attributeType).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch attribute type by ID", "method", "GetAttributeTypeByID", "error", err, "id", id)
		return nil, err
	}

	return &attributeType, nil
}

func (r *AttributeTypeRepository) CreateAttributeType(attributeType *AttributeTypeEntity) (*AttributeTypeEntity, error) {
	err := r.db.Create(attributeType).Error
	if err != nil {
		l.Logger.Error("❌ Failed to create attribute type", "method", "CreateAttributeType", "error", err, "attributeType", attributeType)
		return nil, err
	}

	return attributeType, nil
}

func (r *AttributeTypeRepository) UpdateAttributeType(attributeType *AttributeTypeEntity) error {
	err := r.db.Save(attributeType).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update attribute type", "method", "UpdateAttributeType", "error", err, "attributeType", attributeType)
	}

	return err
}

func (r *AttributeTypeRepository) DeleteAttributeType(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&AttributeTypeEntity{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete attribute type", "method", "DeleteAttributeType", "error", err, "id", id)
	}

	return err
}

func (r *AttributeTypeRepository) UndoDeletedAttributeType(id uint) error {
	var attributeType AttributeTypeEntity
	err := r.db.Unscoped().Where(c.FieldID+" = ?", id).First(&attributeType).Error

	if err != nil {
		l.Logger.Error("❌ Failed to find deleted attribute type", "method", "UndoDeletedAttributeType", "error", err, "id", id)
		return err
	}

	attributeType.DeletedAt = nil
	err = r.db.Unscoped().Save(&attributeType).Error
	if err != nil {
		l.Logger.Error("❌ Failed to undo deleted attribute type", "method", "UndoDeletedAttributeType", "error", err, "id", id)
	}

	return err
}

func (r *AttributeTypeRepository) AttributeTypeExists(id uint, showDeleted *bool) (bool, error) {
	var attributeType AttributeTypeEntity
	query := r.db.Model(&AttributeTypeEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&attributeType).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if attribute type exists", "method", "AttributeTypeExists", "error", err, "id", id)
		return false, err
	}

	return attributeType.ID != 0, nil
}

func (r *AttributeTypeRepository) AttributeTypeExistsByName(name string, excludeID ...uint) (bool, error) {
	var attributeType AttributeTypeEntity
	query := r.db.Model(&AttributeTypeEntity{}).Where(c.AttributeTypeName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(c.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(c.FieldID).Take(&attributeType).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if attribute type exists by name", "method", "AttributeTypeExistsByName", "error", err, "name", name)
		return false, err
	}

	return attributeType.ID != 0, nil
}
