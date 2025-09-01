package category

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/utils"
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

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	if parentID != nil {
		query = query.Where(CategoryParentID+" = ?", *parentID)
	}

	if isActive != nil {
		query = query.Where(CategoryIsActive+" = ?", *isActive)
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetCategoryByID(id uint, include []string) (*Category, error) {
	var category Category

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	err := query.Where(CategoryID+" = ?", id).First(&category).Error
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

func (r *CategoryRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{CategoryID, CategoryTitle, CategoryCreatedAt, CategoryUpdatedAt}
	optionalFields := []string{CategorySubTitle, CategoryImageURL, CategoryIsActive, CategoryParentID}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
