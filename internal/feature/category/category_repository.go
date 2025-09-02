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

	if isActive == nil {
		query = query.Where(CategoryIsActive+" = ?", true)
	} else {
		query = query.Where(CategoryIsActive+" = ?", *isActive)
	}

	if parentID != nil {
		query = query.Where(CategoryParentID+" = ?", *parentID)
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetAllCategoriesPaginated(include []string, parentID *uint, isActive *bool, page, limit int) ([]Category, int, error) {
	var categories []Category
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if isActive == nil {
		query = query.Where(CategoryIsActive+" = ?", true)
	} else {
		query = query.Where(CategoryIsActive+" = ?", *isActive)
	}

	if parentID != nil {
		query = query.Where(CategoryParentID+" = ?", *parentID)
	}

	if err := query.Model(&Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&categories).Error

	return categories, int(total), err
}

func (r *CategoryRepository) GetCategoryByID(id uint, include []string) (*Category, error) {
	var category Category

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	query = query.Where(CategoryIsActive+" = ?", true)
	err := query.Where(CategoryID+" = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetCategoryByIDIncludeInactive(id uint, include []string) (*Category, error) {
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
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&Category{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.deleteSubcategoriesRecursively(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *CategoryRepository) deleteSubcategoriesRecursively(tx *gorm.DB, parentID uint) error {
	var subcategories []Category
	if err := tx.Where(CategoryParentID+" = ?", parentID).Find(&subcategories).Error; err != nil {
		return err
	}

	for _, subcategory := range subcategories {
		if err := tx.Delete(&Category{}, subcategory.ID).Error; err != nil {
			return err
		}

		if err := r.deleteSubcategoriesRecursively(tx, subcategory.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *CategoryRepository) UpdateCategoryStatus(id uint, isActive bool) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&Category{}).Where(CategoryID+" = ?", id).Update(CategoryIsActive, isActive).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateSubcategoriesStatusRecursively(tx, id, isActive); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *CategoryRepository) updateSubcategoriesStatusRecursively(tx *gorm.DB, parentID uint, isActive bool) error {
	var subcategories []Category
	if err := tx.Where(CategoryParentID+" = ?", parentID).Find(&subcategories).Error; err != nil {
		return err
	}

	for _, subcategory := range subcategories {
		if err := tx.Model(&Category{}).Where(CategoryID+" = ?", subcategory.ID).Update(CategoryIsActive, isActive).Error; err != nil {
			return err
		}

		if err := r.updateSubcategoriesStatusRecursively(tx, subcategory.ID, isActive); err != nil {
			return err
		}
	}

	return nil
}

func (r *CategoryRepository) CategoryExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where(CategoryID+" = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) CategoryExistsByTitle(title string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&Category{}).Where(CategoryTitle+" = ?", title)

	if excludeID != nil {
		query = query.Where(CategoryID+" != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) HasChildren(parentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where(CategoryParentID+" = ?", parentID).Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{CategoryID, CategoryTitle, CategoryCreatedAt, CategoryUpdatedAt}
	optionalFields := []string{CategorySubTitle, CategoryImageURL, CategoryIsActive, CategoryParentID}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
