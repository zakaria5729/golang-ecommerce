package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	m "github.com/easy-comerce/backend/pkg/middleware"
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

func (uc *CategoryUseCase) GetAllCategoriesWithSubcategories(showDeleted *bool, subcategoryDepthFilter string, sortBy, sortOrder string) ([]CategorySubcategoriesResponse, error) {
	subcategoryDepth, _ := utils.ParseInt(subcategoryDepthFilter)
	categories, err := uc.repo.GetAllCategoriesWithSubcategories(subcategoryDepth, showDeleted, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch nested categories", "method", "GetAllCategoriesWithSubcategories", "error", err)
		return nil, fmt.Errorf("failed to fetch nested categories: %w", err)
	}

	return categories, nil
}

func (uc *CategoryUseCase) GetAllCategories(showDeleted *bool, parentIDFilter string, priorityLimitFilter string, sortBy, sortOrder string) ([]CategoryResponse, error) {
	parentID, _ := utils.ParseUint(parentIDFilter)
	priorityLimit, _ := utils.ParseInt(priorityLimitFilter)

	categories, err := uc.repo.GetAllCategories(showDeleted, parentID, priorityLimit, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories", "method", "GetAllCategories", "error", err, "parentID", parentID, "priorityLimit", priorityLimit, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return uc.getCategoryResponses(categories), nil
}

func (uc *CategoryUseCase) GetAllCategoriesPaginated(showDeleted *bool, parentIDFilter string, pageStr string, pageSizeStr string, priorityLimitFilter string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	parentID, _ := utils.ParseUint(parentIDFilter)
	priorityLimit, _ := utils.ParseInt(priorityLimitFilter)

	categories, total, err := uc.repo.GetAllCategoriesPaginated(showDeleted, parentID, page, pageSize, priorityLimit, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories paginated", "method", "GetAllCategoriesPaginated", "error", err, "parentID", parentID, "page", page, "pageSize", pageSize, "priorityLimit", priorityLimit, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return utils.BuildPaginatedResponse(uc.getCategoryResponses(categories), total, page, pageSize), nil
}

func (uc *CategoryUseCase) GetCategoryByID(id uint, showDeleted *bool) (*CategoryResponse, error) {
	category, err := uc.repo.GetCategoryByID(id, showDeleted)

	if err != nil {
		logger.Logger.Error("Category not found", "method", "GetCategoryByID", "error", err, "id", id)
		return nil, fmt.Errorf("category not found: %w", err)
	}

	return category.ToResponse(), nil
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	category := &Category{
		Title:    req.Title,
		SubTitle: req.SubTitle,
		ParentID: req.ParentID,
		Priority: req.Priority,
		PathKey:  &req.PathKey,
	}

	category.CreatedBy = m.GetUserIdOnlyFromContext(ctx)
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

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, id uint, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	existingCategory, err := uc.repo.GetCategoryByID(id, nil)
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
	if req.PathKey != "" {
		existingCategory.PathKey = &req.PathKey
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

	existingCategory.UpdatedBy = m.GetUserIdOnlyFromContext(ctx)
	updatedCategory, err := uc.repo.UpdateCategory(existingCategory)
	if err != nil {
		logger.Logger.Error("Failed to update category", "method", "UpdateCategory", "error", err, "category", existingCategory)
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return updatedCategory.ToResponse(), nil
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id uint) error {
	exists, err := uc.repo.CategoryExists(id, nil)
	if err != nil || !exists {
		logger.Logger.Error("Category not found", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := uc.repo.DeleteCategory(ctx, id); err != nil {
		logger.Logger.Error("Failed to delete category", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (uc *CategoryUseCase) UndoDeletedCategory(ctx context.Context, id uint) error {
	showDeleted := true
	exists, err := uc.repo.CategoryExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Category not found", "method", "DeleteCategory", "error", err, "id", id)
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := uc.repo.UndoDeletedCategory(ctx, id); err != nil {
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

func (uc *CategoryUseCase) getCategoryResponses(categories []Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = *cat.ToResponse()
	}
	return responses
}
