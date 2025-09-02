package category

import (
	"errors"
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
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

func (r *CategoryRepository) GetAllCategories(include []string, parentID *uint, isActive *bool, priority *bool) ([]Category, error) {
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

	if priority != nil {
		query = query.Where(CategoryPriority+" = ?", *priority)
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetAllCategoriesPaginated(include []string, parentID *uint, isActive *bool, page, pageSize int, priority *bool) ([]Category, int, error) {
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

	if priority != nil {
		query = query.Where(CategoryPriority+" = ?", *priority)
	}

	if err := query.Model(&Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&categories).Error

	return categories, int(total), err
}

func (r *CategoryRepository) GetCategoryByID(id uint, include []string) (*Category, error) {
	var category Category

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	query = query.Where(CategoryIsActive+" = ?", true)
	err := query.Where(constants.FieldID+" = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetCategoryByIDIncludeInactive(id uint, include []string) (*Category, error) {
	var category Category

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	err := query.Where(constants.FieldID+" = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(category *Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) UpdateCategory(category *Category) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var currentCategory Category
		if err := tx.First(&currentCategory, category.ID).Error; err != nil {
			return err
		}

		if err := tx.Save(category).Error; err != nil {
			return err
		}

		if currentCategory.IsActive != category.IsActive {
			if err := r.updateSubcategoriesStatusRecursively(tx, category.ID, category.IsActive); err != nil {
				return err
			}
		}

		return nil
	})
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Category{}).Where(constants.FieldID+" = ?", id).Update(CategoryIsActive, isActive).Error; err != nil {
			return err
		}

		return r.updateSubcategoriesStatusRecursively(tx, id, isActive)
	})
}

func (r *CategoryRepository) updateSubcategoriesStatusRecursively(tx *gorm.DB, parentID uint, isActive bool) error {
	var subcategories []Category
	if err := tx.Where(CategoryParentID+" = ?", parentID).Find(&subcategories).Error; err != nil {
		return err
	}

	for _, subcategory := range subcategories {
		if err := tx.Model(&Category{}).Where(constants.FieldID+" = ?", subcategory.ID).Update(CategoryIsActive, isActive).Error; err != nil {
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
	err := r.db.Model(&Category{}).Where(constants.FieldID+" = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) CategoryExistsByTitle(title string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&Category{}).Where(CategoryTitle+" = ?", title)

	if excludeID != nil {
		query = query.Where(constants.FieldID+" != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) HasChildren(parentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where(CategoryParentID+" = ?", parentID).Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) IncrementPriority(categoryID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		rootID, err := r.findRootCategoryID(tx, categoryID)
		if err != nil {
			return err
		}

		result := tx.Model(&Category{}).
			Where(constants.FieldID+" = ?", rootID).
			Update(CategoryPriority, gorm.Expr(CategoryPriority+" + 1"))

		return result.Error
	})
}

func (r *CategoryRepository) findRootCategoryID(tx *gorm.DB, categoryID uint) (uint, error) {
	currentID := categoryID
	visited := make(map[uint]bool)

	for {
		if visited[currentID] {
			return 0, fmt.Errorf("circular reference detected in category hierarchy")
		}
		visited[currentID] = true

		var category Category
		err := tx.Select(constants.FieldID, CategoryParentID).
			Where(constants.FieldID+" = ?", currentID).
			First(&category).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("category with ID %d not found", currentID)
			}
			return 0, err
		}

		if category.ParentID == nil {
			return category.ID, nil
		}
		currentID = *category.ParentID
	}
}

func (r *CategoryRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, CategoryTitle, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{CategorySubTitle, CategoryImageURL, CategoryIsActive, CategoryParentID, CategoryPriority}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}

func getPriority() int {
	return 10
}
