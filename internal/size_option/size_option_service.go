package size_option

import (
	"errors"
	"fmt"
	"strconv"

	m "github.com/easy-comerce/backend/internal/size_option/model"
)

type SizeOptionService struct {
	repo *SizeOptionRepository
}

func NewSizeOptionService(repo *SizeOptionRepository) *SizeOptionService {
	return &SizeOptionService{
		repo: repo,
	}
}

func (s *SizeOptionService) GetAllSizeOptions(showDeleted *bool, sizeCategoryIDStr string, sortBy, sortOrder string) ([]SizeOptionEntity, error) {
	var sizeCategoryID *uint
	if sizeCategoryIDStr != "" {
		if id, err := strconv.ParseUint(sizeCategoryIDStr, 10, 32); err == nil {
			categoryID := uint(id)
			sizeCategoryID = &categoryID
		}
	}

	sizeOptions, err := s.repo.GetAllSizeOptions(showDeleted, sizeCategoryID, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch size options: %w", err)
	}

	return sizeOptions, nil
}

func (s *SizeOptionService) GetSizeOptionByID(id uint, showDeleted *bool) (*SizeOptionEntity, error) {
	sizeOption, err := s.repo.GetSizeOptionByID(id, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("size option not found: %w", err)
	}

	return sizeOption, nil
}

func (s *SizeOptionService) CreateSizeOption(req *m.CreateSizeOptionRequest) (*SizeOptionEntity, error) {
	req.Sanitize()

	exists, err := s.repo.SizeOptionExistsByName(req.Name, req.SizeCategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to check size option existence: %w", err)
	}
	if exists {
		return nil, errors.New("size option with this name already exists in this category")
	}

	sizeOption := &SizeOptionEntity{
		Name:           req.Name,
		SortOrder:      req.SortOrder,
		SizeCategoryID: req.SizeCategoryID,
	}

	createdSizeOption, err := s.repo.CreateSizeOption(sizeOption)
	if err != nil {
		return nil, fmt.Errorf("failed to create size option: %w", err)
	}

	return createdSizeOption, nil
}

func (s *SizeOptionService) UpdateSizeOption(id uint, req *m.UpdateSizeOptionRequest) (*SizeOptionEntity, error) {
	req.Sanitize()

	existingSizeOption, err := s.repo.GetSizeOptionByID(id, nil)
	if err != nil {
		return nil, fmt.Errorf("size option not found: %w", err)
	}

	if req.Name != "" && req.Name != existingSizeOption.Name {
		exists, err := s.repo.SizeOptionExistsByName(req.Name, existingSizeOption.SizeCategoryID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check size option name: %w", err)
		}
		if exists {
			return nil, errors.New("size option with this name already exists in this category")
		}
		existingSizeOption.Name = req.Name
	}

	if req.SortOrder != nil {
		existingSizeOption.SortOrder = req.SortOrder
	}

	if err := s.repo.UpdateSizeOption(existingSizeOption); err != nil {
		return nil, fmt.Errorf("failed to update size option: %w", err)
	}

	return existingSizeOption, nil
}

func (s *SizeOptionService) DeleteSizeOption(id uint) error {
	_, err := s.repo.GetSizeOptionByID(id, nil)
	if err != nil {
		return fmt.Errorf("size option not found: %w", err)
	}

	if err := s.repo.DeleteSizeOption(id); err != nil {
		return fmt.Errorf("failed to delete size option: %w", err)
	}

	return nil
}

func (s *SizeOptionService) UndoDeletedSizeOption(id uint) error {
	showDeleted := true
	exists, err := s.repo.SizeOptionExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("size option not found: %w", err)
	}

	if err := s.repo.UndoDeletedSizeOption(id); err != nil {
		return fmt.Errorf("failed to undo deleted size option: %w", err)
	}

	return nil
}
