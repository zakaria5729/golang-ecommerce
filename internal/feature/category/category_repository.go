package category

import (
	"errors"
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
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

func (r *CategoryRepository) GetAllCategories(include []string, parentID *uint, showPriority *bool, sortBy, sortOrder string) ([]Category, error) {
	var categories []Category

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(CategoryIsActive+" = ?", true)

	if parentID != nil {
		query = query.Where(CategoryParentID+" = ?", *parentID)
	}

	if showPriority != nil && *showPriority {
		query = query.Where(CategoryPriority+" >= ?", getPriorityLimit())
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&categories).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch categories", "method", "GetAllCategories", "error", err, "include", include, "parentID", parentID, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return categories, err
}

func (r *CategoryRepository) GetAllCategoriesPaginated(include []string, parentID *uint, page, pageSize int, showPriority *bool, sortBy, sortOrder string) ([]Category, int, error) {
	var categories []Category
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(CategoryIsActive+" = ?", true)

	if parentID != nil {
		query = query.Where(CategoryParentID+" = ?", *parentID)
	}

	if showPriority != nil && *showPriority {
		query = query.Where(CategoryPriority+" >= ?", getPriorityLimit())
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Model(&Category{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count categories", "method", "GetAllCategoriesPaginated", "error", err, "include", include, "parentID", parentID, "page", page, "pageSize", pageSize, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&categories).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch categories paginated", "method", "GetAllCategoriesPaginated", "error", err, "include", include, "parentID", parentID, "page", page, "pageSize", pageSize, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return categories, int(total), err
}

func (r *CategoryRepository) GetCategoryByID(id uint, include []string) (*Category, error) {
	var category Category

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Where(constants.FieldID+" = ?", id).First(&category).Error; err != nil {
		logger.Logger.Error("Failed to fetch category by ID", "method", "GetCategoryByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetCategoryByIDIncludeInactive(id uint, include []string) (*Category, error) {
	var category Category

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Where(constants.FieldID+" = ?", id).First(&category).Error; err != nil {
		logger.Logger.Error("Failed to fetch category by ID (include inactive)", "method", "GetCategoryByIDIncludeInactive", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(category *Category) error {
	err := r.db.Create(category).Error
	if err != nil {
		logger.Logger.Error("Failed to create category", "method", "CreateCategory", "error", err, "category", category)
	}
	return err
}

func (r *CategoryRepository) UpdateCategory(category *Category) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
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
	if err != nil {
		logger.Logger.Error("Failed to update category", "method", "UpdateCategory", "error", err, "category", category)
	}
	return err
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
		logger.Logger.Error("Failed to delete category", "method", "DeleteCategory", "error", err, "id", id)
		return err
	}

	if err := r.deleteSubcategoriesRecursively(tx, id); err != nil {
		tx.Rollback()
		logger.Logger.Error("Failed to delete subcategories", "method", "DeleteCategory", "error", err, "id", id)
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		logger.Logger.Error("Failed to commit category deletion", "method", "DeleteCategory", "error", err, "id", id)
	}
	return err
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
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Category{}).Where(constants.FieldID+" = ?", id).Update(CategoryIsActive, isActive).Error; err != nil {
			return err
		}

		return r.updateSubcategoriesStatusRecursively(tx, id, isActive)
	})
	if err != nil {
		logger.Logger.Error("Failed to update category status", "method", "UpdateCategoryStatus", "error", err, "id", id, "isActive", isActive)
	}
	return err
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
	if err != nil {
		logger.Logger.Error("Failed to check if category exists", "method", "CategoryExists", "error", err, "id", id)
	}
	return count > 0, err
}

func (r *CategoryRepository) CategoryExistsByTitle(title string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&Category{}).Where(CategoryTitle+" = ?", title)

	if excludeID != nil {
		query = query.Where(constants.FieldID+" != ?", *excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if category exists by title", "method", "CategoryExistsByTitle", "error", err, "title", title, "excludeID", excludeID)
	}
	return count > 0, err
}

func (r *CategoryRepository) HasChildren(parentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Category{}).Where(CategoryParentID+" = ?", parentID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if category has children", "method", "HasChildren", "error", err, "parentID", parentID)
	}
	return count > 0, err
}

func (r *CategoryRepository) IncrementPriority(categoryID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		rootID, err := r.findRootCategoryID(tx, categoryID)
		if err != nil {
			return err
		}

		result := tx.Model(&Category{}).
			Where(constants.FieldID+" = ?", rootID).
			Update(CategoryPriority, gorm.Expr(CategoryPriority+" + 1"))

		return result.Error
	})
	if err != nil {
		logger.Logger.Error("Failed to increment priority", "method", "IncrementPriority", "error", err, "categoryID", categoryID)
	}
	return err
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
	defaultFields := []string{constants.FieldID, CategoryTitle, CategoryIsActive, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{CategorySubTitle, CategoryImageURL, CategoryParentID, CategoryPriority}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}

func getPriorityLimit() uint {
	return 10
}
