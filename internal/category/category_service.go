package category

import (
	"context"
	"errors"
	"fmt"

	m "github.com/easy-comerce/backend/internal/category/model"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	r "github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type CategoryService struct {
	repo *CategoryRepository
}

func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) GetAllCategoriesWithSubcategories(showDeleted *bool, subcategoryDepthFilter string, sortBy, sortOrder string) ([]m.CategorySubcategoriesResponse, error) {
	subcategoryDepth, _ := utils.ParseInt(subcategoryDepthFilter)
	categories, err := s.repo.GetAllCategoriesWithSubcategories(subcategoryDepth, showDeleted, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nested categories: %w", err)
	}

	return categories, nil
}

func (s *CategoryService) GetAllCategories(showDeleted *bool, parentIDFilter string, priorityLimitFilter string, sortBy, sortOrder string) ([]m.CategoryResponse, error) {
	parentID, _ := utils.ParseUint(parentIDFilter)
	priorityLimit, _ := utils.ParseInt(priorityLimitFilter)

	categories, err := s.repo.GetAllCategories(showDeleted, parentID, priorityLimit, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return getCategoryResponses(categories), nil
}

func (s *CategoryService) GetAllCategoriesPaginated(showDeleted *bool, parentIDFilter string, pageStr string, pageSizeStr string, priorityLimitFilter string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	parentID, _ := utils.ParseUint(parentIDFilter)
	priorityLimit, _ := utils.ParseInt(priorityLimitFilter)

	categories, total, err := s.repo.GetAllCategoriesPaginated(showDeleted, parentID, page, pageSize, priorityLimit, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return utils.BuildPaginatedResponse(getCategoryResponses(categories), total, page, pageSize), nil
}

func (s *CategoryService) GetCategoryByID(id uint, showDeleted *bool) (*m.CategoryResponse, error) {
	category, err := s.repo.GetCategoryByID(id, showDeleted)

	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	return category.ToResponse(), nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *m.CreateCategoryRequest) (*m.CategoryResponse, error) {
	category := &CategoryEntity{
		Title:    req.Title,
		SubTitle: req.SubTitle,
		ParentID: req.ParentID,
		Priority: req.Priority,
		PathKey:  &req.PathKey,
	}

	category.CreatedBy, _ = cu.GetUserIDFromContext(ctx)
	category.Sanitize()

	if category.ParentID != nil && *category.ParentID > 0 {
		parentExists, err := s.repo.CategoryExists(*category.ParentID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			return nil, errors.New("parent category does not exist")
		}
	}

	exists, err := s.repo.CategoryExistsByTitle(category.Title, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check category title: %w", err)
	}
	if exists {
		return nil, errors.New("category with this title already exists")
	}

	createdCategory, err := s.repo.CreateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return createdCategory.ToResponse(), nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, req *m.UpdateCategoryRequest) (*m.CategoryResponse, error) {
	existingCategory, err := s.repo.GetCategoryByID(id, nil)
	if err != nil || existingCategory == nil {
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
			return nil, errors.New("category cannot be its own parent")
		}

		parentExists, err := s.repo.CategoryExists(*existingCategory.ParentID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			return nil, errors.New("parent category does not exist")
		}
	}

	if req.Title != "" && req.Title != existingCategory.Title {
		exists, err := s.repo.CategoryExistsByTitle(req.Title, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check category title: %w", err)
		}
		if exists {
			return nil, errors.New("category with this title already exists")
		}
	}

	existingCategory.UpdatedBy, _ = cu.GetUserIDFromContext(ctx)
	updatedCategory, err := s.repo.UpdateCategory(existingCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return updatedCategory.ToResponse(), nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	exists, err := s.repo.CategoryExists(id, nil)
	if err != nil || !exists {
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := s.repo.DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (s *CategoryService) UndoDeletedCategory(ctx context.Context, id uint) error {
	showDeleted := true
	exists, err := s.repo.CategoryExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := s.repo.UndoDeletedCategory(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (s *CategoryService) IncrementPriority(categoryID uint) error {
	if err := s.repo.IncrementPriority(categoryID); err != nil {
		return err
	}
	return nil
}

func getCategoryResponses(categories []CategoryEntity) []m.CategoryResponse {
	responses := make([]m.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = *cat.ToResponse()
	}
	return responses
}
