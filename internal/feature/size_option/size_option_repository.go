package size_option

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type SizeOptionRepository struct {
	db *gorm.DB
}

func NewSizeOptionRepository() *SizeOptionRepository {
	return &SizeOptionRepository{
		db: db.GetDB(),
	}
}

func (r *SizeOptionRepository) GetAllSizeOptions(include []string, showDeleted *bool, sizeCategoryID *uint, sortBy, sortOrder string) ([]SizeOption, error) {
	var sizeOptions []SizeOption

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if sizeCategoryID != nil {
		query = query.Where(constants.SizeOptionSizeCategoryID+" = ?", *sizeCategoryID)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&sizeOptions).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch size options", "method", "GetAllSizeOptions", "error", err, "include", include, "sizeCategoryID", sizeCategoryID, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return sizeOptions, err
}

func (r *SizeOptionRepository) GetSizeOptionByID(id uint, include []string, showDeleted *bool) (*SizeOption, error) {
	var sizeOption SizeOption

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&sizeOption).Error; err != nil {
		logger.Logger.Error("Failed to fetch size option by ID", "method", "GetSizeOptionByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &sizeOption, nil
}

func (r *SizeOptionRepository) CreateSizeOption(sizeOption *SizeOption) (*SizeOption, error) {
	err := r.db.Create(sizeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to create size option", "method", "CreateSizeOption", "error", err, "sizeOption", sizeOption)
		return nil, err
	}
	return sizeOption, nil
}

func (r *SizeOptionRepository) UpdateSizeOption(sizeOption *SizeOption) error {
	err := r.db.Save(sizeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to update size option", "method", "UpdateSizeOption", "error", err, "sizeOption", sizeOption)
	}
	return err
}

func (r *SizeOptionRepository) DeleteSizeOption(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&SizeOption{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete size option", "method", "DeleteSizeOption", "error", err, "id", id)
	}

	return err
}

func (r *SizeOptionRepository) UndoDeletedSizeOption(id uint) error {
	var sizeOption SizeOption
	err := r.db.Unscoped().Where(constants.FieldID+" = ?", id).First(&sizeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to find deleted size option", "method", "UndoDeletedSizeOption", "error", err, "id", id)
		return err
	}

	sizeOption.DeletedAt = nil
	err = r.db.Unscoped().Save(&sizeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to undo deleted size option", "method", "UndoDeletedSizeOption", "error", err, "id", id)
	}

	return err
}

func (r *SizeOptionRepository) SizeOptionExists(id uint, showDeleted *bool) (bool, error) {
	var sizeOption SizeOption
	query := r.db.Model(&SizeOption{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(constants.FieldID).Take(&sizeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to check if size option exists", "method", "SizeOptionExists", "error", err, "id", id)
		return false, err
	}

	return sizeOption.ID != 0, nil
}

func (r *SizeOptionRepository) SizeOptionExistsByName(name string, sizeCategoryID uint, excludeID ...uint) (bool, error) {
	var sizeOption SizeOption
	query := r.db.Model(&SizeOption{}).Where(constants.SizeOptionName+" = ? AND "+constants.SizeOptionSizeCategoryID+" = ?", name, sizeCategoryID)

	if len(excludeID) > 0 {
		query = query.Where(constants.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(constants.FieldID).Take(&sizeOption).Error
	if err != nil {
		logger.Logger.Error("Failed to check if size option exists by name", "method", "SizeOptionExistsByName", "error", err, "name", name, "sizeCategoryID", sizeCategoryID)
		return false, err
	}

	return sizeOption.ID != 0, nil
}

func (r *SizeOptionRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.SizeOptionName, constants.SizeOptionSortOrder, constants.SizeOptionSizeCategoryID, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
