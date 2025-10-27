package size_option

import (
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type SizeOptionRepository interface {
	GetAllSizeOptions(showDeleted *bool, sizeCategoryID *uint, sortBy, sortOrder string) ([]SizeOptionEntity, error)
	GetSizeOptionByID(id uint, showDeleted *bool) (*SizeOptionEntity, error)
	CreateSizeOption(sizeOption *SizeOptionEntity) (*SizeOptionEntity, error)
	UpdateSizeOption(sizeOption *SizeOptionEntity) error
	DeleteSizeOption(id uint) error
	UndoDeletedSizeOption(id uint) error
	SizeOptionExists(id uint, showDeleted *bool) (bool, error)
	SizeOptionExistsByName(name string, sizeCategoryID uint, excludeID ...uint) (bool, error)
}

type sizeOptionRepository struct {
	db *gorm.DB
}

func NewSizeOptionRepository(db *gorm.DB) SizeOptionRepository {
	return &sizeOptionRepository{
		db: db,
	}
}

func (r *sizeOptionRepository) GetAllSizeOptions(showDeleted *bool, sizeCategoryID *uint, sortBy, sortOrder string) ([]SizeOptionEntity, error) {
	var sizeOptions []SizeOptionEntity
	query := r.db.Model(&SizeOptionEntity{})

	if sizeCategoryID != nil {
		query = query.Where(c.SizeOptionSizeCategoryID+" = ?", *sizeCategoryID)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&sizeOptions).Error
	if err != nil {
		l.Logger.Error("Failed to fetch size options", "method", "GetAllSizeOptions", "error", err, "sizeCategoryID", sizeCategoryID, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return sizeOptions, err
}

func (r *sizeOptionRepository) GetSizeOptionByID(id uint, showDeleted *bool) (*SizeOptionEntity, error) {
	var sizeOption SizeOptionEntity
	query := r.db.Model(&SizeOptionEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&sizeOption).Error; err != nil {
		l.Logger.Error("Failed to fetch size option by ID", "method", "GetSizeOptionByID", "error", err, "id", id)
		return nil, err
	}

	return &sizeOption, nil
}

func (r *sizeOptionRepository) CreateSizeOption(sizeOption *SizeOptionEntity) (*SizeOptionEntity, error) {
	err := r.db.Create(sizeOption).Error
	if err != nil {
		l.Logger.Error("Failed to create size option", "method", "CreateSizeOption", "error", err, "sizeOption", sizeOption)
		return nil, err
	}
	return sizeOption, nil
}

func (r *sizeOptionRepository) UpdateSizeOption(sizeOption *SizeOptionEntity) error {
	err := r.db.Save(sizeOption).Error
	if err != nil {
		l.Logger.Error("Failed to update size option", "method", "UpdateSizeOption", "error", err, "sizeOption", sizeOption)
	}
	return err
}

func (r *sizeOptionRepository) DeleteSizeOption(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&SizeOptionEntity{}).Error
	if err != nil {
		l.Logger.Error("Failed to delete size option", "method", "DeleteSizeOption", "error", err, "id", id)
	}

	return err
}

func (r *sizeOptionRepository) UndoDeletedSizeOption(id uint) error {
	var sizeOption SizeOptionEntity
	err := r.db.Unscoped().Where(c.FieldID+" = ?", id).First(&sizeOption).Error
	if err != nil {
		l.Logger.Error("Failed to find deleted size option", "method", "UndoDeletedSizeOption", "error", err, "id", id)
		return err
	}

	sizeOption.DeletedAt = nil
	err = r.db.Unscoped().Save(&sizeOption).Error
	if err != nil {
		l.Logger.Error("Failed to undo deleted size option", "method", "UndoDeletedSizeOption", "error", err, "id", id)
	}

	return err
}

func (r *sizeOptionRepository) SizeOptionExists(id uint, showDeleted *bool) (bool, error) {
	var sizeOption SizeOptionEntity
	query := r.db.Model(&SizeOptionEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&sizeOption).Error
	if err != nil {
		l.Logger.Error("Failed to check if size option exists", "method", "SizeOptionExists", "error", err, "id", id)
		return false, err
	}

	return sizeOption.ID != 0, nil
}

func (r *sizeOptionRepository) SizeOptionExistsByName(name string, sizeCategoryID uint, excludeID ...uint) (bool, error) {
	var sizeOption SizeOptionEntity
	query := r.db.Model(&SizeOptionEntity{}).Where(c.SizeOptionName+" = ? AND "+c.SizeOptionSizeCategoryID+" = ?", name, sizeCategoryID)

	if len(excludeID) > 0 {
		query = query.Where(c.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(c.FieldID).Take(&sizeOption).Error
	if err != nil {
		l.Logger.Error("Failed to check if size option exists by name", "method", "SizeOptionExistsByName", "error", err, "name", name, "sizeCategoryID", sizeCategoryID)
		return false, err
	}

	return sizeOption.ID != 0, nil
}
