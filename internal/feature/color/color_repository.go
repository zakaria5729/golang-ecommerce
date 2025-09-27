package color

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type ColorRepository struct {
	db *gorm.DB
}

func NewColorRepository() *ColorRepository {
	return &ColorRepository{
		db: db.GetDB(),
	}
}

func (r *ColorRepository) GetAllColors(include []string, showDeleted *bool, sortBy, sortOrder string) ([]Color, error) {
	var colors []Color

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&colors).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch colors", "method", "GetAllColors", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return colors, err
}

func (r *ColorRepository) GetColorByID(id uint, include []string, showDeleted *bool) (*Color, error) {
	var color Color

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&color).Error; err != nil {
		logger.Logger.Error("Failed to fetch color by ID", "method", "GetColorByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &color, nil
}

func (r *ColorRepository) CreateColor(color *Color) (*Color, error) {
	err := r.db.Create(color).Error
	if err != nil {
		logger.Logger.Error("Failed to create color", "method", "CreateColor", "error", err, "color", color)
		return nil, err
	}
	return color, nil
}

func (r *ColorRepository) UpdateColor(color *Color) error {
	err := r.db.Save(color).Error
	if err != nil {
		logger.Logger.Error("Failed to update color", "method", "UpdateColor", "error", err, "color", color)
	}
	return err
}

func (r *ColorRepository) DeleteColor(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&Color{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete color", "method", "DeleteColor", "error", err, "id", id)
	}

	return err
}

func (r *ColorRepository) UndoDeletedColor(id uint) error {
	err := r.db.Unscoped().Model(&Color{}).Where(constants.FieldID+" = ?", id).Update(constants.FieldDeletedAt, nil).Error

	if err != nil {
		logger.Logger.Error("Failed to undo deleted color", "method", "UndoDeletedColor", "error", err, "id", id)
	}

	return err
}

func (r *ColorRepository) ColorExists(id uint, showDeleted *bool) (bool, error) {
	var color Color
	query := r.db.Model(&Color{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(constants.FieldID).Take(&color).Error
	if err != nil {
		logger.Logger.Error("Failed to check if color exists", "method", "ColorExists", "error", err, "id", id)
		return false, err
	}

	return color.ID != 0, nil
}

func (r *ColorRepository) ColorExistsByName(name string, excludeID ...uint) (bool, error) {
	var color Color
	query := r.db.Model(&Color{}).Where(constants.ColorName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(constants.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(constants.FieldID).Take(&color).Error
	if err != nil {
		logger.Logger.Error("Failed to check if color exists by name", "method", "ColorExistsByName", "error", err, "name", name)
		return false, err
	}

	return color.ID != 0, nil
}

func (r *ColorRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.ColorName, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
