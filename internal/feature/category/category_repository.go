package category

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: db.GetDB(),
	}
}

func (r *CategoryRepository) GetAllCategories(include []string, parentID *uint, isActive *bool) ([]Category, error) {
	var categories []Category

	// Build select fields dynamically
	selectFields := []string{"id", "title", "created_at", "updated_at"}

	// Add include fields to select
	for _, field := range include {
		switch field {
		case "sub_title", "image_url", "is_active", "parent_id":
			selectFields = append(selectFields, field)
		}
	}

	// Create query with selected fields
	query := r.db.Select(strings.Join(selectFields, ", "))

	// Apply filters
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetCategoryByID(id uint, include []string) (*Category, error) {
	var category Category

	// Build select fields dynamically
	selectFields := []string{"id", "title", "created_at", "updated_at"}

	// Add include fields to select
	for _, field := range include {
		switch field {
		case "sub_title", "image_url", "is_active", "parent_id":
			selectFields = append(selectFields, field)
		}
	}

	// Create query with selected fields
	query := r.db.Select(strings.Join(selectFields, ", "))

	err := query.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(category *Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) UpdateCategory(category *Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&Category{}, id).Error
}
