package category

import (
	"errors"
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/constants"
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

func (r *CategoryRepository) GetAllCategories(isActive *bool, showDeleted *bool, include []string, parentID *uint, priorityLimit *int, sortBy, sortOrder string) ([]Category, error) {
	var categories []Category

	var maxPriorityLimit int = c.MaxPriorityLimit
	if priorityLimit != nil && *priorityLimit > maxPriorityLimit {
		priorityLimit = &maxPriorityLimit
	}

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if isActive != nil {
		query = query.Where(c.CategoryIsActive+" = ?", *isActive)
	}

	if parentID != nil {
		query = query.Where(c.CategoryParentID+" = ?", *parentID)
	}

	if priorityLimit != nil && *priorityLimit > 0 {
		query = query.Where(c.CategoryPriority+" >= ?", *priorityLimit)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&categories).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch categories", "method", "GetAllCategories", "error", err, "include", include, "parentID", parentID, "priorityLimit", priorityLimit, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return categories, err
}

func (r *CategoryRepository) GetAllCategoriesPaginated(isActive *bool, showDeleted *bool, include []string, parentID *uint, page int, pageSize int, priorityLimit *int, sortBy, sortOrder string) ([]Category, int, error) {
	var categories []Category
	var total int64

	var maxPriorityLimit int = c.MaxPriorityLimit
	if priorityLimit != nil && *priorityLimit > maxPriorityLimit {
		priorityLimit = &maxPriorityLimit
	}

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if isActive != nil {
		query = query.Where(c.CategoryIsActive+" = ?", *isActive)
	}

	if parentID != nil {
		query = query.Where(c.CategoryParentID+" = ?", *parentID)
	}

	if priorityLimit != nil && *priorityLimit > 0 {
		query = query.Where(c.CategoryPriority+" >= ?", *priorityLimit)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Model(&Category{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count categories", "method", "GetAllCategoriesPaginated", "error", err, "include", include, "parentID", parentID, "page", page, "pageSize", pageSize, "priorityLimit", priorityLimit, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&categories).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch categories paginated", "method", "GetAllCategoriesPaginated", "error", err, "include", include, "parentID", parentID, "page", page, "pageSize", pageSize, "priorityLimit", priorityLimit, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return categories, int(total), err
}

func (r *CategoryRepository) GetCategoryByID(id uint, isActive *bool, showDeleted *bool) (*Category, error) {
	var category Category
	query := r.db.Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if isActive != nil {
		query = query.Where(c.CategoryIsActive+" = ?", *isActive)
	}

	if err := query.First(&category).Error; err != nil {
		logger.Logger.Error("Failed to fetch category by ID", "method", "GetCategoryByID", "error", err, "id", id, "isActive", isActive)
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(category *Category) (*Category, error) {
	err := r.db.Create(category).Error
	if err != nil {
		logger.Logger.Error("Failed to create category", "method", "CreateCategory", "error", err, "category", category)
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) UpdateCategory(category *Category, newIsActive *bool) (*Category, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(category).Error; err != nil {
			return err
		}

		if newIsActive != nil && *newIsActive != category.IsActive {
			return r.updateSubCategoriesIsActiveRecursively(tx, category.ID, *newIsActive)
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error("Failed to update category", "method", "UpdateCategory", "error", err, "category", category)
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) UndoDeletedCategory(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Update(c.FieldDeletedAt, nil).Where(c.FieldID+" = ?", id).Error; err != nil {
			logger.Logger.Error("Failed to delete category", "method", "DeleteCategory", "error", err, "id", id)
			return err
		}

		if err := tx.Unscoped().Where(c.CategoryParentID+" = ?", id).Update(c.FieldDeletedAt, nil).Error; err != nil {
			logger.Logger.Error("Failed to delete subcategories", "method", "DeleteCategory", "error", err, "id", id)
			return err
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error("Failed to delete category transaction", "method", "DeleteCategory", "error", err, "id", id)
	}
	return err
}

func (r *CategoryRepository) DeleteCategory(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Category{}, id).Error; err != nil {
			logger.Logger.Error("Failed to delete category", "method", "DeleteCategory", "error", err, "id", id)
			return err
		}

		if err := tx.Where(c.CategoryParentID+" = ?", id).Delete(&Category{}).Error; err != nil {
			logger.Logger.Error("Failed to delete subcategories", "method", "DeleteCategory", "error", err, "id", id)
			return err
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error("Failed to delete category transaction", "method", "DeleteCategory", "error", err, "id", id)
	}
	return err
}

func (r *CategoryRepository) ToggleCategoryIsActive(id uint, isActive bool) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Category{}).Where(c.FieldID+" = ?", id).Update(c.CategoryIsActive, isActive).Error; err != nil {
			return err
		}

		return r.updateSubCategoriesIsActiveRecursively(tx, id, isActive)
	})
	if err != nil {
		logger.Logger.Error("Failed to update category status", "method", "UpdateCategoryStatus", "error", err, "id", id, "isActive", isActive)
	}
	return err
}

func (r *CategoryRepository) updateSubCategoriesIsActiveRecursively(tx *gorm.DB, parentID uint, isActive bool) error {
	if err := tx.Where(c.CategoryParentID+" = ?", parentID).Update(c.CategoryIsActive, isActive).Error; err != nil {
		logger.Logger.Error("Failed to update subcategories is active", "method", "updateSubCategoriesIsActiveRecursively", "error", err, "parentID", parentID, "isActive", isActive)
		return err
	}

	return nil
}

func (r *CategoryRepository) CategoryExists(id uint, showDeleted *bool) (bool, error) {
	var category Category
	query := r.db.Model(&Category{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&category).Error
	if err != nil {
		logger.Logger.Error("Failed to check if category exists", "method", "CategoryExists", "error", err, "id", id)
		return false, err
	}

	return category.ID != 0, nil
}

func (r *CategoryRepository) CategoryExistsByTitle(title string, excludeID *uint) (bool, error) {
	var category Category
	query := r.db.Model(&Category{}).Where(c.CategoryTitle+" = ?", title)

	if excludeID != nil {
		query = query.Where(c.FieldID+" != ?", *excludeID)
	}

	err := query.Select(c.FieldID).Take(&category).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Failed to check if category exists by title", "method", "CategoryExistsByTitle", "error", err, "title", title, "excludeID", excludeID)
		return false, err
	}

	return category.ID != 0, nil
}

func (r *CategoryRepository) HasChildren(parentID uint) (bool, error) {
	var category Category
	err := r.db.Model(&Category{}).Where(c.CategoryParentID+" = ?", parentID).Select(c.FieldID).Take(&category).Error

	if err != nil {
		logger.Logger.Error("Failed to check if category has children", "method", "HasChildren", "error", err, "parentID", parentID)
		return false, err
	}

	return category.ID != 0, nil
}

func (r *CategoryRepository) IncrementPriority(categoryID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		rootID, err := r.findRootCategoryID(tx, categoryID)
		if err != nil {
			return err
		}

		result := tx.Model(&Category{}).
			Where(c.FieldID+" = ?", rootID).
			Update(c.CategoryPriority, gorm.Expr(c.CategoryPriority+" + 1"))

		return result.Error
	})
	if err != nil {
		logger.Logger.Error("Failed to increment priority", "method", "IncrementPriority", "error", err, "categoryID", categoryID)
	}
	return err
}

func (r *CategoryRepository) GetAllCategoriesWithSubcategories(subcategoryDepth *int, isActive *bool, showDeleted *bool, include []string, showPriority *bool, sortBy, sortOrder string) ([]CategorySubcategoriesResponse, error) {
	var categories []Category

	query := `
        WITH RECURSIVE category_tree AS (
            SELECT 
                c.*, 
                1 AS level
            FROM categories c
            WHERE c.parent_id IS NULL
            AND (?::boolean IS NULL OR c.is_active = ?)
            AND (?::boolean IS NULL OR ?::boolean OR c.deleted_at IS NULL)
            
            UNION ALL
            
            SELECT 
                c.*,
                ct.level + 1
            FROM categories c
            JOIN category_tree ct ON c.parent_id = ct.id
            WHERE ct.level < ?
            AND (?::boolean IS NULL OR c.is_active = ?)
            AND (?::boolean IS NULL OR ?::boolean OR c.deleted_at IS NULL)
        )
        SELECT * FROM category_tree
        ORDER BY level, 
            CASE WHEN ? = 'id' AND ? = 'asc' THEN id END ASC,
            CASE WHEN ? = 'id' AND ? = 'desc' THEN id END DESC,
            CASE WHEN ? = 'title' AND ? = 'asc' THEN title END ASC,
            CASE WHEN ? = 'title' AND ? = 'desc' THEN title END DESC,
            CASE WHEN ? = 'created_at' AND ? = 'asc' THEN created_at END ASC,
            CASE WHEN ? = 'created_at' AND ? = 'desc' THEN created_at END DESC
    `

	var subcategoryDepthInt int = c.SubcategoryDepthLimit
	if subcategoryDepth == nil || *subcategoryDepth > c.SubcategoryDepthLimit {
		subcategoryDepth = &subcategoryDepthInt
	}

	err := r.db.Raw(query,
		// Initial WHERE
		isActive, isActive, // is_active status (NULL + value)
		showDeleted, showDeleted, // deleted_at (NULL + showDeleted)

		// UNION ALL WHERE
		subcategoryDepth,   // subcategoryDepth
		isActive, isActive, // is_active status (NULL + value)
		showDeleted, showDeleted, // deleted_at (NULL + showDeleted)

		// ORDER BY
		sortBy, sortOrder, // id ASC
		sortBy, sortOrder, // id DESC
		sortBy, sortOrder, // title ASC
		sortBy, sortOrder, // title DESC
		sortBy, sortOrder, // created_at ASC
		sortBy, sortOrder, // created_at DESC
	).Scan(&categories).Error

	if err != nil {
		logger.Logger.Error("Failed to fetch nested categories", "method", "GetNestedCategories", "error", err)
		return nil, err
	}

	lookup := make(map[uint]*Category)
	for i := range categories {
		categories[i].ImageURL = utils.BuildFullImageURL(categories[i].PathKey)
		lookup[categories[i].ID] = &categories[i]
	}

	var buildSubcategories func(parentID uint) []CategorySubcategoriesResponse
	buildSubcategories = func(parentID uint) []CategorySubcategoriesResponse {
		var subs []CategorySubcategoriesResponse
		for _, cat := range categories {
			if cat.ParentID != nil && *cat.ParentID == parentID {
				sub := CategorySubcategoriesResponse{
					Category:      *lookup[cat.ID].ToResponse(),
					Subcategories: buildSubcategories(cat.ID),
				}
				subs = append(subs, sub)
			}
		}
		return subs
	}

	var roots []CategorySubcategoriesResponse
	for i := range categories {
		if categories[i].ParentID == nil {
			root := CategorySubcategoriesResponse{
				Category:      *lookup[categories[i].ID].ToResponse(),
				Subcategories: buildSubcategories(categories[i].ID),
			}
			roots = append(roots, root)
		}
	}

	return roots, nil
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
		err := tx.Select(c.FieldID, c.CategoryParentID).
			Where(c.FieldID+" = ?", currentID).
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
	defaultFields := []string{c.FieldID, c.CategoryTitle, c.CategoryIsActive, c.FieldCreatedAt, c.FieldUpdatedAt}
	optionalFields := []string{c.CategorySubTitle, c.CategoryImageURL, c.CategoryParentID, c.CategoryPriority}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
