package size_category

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SizeCategoryUseCase struct {
	sizeCategoryRepo *SizeCategoryRepository
}

func NewSizeCategoryUseCase() *SizeCategoryUseCase {
	return &SizeCategoryUseCase{
		sizeCategoryRepo: NewSizeCategoryRepository(),
	}
}

func (uc *SizeCategoryUseCase) GetAllSizeCategories(includeStr string, showDeleted *bool, sortBy, sortOrder string) ([]SizeCategory, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	sizeCategories, err := uc.sizeCategoryRepo.GetAllSizeCategories(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch size categories", "method", "GetAllSizeCategories", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch size categories: %w", err)
	}

	return sizeCategories, nil
}

func (uc *SizeCategoryUseCase) GetSizeCategoryByID(id uint, includeStr string, showDeleted *bool) (*SizeCategory, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	sizeCategory, err := uc.sizeCategoryRepo.GetSizeCategoryByID(id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Size category not found", "method", "GetSizeCategoryByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("size category not found: %w", err)
	}

	return sizeCategory, nil
}

func (uc *SizeCategoryUseCase) CreateSizeCategory(req *CreateSizeCategoryRequest) (*SizeCategory, error) {
	req.Sanitize()

	exists, err := uc.sizeCategoryRepo.SizeCategoryExistsByName(req.Name)
	if err != nil {
		logger.Logger.Error("Failed to check if size category exists", "method", "CreateSizeCategory", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to check size category existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Size category already exists", "method", "CreateSizeCategory", "name", req.Name)
		return nil, errors.New("size category with this name already exists")
	}

	sizeCategory := &SizeCategory{
		Name: req.Name,
	}

	if sizeCategory, err := uc.sizeCategoryRepo.CreateSizeCategory(sizeCategory); err != nil {
		logger.Logger.Error("Failed to create size category", "method", "CreateSizeCategory", "error", err, "sizeCategory", sizeCategory)
		return nil, fmt.Errorf("failed to create size category: %w", err)
	}

	return sizeCategory, nil
}

func (uc *SizeCategoryUseCase) UpdateSizeCategory(id uint, req *UpdateSizeCategoryRequest) (*SizeCategory, error) {
	req.Sanitize()

	existingSizeCategory, err := uc.sizeCategoryRepo.GetSizeCategoryByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Size category not found", "method", "UpdateSizeCategory", "error", err, "id", id)
		return nil, fmt.Errorf("size category not found: %w", err)
	}

	if req.Name != "" && req.Name != existingSizeCategory.Name {
		exists, err := uc.sizeCategoryRepo.SizeCategoryExistsByName(req.Name, id)
		if err != nil {
			logger.Logger.Error("Failed to check if size category name exists", "method", "UpdateSizeCategory", "error", err, "id", id, "name", req.Name)
			return nil, fmt.Errorf("failed to check size category name: %w", err)
		}
		if exists {
			logger.Logger.Error("Size category name already exists", "method", "UpdateSizeCategory", "id", id, "name", req.Name)
			return nil, errors.New("size category with this name already exists")
		}
		existingSizeCategory.Name = req.Name
	}

	if err := uc.sizeCategoryRepo.UpdateSizeCategory(existingSizeCategory); err != nil {
		logger.Logger.Error("Failed to update size category", "method", "UpdateSizeCategory", "error", err, "sizeCategory", existingSizeCategory)
		return nil, fmt.Errorf("failed to update size category: %w", err)
	}

	return existingSizeCategory, nil
}

func (uc *SizeCategoryUseCase) DeleteSizeCategory(id uint) error {
	_, err := uc.sizeCategoryRepo.GetSizeCategoryByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Size category not found", "method", "DeleteSizeCategory", "error", err, "id", id)
		return fmt.Errorf("size category not found: %w", err)
	}

	if err := uc.sizeCategoryRepo.DeleteSizeCategory(id); err != nil {
		logger.Logger.Error("Failed to delete size category", "method", "DeleteSizeCategory", "error", err, "id", id)
		return fmt.Errorf("failed to delete size category: %w", err)
	}

	return nil
}

func (uc *SizeCategoryUseCase) UndoDeletedSizeCategory(id uint) error {
	showDeleted := true
	exists, err := uc.sizeCategoryRepo.SizeCategoryExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Size category not found", "method", "UndoDeletedSizeCategory", "error", err, "id", id)
		return fmt.Errorf("size category not found: %w", err)
	}

	if err := uc.sizeCategoryRepo.UndoDeletedSizeCategory(id); err != nil {
		logger.Logger.Error("Failed to undo deleted size category", "method", "UndoDeletedSizeCategory", "error", err, "id", id)
		return fmt.Errorf("failed to undo deleted size category: %w", err)
	}

	return nil
}
