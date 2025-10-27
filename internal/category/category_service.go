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

type CategoryService interface {
	GetAllCategoriesWithSubcategories(showDeleted *bool, subcategoryDepth *int, sortBy, sortOrder string) ([]m.CategorySubcategoriesResponse, error)
	GetAllCategories(showDeleted *bool, parentID *uint, priorityLimit *int, sortBy, sortOrder string) ([]m.CategoryResponse, error)
	GetAllCategoriesPaginated(showDeleted *bool, parentID *uint, pageStr string, pageSizeStr string, priorityLimit *int, sortBy, sortOrder string) (*r.PaginatedResponse, error)
	GetCategoryByID(id uint, showDeleted *bool) (*m.CategoryResponse, error)
	CreateCategory(ctx context.Context, req *m.CreateCategoryRequest) (*m.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req *m.UpdateCategoryRequest) (*m.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uint) error
	UndoDeleteCategory(ctx context.Context, id uint) error
}

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) GetAllCategoriesWithSubcategories(showDeleted *bool, subcategoryDepth *int, sortBy, sortOrder string) ([]m.CategorySubcategoriesResponse, error) {
	categories, err := s.repo.GetAllCategoriesWithSubcategories(subcategoryDepth, showDeleted, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nested categories: %w", err)
	}

	return categories, nil
}

func (s *categoryService) GetAllCategories(showDeleted *bool, parentID *uint, priorityLimit *int, sortBy, sortOrder string) ([]m.CategoryResponse, error) {
	categories, err := s.repo.GetAllCategories(showDeleted, parentID, priorityLimit, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return getCategoryResponses(categories), nil
}

func (s *categoryService) GetAllCategoriesPaginated(showDeleted *bool, parentID *uint, pageStr string, pageSizeStr string, priorityLimit *int, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	categories, total, err := s.repo.GetAllCategoriesPaginated(showDeleted, parentID, page, pageSize, priorityLimit, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return utils.BuildPaginatedResponse(getCategoryResponses(categories), total, page, pageSize), nil
}

func (s *categoryService) GetCategoryByID(id uint, showDeleted *bool) (*m.CategoryResponse, error) {
	category, err := s.repo.GetCategoryByID(id, showDeleted)

	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	return category.ToResponse(), nil
}

func (s *categoryService) CreateCategory(ctx context.Context, req *m.CreateCategoryRequest) (*m.CategoryResponse, error) {
	category := &CategoryEntity{
		Title:    req.Title,
		SubTitle: req.SubTitle,
		ParentID: req.ParentID,
		Priority: req.Priority,
		PathKey:  &req.PathKey,
	}

	category.CreatedBy, _ = cu.GetUserIDFromContext(ctx)
	category.Sanitize()

	if req.ParentID != nil && *req.ParentID != 0 {
		parentExists, err := s.repo.CategoryExists(*req.ParentID, nil)
		if err != nil || !parentExists {
			return nil, fmt.Errorf("parent category not found with ID: %d", *req.ParentID)
		}

		parentCategory, err := s.repo.GetCategoryByID(*req.ParentID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent category: %w", err)
		}

		if parentCategory.ParentID != nil {
			return nil, errors.New("nested subcategories are not allowed")
		}

		category.PathKey = parentCategory.PathKey
	}

	titleExists, err := s.repo.CategoryExistsByTitle(req.Title, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check if category exists: %w", err)
	}
	if titleExists {
		return nil, fmt.Errorf("category with title '%s' already exists", req.Title)
	}

	createdCategory, err := s.repo.CreateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return createdCategory.ToResponse(), nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uint, req *m.UpdateCategoryRequest) (*m.CategoryResponse, error) {
	existingCategory, err := s.repo.GetCategoryByID(id, nil)
	if err != nil || existingCategory == nil {
		return nil, fmt.Errorf("category not found with category ID: %w", err)
	}

	if req.Title != "" {
		exists, err := s.repo.CategoryExistsByTitle(req.Title, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check if category exists: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("category with title '%s' already exists", req.Title)
		}
		existingCategory.Title = req.Title
	}

	if req.SubTitle != nil {
		existingCategory.SubTitle = req.SubTitle
	}

	if req.ParentID != nil {
		if *req.ParentID == 0 {
			existingCategory.ParentID = nil
		} else {
			parentExists, err := s.repo.CategoryExists(*req.ParentID, nil)
			if err != nil || !parentExists {
				return nil, fmt.Errorf("parent category not found with ID: %d", *req.ParentID)
			}

			parentCategory, err := s.repo.GetCategoryByID(*req.ParentID, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get parent category: %w", err)
			}

			if parentCategory.ParentID != nil {
				return nil, errors.New("nested subcategories are not allowed")
			}

			existingCategory.ParentID = req.ParentID
			existingCategory.PathKey = parentCategory.PathKey
		}
	}

	if req.Priority != nil {
		tempPriority := *req.Priority
		existingCategory.Priority = &tempPriority
	}

	existingCategory.UpdatedBy, _ = cu.GetUserIDFromContext(ctx)
	existingCategory.Sanitize()

	updatedCategory, err := s.repo.UpdateCategory(existingCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return updatedCategory.ToResponse(), nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	exists, err := s.repo.CategoryExists(id, nil)
	if err != nil || !exists {
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := s.repo.DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (s *categoryService) UndoDeleteCategory(ctx context.Context, id uint) error {
	showDeleted := true
	exists, err := s.repo.CategoryExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("category not found with category ID: %w", err)
	}

	if err := s.repo.UndoDeletedCategory(ctx, id); err != nil {
		return fmt.Errorf("failed to undo delete category: %w", err)
	}

	return nil
}

func (s *categoryService) IncrementPriority(categoryID uint) error {
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
