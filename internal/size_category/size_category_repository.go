package size_category

import (
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type SizeCategoryRepository struct {
	db *gorm.DB
}

func NewSizeCategoryRepository(db *gorm.DB) *SizeCategoryRepository {
	return &SizeCategoryRepository{
		db: db,
	}
}

func (r *SizeCategoryRepository) GetAllSizeCategories(showDeleted *bool, sortBy, sortOrder string) ([]SizeCategoryEntity, error) {
	var sizeCategories []SizeCategoryEntity
	query := r.db.Model(&SizeCategoryEntity{})

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&sizeCategories).Error
	if err != nil {
		l.Logger.Error("Failed to fetch size categories", "method", "GetAllSizeCategories", "error", err, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return sizeCategories, err
}

func (r *SizeCategoryRepository) GetSizeCategoryByID(id uint, showDeleted *bool) (*SizeCategoryEntity, error) {
	var sizeCategory SizeCategoryEntity
	query := r.db.Model(&SizeCategoryEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&sizeCategory).Error; err != nil {
		l.Logger.Error("Failed to fetch size category by ID", "method", "GetSizeCategoryByID", "error", err, "id", id)
		return nil, err
	}

	return &sizeCategory, nil
}

func (r *SizeCategoryRepository) CreateSizeCategory(sizeCategory *SizeCategoryEntity) (*SizeCategoryEntity, error) {
	err := r.db.Create(sizeCategory).Error
	if err != nil {
		l.Logger.Error("Failed to create size category", "method", "CreateSizeCategory", "error", err, "sizeCategory", sizeCategory)
		return nil, err
	}
	return sizeCategory, nil
}

func (r *SizeCategoryRepository) UpdateSizeCategory(sizeCategory *SizeCategoryEntity) error {
	err := r.db.Save(sizeCategory).Error
	if err != nil {
		l.Logger.Error("Failed to update size category", "method", "UpdateSizeCategory", "error", err, "sizeCategory", sizeCategory)
	}
	return err
}

func (r *SizeCategoryRepository) DeleteSizeCategory(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&SizeCategoryEntity{}).Error
	if err != nil {
		l.Logger.Error("Failed to delete size category", "method", "DeleteSizeCategory", "error", err, "id", id)
	}

	return err
}

func (r *SizeCategoryRepository) UndoDeletedSizeCategory(id uint) error {
	var sizeCategory SizeCategoryEntity
	err := r.db.Unscoped().Where(c.FieldID+" = ?", id).First(&sizeCategory).Error
	if err != nil {
		l.Logger.Error("Failed to find deleted size category", "method", "UndoDeletedSizeCategory", "error", err, "id", id)
		return err
	}

	sizeCategory.DeletedAt = nil
	err = r.db.Unscoped().Save(&sizeCategory).Error
	if err != nil {
		l.Logger.Error("Failed to undo deleted size category", "method", "UndoDeletedSizeCategory", "error", err, "id", id)
	}

	return err
}

func (r *SizeCategoryRepository) SizeCategoryExists(id uint, showDeleted *bool) (bool, error) {
	var sizeCategory SizeCategoryEntity
	query := r.db.Model(&SizeCategoryEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&sizeCategory).Error
	if err != nil {
		l.Logger.Error("Failed to check if size category exists", "method", "SizeCategoryExists", "error", err, "id", id)
		return false, err
	}

	return sizeCategory.ID != 0, nil
}

func (r *SizeCategoryRepository) SizeCategoryExistsByName(name string, excludeID ...uint) (bool, error) {
	var sizeCategory SizeCategoryEntity
	query := r.db.Model(&SizeCategoryEntity{}).Where(c.SizeCategoryName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(c.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(c.FieldID).Take(&sizeCategory).Error
	if err != nil {
		l.Logger.Error("Failed to check if size category exists by name", "method", "SizeCategoryExistsByName", "error", err, "name", name)
		return false, err
	}

	return sizeCategory.ID != 0, nil
}
