package brand

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type BrandRepository struct {
	db *gorm.DB
}

func NewBrandRepository() *BrandRepository {
	return &BrandRepository{
		db: db.GetDB(),
	}
}

func (r *BrandRepository) GetAllBrands(include []string, showDeleted *bool, sortBy, sortOrder string) ([]Brand, error) {
	var brands []Brand

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&brands).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch brands", "method", "GetAllBrands", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return brands, err
}

func (r *BrandRepository) GetBrandByID(id uint, include []string, showDeleted *bool) (*Brand, error) {
	var brand Brand

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&brand).Error; err != nil {
		logger.Logger.Error("Failed to fetch brand by ID", "method", "GetBrandByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &brand, nil
}

func (r *BrandRepository) CreateBrand(brand *Brand) (*Brand, error) {
	err := r.db.Create(brand).Error
	if err != nil {
		logger.Logger.Error("Failed to create brand", "method", "CreateBrand", "error", err, "brand", brand)
		return nil, err
	}
	return brand, nil
}

func (r *BrandRepository) UpdateBrand(brand *Brand) error {
	err := r.db.Save(brand).Error
	if err != nil {
		logger.Logger.Error("Failed to update brand", "method", "UpdateBrand", "error", err, "brand", brand)
	}
	return err
}

func (r *BrandRepository) DeleteBrand(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&Brand{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete brand", "method", "DeleteBrand", "error", err, "id", id)
	}

	return err
}

func (r *BrandRepository) UndoDeletedBrand(id uint) error {
	err := r.db.Unscoped().Model(&Brand{}).Where(constants.FieldID+" = ?", id).Update(constants.FieldDeletedAt, nil).Error

	if err != nil {
		logger.Logger.Error("Failed to undo deleted brand", "method", "UndoDeletedBrand", "error", err, "id", id)
	}

	return err
}

func (r *BrandRepository) BrandExists(id uint, showDeleted *bool) (bool, error) {
	var brand Brand
	query := r.db.Model(&Brand{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(constants.FieldID).Take(&brand).Error
	if err != nil {
		logger.Logger.Error("Failed to check if brand exists", "method", "BrandExists", "error", err, "id", id)
		return false, err
	}

	return brand.ID != 0, nil
}

func (r *BrandRepository) BrandExistsByName(name string, excludeID ...uint) (bool, error) {
	var brand Brand
	query := r.db.Model(&Brand{}).Where(constants.BrandName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(constants.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(constants.FieldID).Take(&brand).Error
	if err != nil {
		logger.Logger.Error("Failed to check if brand exists by name", "method", "BrandExistsByName", "error", err, "name", name)
		return false, err
	}

	return brand.ID != 0, nil
}

func (r *BrandRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.BrandName, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{constants.BrandDescription}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
