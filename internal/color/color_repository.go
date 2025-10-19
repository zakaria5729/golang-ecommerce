package color

import (
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type ColorRepository struct {
	db *gorm.DB
}

func NewColorRepository(db *gorm.DB) *ColorRepository {
	return &ColorRepository{
		db: db,
	}
}

func (r *ColorRepository) GetAllColors(showDeleted *bool, sortBy, sortOrder string) ([]ColorEntity, error) {
	var colors []ColorEntity
	query := r.db.Model(&ColorEntity{})

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&colors).Error
	if err != nil {
		l.Logger.Error("❌ Failed to fetch colors", "method", "GetAllColors", "error", err, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return colors, err
}

func (r *ColorRepository) GetColorByID(id uint, showDeleted *bool) (*ColorEntity, error) {
	var color ColorEntity
	query := r.db.Model(&ColorEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&color).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch color by ID", "method", "GetColorByID", "error", err, "id", id)
		return nil, err
	}

	return &color, nil
}

func (r *ColorRepository) CreateColor(color *ColorEntity) (*ColorEntity, error) {
	err := r.db.Create(color).Error
	if err != nil {
		l.Logger.Error("❌ Failed to create color", "method", "CreateColor", "error", err, "color", color)
		return nil, err
	}

	return color, nil
}

func (r *ColorRepository) UpdateColor(color *ColorEntity) error {
	err := r.db.Save(color).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update color", "method", "UpdateColor", "error", err, "color", color)
	}

	return err
}

func (r *ColorRepository) DeleteColor(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&ColorEntity{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete color", "method", "DeleteColor", "error", err, "id", id)
	}

	return err
}

func (r *ColorRepository) UndoDeletedColor(id uint) error {
	err := r.db.Unscoped().Model(&ColorEntity{}).Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil).Error
	if err != nil {
		l.Logger.Error("❌ Failed to undo deleted color", "method", "UndoDeletedColor", "error", err, "id", id)
	}

	return err
}

func (r *ColorRepository) ColorExists(id uint, showDeleted *bool) (bool, error) {
	var color ColorEntity
	query := r.db.Model(&ColorEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&color).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if color exists", "method", "ColorExists", "error", err, "id", id)
		return false, err
	}

	return color.ID != 0, nil
}

func (r *ColorRepository) ColorExistsByName(name string, excludeID ...uint) (bool, error) {
	var color ColorEntity
	query := r.db.Model(&ColorEntity{}).Where(c.ColorName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(c.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(c.FieldID).Take(&color).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if color exists by name", "method", "ColorExistsByName", "error", err, "name", name)
		return false, err
	}

	return color.ID != 0, nil
}
