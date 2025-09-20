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

func (uc *CategoryUseCase) GetAllCategories(isActive *bool, showDeleted *bool, includeStr string, parentIDFilter string, priorityFilter string, sortBy, sortOrder string) ([]CategoryResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	parentID, _ := utils.ParseUint(parentIDFilter)
	showPriority := utils.ParseBoolPtr(priorityFilter)

	categories, err := uc.repo.GetAllCategories(isActive, showDeleted, include, parentID, showPriority, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories", "method", "GetAllCategories", "error", err, "include", include, "parentID", parentID, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return uc.getCategoryResponses(categories), nil
}

func (uc *CategoryUseCase) GetAllCategoriesPaginated(isActive *bool, showDeleted *bool, includeStr string, parentIDFilter string, pageStr string, pageSizeStr string, priorityFilter string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)
	parentID, _ := utils.ParseUint(parentIDFilter)
	showPriority := utils.ParseBoolPtr(priorityFilter)

	categories, total, err := uc.repo.GetAllCategoriesPaginated(isActive, showDeleted, include, parentID, page, pageSize, showPriority, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories paginated", "method", "GetAllCategoriesPaginated", "error", err, "include", include, "parentID", parentID, "page", page, "pageSize", pageSize, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return utils.BuildPaginatedResponse(uc.getCategoryResponses(categories), total, page, pageSize), nil
}

func (uc *CategoryUseCase) GetCategoryByID(id uint, isActive *bool, showDeleted *bool) (*CategoryResponse, error) {
	category, err := uc.repo.GetCategoryByID(id, isActive, showDeleted)

	if err != nil {
		logger.Logger.Error("Category not found", "method", "GetCategoryByID", "error", err, "id", id)
		return nil, fmt.Errorf("category not found: %w", err)
	}

	return category.ToResponse(), nil
}

func (uc *CategoryUseCase) CreateCategory(req *CreateCategoryRequest) (*CategoryResponse, error) {
	category := &Category{
		Title:    req.Title,
		SubTitle: req.SubTitle,
		ParentID: req.ParentID,
		Priority: req.Priority,
		IsActive: true,
	}
	category.Sanitize()

	if category.ParentID != nil && *category.ParentID > 0 {
		parentExists, err := uc.repo.CategoryExists(*category.ParentID, nil)
		if err != nil {
			logger.Logger.Error("Failed to check parent category", "method", "CreateCategory", "error", err, "parentID", *category.ParentID)
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			logger.Logger.Error("Parent category does not exist", "method", "CreateCategory", "parentID", *category.ParentID)
			return nil, errors.New("parent category does not exist")
		}
	}

	exists, err := uc.repo.CategoryExistsByTitle(category.Title, nil)
	if err != nil {
		logger.Logger.Error("Failed to check category title", "method", "CreateCategory", "error", err, "title", category.Title)
		return nil, fmt.Errorf("failed to check category title: %w", err)
	}
	if exists {
		logger.Logger.Error("Category with this title already exists", "method", "CreateCategory", "title", category.Title)
		return nil, errors.New("category with this title already exists")
	}

	createdCategory, err := uc.repo.CreateCategory(category)
	if err != nil {
		logger.Logger.Error("Failed to create category", "method", "CreateCategory", "error", err, "category", category)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return createdCategory.ToResponse(), nil
}

func (uc *CategoryUseCase) UpdateCategory(id uint, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	existingCategory, err := uc.repo.GetCategoryByID(id, nil, nil)
	if err != nil || existingCategory == nil {
		logger.Logger.Error("Category not found", "method", "UpdateCategory", "error", err, "id", id)
		return nil, fmt.Errorf("category not found with category ID: %w", err)
	}

	if req.Title != "" {
		existingCategory.Title = req.Title
	}
	if req.SubTitle != nil {
		existingCategory.SubTitle = req.SubTitle
	}
	if req.ParentID != nil {
		existingCategory.ParentID = req.ParentID
	}
	if req.Priority != nil {
		existingCategory.Priority = req.Priority
	}
	if req.IsActive != nil {
		existingCategory.IsActive = *req.IsActive
	}
	existingCategory.Sanitize()

	if existingCategory.ParentID != nil && *existingCategory.ParentID > 0 {
		if *existingCategory.ParentID == id {
			logger.Logger.Error("Category cannot be its own parent", "method", "UpdateCategory", "id", id, "parentID", *existingCategory.ParentID)
			return nil, errors.New("category cannot be its own parent")
		}

		parentExists, err := uc.repo.CategoryExists(*existingCategory.ParentID, nil)
		if err != nil {
			logger.Logger.Error("Failed to check parent category", "method", "UpdateCategory", "error", err, "parentID", *existingCategory.ParentID)
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			logger.Logger.Error("Parent category does not exist", "method", "UpdateCategory", "parentID", *existingCategory.ParentID)
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

	updatedCategory, err := uc.repo.UpdateCategory(existingCategory, req.IsActive)
	if err != nil {
		logger.Logger.Error("Failed to update category", "method", "UpdateCategory", "error", err, "category", existingCategory)
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return updatedCategory.ToResponse(), nil
}

func (uc *CategoryUseCase) DeleteCategory(id uint) error {
	exists, err := uc.repo.CategoryExists(id, nil)
	if err != nil || !exists {
		logger.Logger.Error("Category not found", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := uc.repo.DeleteCategory(id); err != nil {
		logger.Logger.Error("Failed to delete category", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (uc *CategoryUseCase) UndoDeletedCategory(id uint) error {
	showDeleted := true
	exists, err := uc.repo.CategoryExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Category not found", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := uc.repo.UndoDeletedCategory(id); err != nil {
		logger.Logger.Error("Failed to undo delete category", "method", "UndoDeletedCategory", "error", err, "id", id)
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

func (uc *CategoryUseCase) ToggleCategoryIsActive(id uint) (*CategoryResponse, error) {
	category, err := uc.repo.GetCategoryByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Category not found", "method", "ToggleCategoryStatus", "error", err, "id", id)
		return nil, fmt.Errorf("category not found: %w", err)
	}

	newStatus := !category.IsActive
	if err := uc.repo.ToggleCategoryIsActive(id, newStatus); err != nil {
		logger.Logger.Error("Failed to update category status", "method", "ToggleCategoryStatus", "error", err, "id", id, "newStatus", newStatus)
		return nil, fmt.Errorf("failed to update category status: %w", err)
	}

	category.IsActive = newStatus
	return category.ToResponse(), nil
}

func (uc *CategoryUseCase) getCategoryResponses(categories []Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = *cat.ToResponse()
	}
	return responses
}
