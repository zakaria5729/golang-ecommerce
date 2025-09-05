package category

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type CategoryUseCase struct {
	repo *CategoryRepository
}

func NewCategoryUseCase() *CategoryUseCase {
	return &CategoryUseCase{
		repo: NewCategoryRepository(),
	}
}

func (uc *CategoryUseCase) GetAllCategories(includeStr string, parentIDFilter string, priorityFilter string, sortBy, sortOrder string) ([]Category, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	parentID, _ := utils.ParseUint(parentIDFilter)
	showPriority := utils.ParseBoolPtr(priorityFilter)

	categories, err := uc.repo.GetAllCategories(include, parentID, showPriority, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories", "method", "GetAllCategories", "error", err, "include", include, "parentID", parentID, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return categories, nil
}

func (uc *CategoryUseCase) GetAllCategoriesPaginated(includeStr string, parentIDFilter string, pageStr string, pageSizeStr string, priorityFilter string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	include := utils.ParseCommaSeparatedString(includeStr)
	parentID, _ := utils.ParseUint(parentIDFilter)
	showPriority := utils.ParseBoolPtr(priorityFilter)

	categories, total, err := uc.repo.GetAllCategoriesPaginated(include, parentID, page, pageSize, showPriority, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories paginated", "method", "GetAllCategoriesPaginated", "error", err, "include", include, "parentID", parentID, "page", page, "pageSize", pageSize, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	var categoryPtrs []*Category
	for i := range categories {
		categoryPtrs = append(categoryPtrs, &categories[i])
	}

	return utils.BuildPaginatedResponse(categoryPtrs, total, page, pageSize), nil
}

func (uc *CategoryUseCase) GetCategoryByID(id uint, includeStr string) (*Category, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	category, err := uc.repo.GetCategoryByID(id, include)
	if err != nil {
		logger.Logger.Error("Category not found", "method", "GetCategoryByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("category not found: %w", err)
	}

	return category, nil
}

func (uc *CategoryUseCase) CreateCategory(req *Category) (*Category, error) {

	req.Sanitize()

	if req.ParentID != nil && *req.ParentID > 0 {
		parentExists, err := uc.repo.CategoryExists(*req.ParentID)
		if err != nil {
			logger.Logger.Error("Failed to check parent category", "method", "CreateCategory", "error", err, "parentID", *req.ParentID)
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			logger.Logger.Error("Parent category does not exist", "method", "CreateCategory", "parentID", *req.ParentID)
			return nil, errors.New("parent category does not exist")
		}
	}

	exists, err := uc.repo.CategoryExistsByTitle(req.Title, nil)
	if err != nil {
		logger.Logger.Error("Failed to check category title", "method", "CreateCategory", "error", err, "title", req.Title)
		return nil, fmt.Errorf("failed to check category title: %w", err)
	}
	if exists {
		logger.Logger.Error("Category with this title already exists", "method", "CreateCategory", "title", req.Title)
		return nil, errors.New("category with this title already exists")
	}

	if req.IsActive {
		req.IsActive = true
	}

	category := &Category{
		Title:    req.Title,
		SubTitle: req.SubTitle,
		ImageURL: req.ImageURL,
		ParentID: req.ParentID,
		IsActive: req.IsActive,
	}

	if err := uc.repo.CreateCategory(category); err != nil {
		logger.Logger.Error("Failed to create category", "method", "CreateCategory", "error", err, "category", category)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (uc *CategoryUseCase) UpdateCategory(id uint, req *Category) (*Category, error) {

	req.Sanitize()

	existingCategory, err := uc.repo.GetCategoryByIDIncludeInactive(id, nil)
	if err != nil {
		logger.Logger.Error("Category not found", "method", "UpdateCategory", "error", err, "id", id)
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if req.ParentID != nil && *req.ParentID > 0 {
		if *req.ParentID == id {
			logger.Logger.Error("Category cannot be its own parent", "method", "UpdateCategory", "id", id, "parentID", *req.ParentID)
			return nil, errors.New("category cannot be its own parent")
		}

		parentExists, err := uc.repo.CategoryExists(*req.ParentID)
		if err != nil {
			logger.Logger.Error("Failed to check parent category", "method", "UpdateCategory", "error", err, "parentID", *req.ParentID)
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			logger.Logger.Error("Parent category does not exist", "method", "UpdateCategory", "parentID", *req.ParentID)
			return nil, errors.New("parent category does not exist")
		}
	}

	if req.Title != "" && req.Title != existingCategory.Title {
		exists, err := uc.repo.CategoryExistsByTitle(req.Title, &id)
		if err != nil {
			logger.Logger.Error("Failed to check category title", "method", "UpdateCategory", "error", err, "title", req.Title, "id", id)
			return nil, fmt.Errorf("failed to check category title: %w", err)
		}
		if exists {
			logger.Logger.Error("Category with this title already exists", "method", "UpdateCategory", "title", req.Title, "id", id)
			return nil, errors.New("category with this title already exists")
		}
	}

	if req.Title != "" {
		existingCategory.Title = req.Title
	}
	if req.SubTitle != nil {
		existingCategory.SubTitle = req.SubTitle
	}
	if req.ImageURL != nil {
		existingCategory.ImageURL = req.ImageURL
	}
	if req.ParentID != nil {
		existingCategory.ParentID = req.ParentID
	}

	if req.IsActive != existingCategory.IsActive {
		if err := uc.repo.UpdateCategoryStatus(id, req.IsActive); err != nil {
			logger.Logger.Error("Failed to update category status", "method", "UpdateCategory", "error", err, "id", id, "newStatus", req.IsActive)
			return nil, fmt.Errorf("failed to update category status: %w", err)
		}
		existingCategory.IsActive = req.IsActive
	}

	if err := uc.repo.UpdateCategory(existingCategory); err != nil {
		logger.Logger.Error("Failed to update category", "method", "UpdateCategory", "error", err, "category", existingCategory)
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return existingCategory, nil
}

func (uc *CategoryUseCase) DeleteCategory(id uint) error {
	_, err := uc.repo.GetCategoryByIDIncludeInactive(id, nil)
	if err != nil {
		logger.Logger.Error("Category not found", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("category not found: %w", err)
	}

	if err := uc.repo.DeleteCategory(id); err != nil {
		logger.Logger.Error("Failed to delete category", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (uc *CategoryUseCase) IncrementPriority(categoryID uint) error {
	if err := uc.repo.IncrementPriority(categoryID); err != nil {
		logger.Logger.Error("Failed to increment priority for category", "method", "IncrementPriority", "error", err, "categoryID", categoryID)
		return err
	}
	return nil
}

func (uc *CategoryUseCase) ToggleCategoryStatus(id uint) (*Category, error) {
	category, err := uc.repo.GetCategoryByIDIncludeInactive(id, nil)
	if err != nil {
		logger.Logger.Error("Category not found", "method", "ToggleCategoryStatus", "error", err, "id", id)
		return nil, fmt.Errorf("category not found: %w", err)
	}

	newStatus := !category.IsActive
	if err := uc.repo.UpdateCategoryStatus(id, newStatus); err != nil {
		logger.Logger.Error("Failed to update category status", "method", "ToggleCategoryStatus", "error", err, "id", id, "newStatus", newStatus)
		return nil, fmt.Errorf("failed to update category status: %w", err)
	}

	category.IsActive = newStatus
	return category, nil
}
