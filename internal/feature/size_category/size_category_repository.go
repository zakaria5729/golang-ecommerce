package size_category

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type SizeCategoryRepository struct {
	db *gorm.DB
}

func NewSizeCategoryRepository() *SizeCategoryRepository {
	return &SizeCategoryRepository{
		db: db.GetDB(),
	}
}

func (r *SizeCategoryRepository) GetAllSizeCategories(include []string, showDeleted *bool, sortBy, sortOrder string) ([]SizeCategory, error) {
	var sizeCategories []SizeCategory

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&sizeCategories).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch size categories", "method", "GetAllSizeCategories", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return sizeCategories, err
}

func (r *SizeCategoryRepository) GetSizeCategoryByID(id uint, include []string, showDeleted *bool) (*SizeCategory, error) {
	var sizeCategory SizeCategory

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&sizeCategory).Error; err != nil {
		logger.Logger.Error("Failed to fetch size category by ID", "method", "GetSizeCategoryByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &sizeCategory, nil
}

func (r *SizeCategoryRepository) CreateSizeCategory(sizeCategory *SizeCategory) (*SizeCategory, error) {
	err := r.db.Create(sizeCategory).Error
	if err != nil {
		logger.Logger.Error("Failed to create size category", "method", "CreateSizeCategory", "error", err, "sizeCategory", sizeCategory)
		return nil, err
	}
	return sizeCategory, nil
}

func (r *SizeCategoryRepository) UpdateSizeCategory(sizeCategory *SizeCategory) error {
	err := r.db.Save(sizeCategory).Error
	if err != nil {
		logger.Logger.Error("Failed to update size category", "method", "UpdateSizeCategory", "error", err, "sizeCategory", sizeCategory)
	}
	return err
}

func (r *SizeCategoryRepository) DeleteSizeCategory(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&SizeCategory{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete size category", "method", "DeleteSizeCategory", "error", err, "id", id)
	}

	return err
}

func (r *SizeCategoryRepository) UndoDeletedSizeCategory(id uint) error {
	var sizeCategory SizeCategory
	err := r.db.Unscoped().Where(constants.FieldID+" = ?", id).First(&sizeCategory).Error
	if err != nil {
		logger.Logger.Error("Failed to find deleted size category", "method", "UndoDeletedSizeCategory", "error", err, "id", id)
		return err
	}

	sizeCategory.DeletedAt = nil
	err = r.db.Unscoped().Save(&sizeCategory).Error
	if err != nil {
		logger.Logger.Error("Failed to undo deleted size category", "method", "UndoDeletedSizeCategory", "error", err, "id", id)
	}

	return err
}

func (r *SizeCategoryRepository) SizeCategoryExists(id uint, showDeleted *bool) (bool, error) {
	var sizeCategory SizeCategory
	query := r.db.Model(&SizeCategory{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(constants.FieldID).Take(&sizeCategory).Error
	if err != nil {
		logger.Logger.Error("Failed to check if size category exists", "method", "SizeCategoryExists", "error", err, "id", id)
		return false, err
	}

	return sizeCategory.ID != 0, nil
}

func (r *SizeCategoryRepository) SizeCategoryExistsByName(name string, excludeID ...uint) (bool, error) {
	var sizeCategory SizeCategory
	query := r.db.Model(&SizeCategory{}).Where(constants.SizeCategoryName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(constants.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(constants.FieldID).Take(&sizeCategory).Error
	if err != nil {
		logger.Logger.Error("Failed to check if size category exists by name", "method", "SizeCategoryExistsByName", "error", err, "name", name)
		return false, err
	}

	return sizeCategory.ID != 0, nil
}

func (r *SizeCategoryRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.SizeCategoryName, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
