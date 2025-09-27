package brand

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type BrandUseCase struct {
	brandRepo *BrandRepository
}

func NewBrandUseCase() *BrandUseCase {
	return &BrandUseCase{
		brandRepo: NewBrandRepository(),
	}
}

func (uc *BrandUseCase) GetAllBrands(includeStr string, showDeleted *bool, sortBy, sortOrder string) ([]Brand, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	brands, err := uc.brandRepo.GetAllBrands(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch brands", "method", "GetAllBrands", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch brands: %w", err)
	}

	return brands, nil
}

func (uc *BrandUseCase) GetBrandByID(id uint, includeStr string, showDeleted *bool) (*Brand, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	brand, err := uc.brandRepo.GetBrandByID(id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Brand not found", "method", "GetBrandByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("brand not found: %w", err)
	}

	return brand, nil
}

func (uc *BrandUseCase) CreateBrand(req *CreateBrandRequest) (*Brand, error) {
	req.Sanitize()

	exists, err := uc.brandRepo.BrandExistsByName(req.Name)
	if err != nil {
		logger.Logger.Error("Failed to check if brand exists", "method", "CreateBrand", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to check brand existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Brand already exists", "method", "CreateBrand", "name", req.Name)
		return nil, errors.New("brand with this name already exists")
	}

	brand := &Brand{
		Name:        req.Name,
		Description: req.Description,
	}

	if brand, err := uc.brandRepo.CreateBrand(brand); err != nil {
		logger.Logger.Error("Failed to create brand", "method", "CreateBrand", "error", err, "brand", brand)
		return nil, fmt.Errorf("failed to create brand: %w", err)
	}

	return brand, nil
}

func (uc *BrandUseCase) UpdateBrand(id uint, req *UpdateBrandRequest) (*Brand, error) {
	req.Sanitize()

	existingBrand, err := uc.brandRepo.GetBrandByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Brand not found", "method", "UpdateBrand", "error", err, "id", id)
		return nil, fmt.Errorf("brand not found: %w", err)
	}

	if req.Name != "" && req.Name != existingBrand.Name {
		exists, err := uc.brandRepo.BrandExistsByName(req.Name, id)
		if err != nil {
			logger.Logger.Error("Failed to check if brand name exists", "method", "UpdateBrand", "error", err, "id", id, "name", req.Name)
			return nil, fmt.Errorf("failed to check brand name: %w", err)
		}
		if exists {
			logger.Logger.Error("Brand name already exists", "method", "UpdateBrand", "id", id, "name", req.Name)
			return nil, errors.New("brand with this name already exists")
		}
		existingBrand.Name = req.Name
	}

	if req.Description != nil {
		existingBrand.Description = req.Description
	}

	if err := uc.brandRepo.UpdateBrand(existingBrand); err != nil {
		logger.Logger.Error("Failed to update brand", "method", "UpdateBrand", "error", err, "brand", existingBrand)
		return nil, fmt.Errorf("failed to update brand: %w", err)
	}

	return existingBrand, nil
}

func (uc *BrandUseCase) DeleteBrand(id uint) error {
	_, err := uc.brandRepo.GetBrandByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Brand not found", "method", "DeleteBrand", "error", err, "id", id)
		return fmt.Errorf("brand not found: %w", err)
	}

	if err := uc.brandRepo.DeleteBrand(id); err != nil {
		logger.Logger.Error("Failed to delete brand", "method", "DeleteBrand", "error", err, "id", id)
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	return nil
}

func (uc *BrandUseCase) UndoDeletedBrand(id uint) error {
	showDeleted := true
	exists, err := uc.brandRepo.BrandExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Brand not found", "method", "UndoDeletedBrand", "error", err, "id", id)
		return fmt.Errorf("brand not found: %w", err)
	}

	if err := uc.brandRepo.UndoDeletedBrand(id); err != nil {
		logger.Logger.Error("Failed to undo deleted brand", "method", "UndoDeletedBrand", "error", err, "id", id)
		return fmt.Errorf("failed to undo deleted brand: %w", err)
	}

	return nil
}
