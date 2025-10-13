package brand

import (
	"errors"
	"fmt"
)

type BrandService struct {
	repo *BrandRepository
}

func NewBrandService(repo *BrandRepository) *BrandService {
	return &BrandService{
		repo: repo,
	}
}

func (s *BrandService) GetAllBrands(showDeleted *bool, sortBy, sortOrder string) ([]Brand, error) {
	brands, err := s.repo.GetAllBrands(showDeleted, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch brands: %w", err)
	}

	return brands, nil
}

func (s *BrandService) GetBrandByID(id uint, showDeleted *bool) (*Brand, error) {
	brand, err := s.repo.GetBrandByID(id, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}

	return brand, nil
}

func (s *BrandService) CreateBrand(req *CreateBrandRequest) (*Brand, error) {
	req.Sanitize()

	exists, err := s.repo.BrandExistsByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check brand existence: %w", err)
	}
	if exists {
		return nil, errors.New("brand with this name already exists")
	}

	brand := &Brand{
		Name:        req.Name,
		Description: req.Description,
	}

	brand, err = s.repo.CreateBrand(brand)
	if err != nil {
		return nil, fmt.Errorf("failed to create brand: %w", err)
	}

	return brand, nil
}

func (s *BrandService) UpdateBrand(id uint, req *UpdateBrandRequest) (*Brand, error) {
	req.Sanitize()

	existingBrand, err := s.repo.GetBrandByID(id, nil)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}

	if req.Name != "" && req.Name != existingBrand.Name {
		exists, err := s.repo.BrandExistsByName(req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check brand name: %w", err)
		}
		if exists {
			return nil, errors.New("brand with this name already exists")
		}
		existingBrand.Name = req.Name
	}

	if req.Description != nil {
		existingBrand.Description = req.Description
	}

	if err := s.repo.UpdateBrand(existingBrand); err != nil {
		return nil, fmt.Errorf("failed to update brand: %w", err)
	}

	return existingBrand, nil
}

func (s *BrandService) DeleteBrand(id uint) error {
	_, err := s.repo.GetBrandByID(id, nil)
	if err != nil {
		return fmt.Errorf("brand not found: %w", err)
	}

	if err := s.repo.DeleteBrand(id); err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	return nil
}

func (s *BrandService) UndoDeletedBrand(id uint) error {
	showDeleted := true
	exists, err := s.repo.BrandExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("brand not found: %w", err)
	}

	if err := s.repo.UndoDeletedBrand(id); err != nil {
		return fmt.Errorf("failed to undo deleted brand: %w", err)
	}

	return nil
}
