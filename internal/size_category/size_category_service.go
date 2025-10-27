package size_category

import (
	"errors"
	"fmt"

	m "github.com/easy-comerce/backend/internal/size_category/model"
)

type SizeCategoryService interface {
	GetAllSizeCategories(showDeleted *bool, sortBy, sortOrder string) ([]SizeCategoryEntity, error)
	GetSizeCategoryByID(id uint, showDeleted *bool) (*SizeCategoryEntity, error)
	CreateSizeCategory(req *m.CreateSizeCategoryRequest) (*SizeCategoryEntity, error)
	UpdateSizeCategory(id uint, req *m.UpdateSizeCategoryRequest) (*SizeCategoryEntity, error)
	DeleteSizeCategory(id uint) error
	UndoDeletedSizeCategory(id uint) error
}

type sizeCategoryService struct {
	repo SizeCategoryRepository
}

func NewSizeCategoryService(repo SizeCategoryRepository) SizeCategoryService {
	return &sizeCategoryService{
		repo: repo,
	}
}

func (s *sizeCategoryService) GetAllSizeCategories(showDeleted *bool, sortBy, sortOrder string) ([]SizeCategoryEntity, error) {
	sizeCategories, err := s.repo.GetAllSizeCategories(showDeleted, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch size categories: %w", err)
	}

	return sizeCategories, nil
}

func (s *sizeCategoryService) GetSizeCategoryByID(id uint, showDeleted *bool) (*SizeCategoryEntity, error) {
	sizeCategory, err := s.repo.GetSizeCategoryByID(id, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("size category not found: %w", err)
	}

	return sizeCategory, nil
}

func (s *sizeCategoryService) CreateSizeCategory(req *m.CreateSizeCategoryRequest) (*SizeCategoryEntity, error) {
	req.Sanitize()

	exists, err := s.repo.SizeCategoryExistsByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check size category existence: %w", err)
	}
	if exists {
		return nil, errors.New("size category with this name already exists")
	}

	sizeCategory := &SizeCategoryEntity{
		Name: req.Name,
	}

	sizeCategory, err = s.repo.CreateSizeCategory(sizeCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to create size category: %w", err)
	}

	return sizeCategory, nil
}

func (s *sizeCategoryService) UpdateSizeCategory(id uint, req *m.UpdateSizeCategoryRequest) (*SizeCategoryEntity, error) {
	req.Sanitize()

	existingSizeCategory, err := s.repo.GetSizeCategoryByID(id, nil)
	if err != nil {
		return nil, fmt.Errorf("size category not found: %w", err)
	}

	if req.Name != "" && req.Name != existingSizeCategory.Name {
		exists, err := s.repo.SizeCategoryExistsByName(req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check size category name: %w", err)
		}
		if exists {
			return nil, errors.New("size category with this name already exists")
		}
		existingSizeCategory.Name = req.Name
	}

	if err := s.repo.UpdateSizeCategory(existingSizeCategory); err != nil {
		return nil, fmt.Errorf("failed to update size category: %w", err)
	}

	return existingSizeCategory, nil
}

func (s *sizeCategoryService) DeleteSizeCategory(id uint) error {
	_, err := s.repo.GetSizeCategoryByID(id, nil)
	if err != nil {
		return fmt.Errorf("size category not found: %w", err)
	}

	if err := s.repo.DeleteSizeCategory(id); err != nil {
		return fmt.Errorf("failed to delete size category: %w", err)
	}

	return nil
}

func (s *sizeCategoryService) UndoDeletedSizeCategory(id uint) error {
	showDeleted := true
	exists, err := s.repo.SizeCategoryExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("size category not found: %w", err)
	}

	if err := s.repo.UndoDeletedSizeCategory(id); err != nil {
		return fmt.Errorf("failed to undo deleted size category: %w", err)
	}

	return nil
}
