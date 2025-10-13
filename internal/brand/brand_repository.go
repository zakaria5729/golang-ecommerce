package brand

import (
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type BrandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) *BrandRepository {
	return &BrandRepository{
		db: db,
	}
}

func (r *BrandRepository) GetAllBrands(showDeleted *bool, sortBy, sortOrder string) ([]Brand, error) {
	var brands []Brand
	query := r.db.Model(&Brand{})

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&brands).Error
	if err != nil {
		l.Logger.Error("❌ Failed to fetch brands", "method", "GetAllBrands", "error", err, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return brands, err
}

func (r *BrandRepository) GetBrandByID(id uint, showDeleted *bool) (*Brand, error) {
	var brand Brand
	query := r.db.Model(&Brand{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&brand).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch brand by ID", "method", "GetBrandByID", "error", err, "id", id)
		return nil, err
	}

	return &brand, nil
}

func (r *BrandRepository) CreateBrand(brand *Brand) (*Brand, error) {
	err := r.db.Create(brand).Error
	if err != nil {
		l.Logger.Error("❌ Failed to create brand", "method", "CreateBrand", "error", err, "brand", brand)
		return nil, err
	}

	return brand, nil
}

func (r *BrandRepository) UpdateBrand(brand *Brand) error {
	err := r.db.Save(brand).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update brand", "method", "UpdateBrand", "error", err, "brand", brand)
	}

	return err
}

func (r *BrandRepository) DeleteBrand(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&Brand{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete brand", "method", "DeleteBrand", "error", err, "id", id)
	}

	return err
}

func (r *BrandRepository) UndoDeletedBrand(id uint) error {
	err := r.db.Unscoped().Model(&Brand{}).Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil).Error
	if err != nil {
		l.Logger.Error("❌ Failed to undo deleted brand", "method", "UndoDeletedBrand", "error", err, "id", id)
	}

	return err
}

func (r *BrandRepository) BrandExists(id uint, showDeleted *bool) (bool, error) {
	var brand Brand
	query := r.db.Model(&Brand{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&brand).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if brand exists", "method", "BrandExists", "error", err, "id", id)
		return false, err
	}

	return brand.ID != 0, nil
}

func (r *BrandRepository) BrandExistsByName(name string, excludeID ...uint) (bool, error) {
	var brand Brand
	query := r.db.Model(&Brand{}).Where(c.BrandName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(c.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(c.FieldID).Take(&brand).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if brand exists by name", "method", "BrandExistsByName", "error", err, "name", name)
		return false, err
	}

	return brand.ID != 0, nil
}
