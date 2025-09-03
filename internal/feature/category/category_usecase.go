package category

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
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
		logger.Logger.Error("Failed to fetch categories", "error", err, "include", include, "parentID", parentID, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
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
		logger.Logger.Error("Failed to fetch categories paginated", "error", err, "include", include, "parentID", parentID, "page", page, "pageSize", pageSize, "showPriority", showPriority, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	var categoryPtrs []*Category
	for i := range categories {
		categoryPtrs = append(categoryPtrs, &categories[i])
	}

	return utils.BuildPaginatedResponse(categoryPtrs, total, page, pageSize), nil
}

func (uc *CategoryUseCase) GetCategoryByID(id uint, includeStr string) (*Category, error) {
	if id == 0 {
		logger.Logger.Error("Invalid category ID provided", "id", id)
		return nil, errors.New("invalid category ID")
	}

	include := utils.ParseCommaSeparatedString(includeStr)
	category, err := uc.repo.GetCategoryByID(id, include)
	if err != nil {
		logger.Logger.Error("Category not found", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("category not found: %w", err)
	}

	return category, nil
}

func (uc *CategoryUseCase) CreateCategory(req *Category) (*Category, error) {
	if err := uc.validateCreateRequest(req); err.HasErrors() {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	req.Sanitize()

	if req.ParentID != nil && *req.ParentID > 0 {
		parentExists, err := uc.repo.CategoryExists(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			return nil, errors.New("parent category does not exist")
		}
	}

	exists, err := uc.repo.CategoryExistsByTitle(req.Title, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check category title: %w", err)
	}
	if exists {
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
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (uc *CategoryUseCase) UpdateCategory(id uint, req *Category) (*Category, error) {
	if id == 0 {
		return nil, errors.New("invalid category ID")
	}

	if err := uc.validateUpdateRequest(req); err.HasErrors() {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	req.Sanitize()

	existingCategory, err := uc.repo.GetCategoryByIDIncludeInactive(id, nil)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if req.ParentID != nil && *req.ParentID > 0 {
		if *req.ParentID == id {
			return nil, errors.New("category cannot be its own parent")
		}

		parentExists, err := uc.repo.CategoryExists(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent category: %w", err)
		}
		if !parentExists {
			return nil, errors.New("parent category does not exist")
		}
	}

	if req.Title != "" && req.Title != existingCategory.Title {
		exists, err := uc.repo.CategoryExistsByTitle(req.Title, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check category title: %w", err)
		}
		if exists {
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
			return nil, fmt.Errorf("failed to update category status: %w", err)
		}
		existingCategory.IsActive = req.IsActive
	}

	if err := uc.repo.UpdateCategory(existingCategory); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return existingCategory, nil
}

func (uc *CategoryUseCase) DeleteCategory(id uint) error {
	if id == 0 {
		return errors.New("invalid category ID")
	}

	_, err := uc.repo.GetCategoryByIDIncludeInactive(id, nil)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	if err := uc.repo.DeleteCategory(id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (uc *CategoryUseCase) IncrementPriority(categoryID uint) error {
	if err := uc.repo.IncrementPriority(categoryID); err != nil {
		logger.Logger.Error("Failed to increment priority for category", "categoryID", categoryID, "error", err)
		return nil
	}
	return nil
}

func (uc *CategoryUseCase) ToggleCategoryStatus(id uint) (*Category, error) {
	if id == 0 {
		return nil, errors.New("invalid category ID")
	}

	category, err := uc.repo.GetCategoryByIDIncludeInactive(id, nil)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	newStatus := !category.IsActive

	if err := uc.repo.UpdateCategoryStatus(id, newStatus); err != nil {
		return nil, fmt.Errorf("failed to update category status: %w", err)
	}

	category.IsActive = newStatus
	return category, nil
}

func (uc *CategoryUseCase) validateCreateRequest(req *Category) validator.ValidationErrors {
	var errors validator.ValidationErrors

	errors = validator.MergeValidationErrors(
		errors,
		validator.ValidateRequired(req.Title, "title"),
	)

	if req.Title != "" {
		errors = validator.MergeValidationErrors(
			errors,
			validator.ValidateMinLength(req.Title, "title", 2),
			validator.ValidateMaxLength(req.Title, "title", 100),
		)
	}

	if req.SubTitle != nil && *req.SubTitle != "" {
		errors = validator.MergeValidationErrors(
			errors,
			validator.ValidateMaxLength(*req.SubTitle, "sub_title", 200),
		)
	}

	if req.ImageURL != nil && *req.ImageURL != "" {
		errors = validator.MergeValidationErrors(
			errors,
			validator.ValidateURL(*req.ImageURL, "image_url"),
		)
	}

	if req.ParentID != nil && *req.ParentID == 0 {
		errors.AddError("parent_id", "parent_id must be a positive integer")
	}

	return errors
}

func (uc *CategoryUseCase) validateUpdateRequest(req *Category) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if req.Title != "" {
		errors = validator.MergeValidationErrors(
			errors,
			validator.ValidateMinLength(req.Title, "title", 2),
			validator.ValidateMaxLength(req.Title, "title", 120),
		)
	}

	if req.SubTitle != nil && *req.SubTitle != "" {
		errors = validator.MergeValidationErrors(
			errors,
			validator.ValidateMaxLength(*req.SubTitle, "sub_title", 255),
		)
	}

	if req.ImageURL != nil && *req.ImageURL != "" {
		errors = validator.MergeValidationErrors(
			errors,
			validator.ValidateURL(*req.ImageURL, "image_url"),
		)
	}

	if req.ParentID != nil && *req.ParentID == 0 {
		errors.AddError("parent_id", "parent_id must be a positive integer")
	}

	return errors
}
