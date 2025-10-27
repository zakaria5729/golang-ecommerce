package brand

import (
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type BrandRepository interface {
	GetAllBrands(showDeleted *bool, sortBy, sortOrder string) ([]BrandEntity, error)
	GetBrandByID(id uint, showDeleted *bool) (*BrandEntity, error)
	CreateBrand(brand *BrandEntity) (*BrandEntity, error)
	UpdateBrand(brand *BrandEntity) error
	DeleteBrand(id uint) error
	UndoDeletedBrand(id uint) error
	BrandExists(id uint, showDeleted *bool) (bool, error)
	BrandExistsByName(name string, excludeID ...uint) (bool, error)
}

type brandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{
		db: db,
	}
}

func (r *brandRepository) GetAllBrands(showDeleted *bool, sortBy, sortOrder string) ([]BrandEntity, error) {
	var brands []BrandEntity
	query := r.db.Model(&BrandEntity{})

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

func (r *brandRepository) GetBrandByID(id uint, showDeleted *bool) (*BrandEntity, error) {
	var brand BrandEntity
	query := r.db.Model(&BrandEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&brand).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch brand by ID", "method", "GetBrandByID", "error", err, "id", id)
		return nil, err
	}

	return &brand, nil
}

func (r *brandRepository) CreateBrand(brand *BrandEntity) (*BrandEntity, error) {
	err := r.db.Create(brand).Error
	if err != nil {
		l.Logger.Error("❌ Failed to create brand", "method", "CreateBrand", "error", err, "brand", brand)
		return nil, err
	}

	return brand, nil
}

func (r *brandRepository) UpdateBrand(brand *BrandEntity) error {
	err := r.db.Save(brand).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update brand", "method", "UpdateBrand", "error", err, "brand", brand)
	}

	return err
}

func (r *brandRepository) DeleteBrand(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&BrandEntity{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete brand", "method", "DeleteBrand", "error", err, "id", id)
	}

	return err
}

func (r *brandRepository) UndoDeletedBrand(id uint) error {
	err := r.db.Unscoped().Model(&BrandEntity{}).Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil).Error
	if err != nil {
		l.Logger.Error("❌ Failed to undo deleted brand", "method", "UndoDeletedBrand", "error", err, "id", id)
	}

	return err
}

func (r *brandRepository) BrandExists(id uint, showDeleted *bool) (bool, error) {
	var brand BrandEntity
	query := r.db.Model(&BrandEntity{}).Where(c.FieldID+" = ?", id)

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

func (r *brandRepository) BrandExistsByName(name string, excludeID ...uint) (bool, error) {
	var brand BrandEntity
	query := r.db.Model(&BrandEntity{}).Where(c.BrandName+" = ?", name)

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
