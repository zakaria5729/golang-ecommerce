package color

import (
	"errors"
	"fmt"

	m "github.com/easy-comerce/backend/internal/color/model"
)

type ColorService interface {
	GetAllColors(showDeleted *bool, sortBy, sortOrder string) ([]ColorEntity, error)
	GetColorByID(id uint, showDeleted *bool) (*ColorEntity, error)
	CreateColor(req *m.CreateColorRequest) (*ColorEntity, error)
	UpdateColor(id uint, req *m.UpdateColorRequest) (*ColorEntity, error)
	DeleteColor(id uint) error
	UndoDeletedColor(id uint) error
}

type colorService struct {
	colorRepo ColorRepository
}

func NewColorService(repo ColorRepository) ColorService {
	return &colorService{
		colorRepo: repo,
	}
}

func (s *colorService) GetAllColors(showDeleted *bool, sortBy, sortOrder string) ([]ColorEntity, error) {
	colors, err := s.colorRepo.GetAllColors(showDeleted, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch colors: %w", err)
	}

	return colors, nil
}

func (s *colorService) GetColorByID(id uint, showDeleted *bool) (*ColorEntity, error) {
	color, err := s.colorRepo.GetColorByID(id, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("color not found: %w", err)
	}

	return color, nil
}

func (s *colorService) CreateColor(req *m.CreateColorRequest) (*ColorEntity, error) {
	req.Sanitize()

	exists, err := s.colorRepo.ColorExistsByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check color existence: %w", err)
	}
	if exists {
		return nil, errors.New("color with this name already exists")
	}

	color := &ColorEntity{
		Name: req.Name,
	}

	color, err = s.colorRepo.CreateColor(color)
	if err != nil {
		return nil, fmt.Errorf("failed to create color: %w", err)
	}

	return color, nil
}

func (s *colorService) UpdateColor(id uint, req *m.UpdateColorRequest) (*ColorEntity, error) {
	req.Sanitize()

	existingColor, err := s.colorRepo.GetColorByID(id, nil)
	if err != nil {
		return nil, fmt.Errorf("color not found: %w", err)
	}

	if req.Name != "" && req.Name != existingColor.Name {
		exists, err := s.colorRepo.ColorExistsByName(req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check color name: %w", err)
		}
		if exists {
			return nil, errors.New("color with this name already exists")
		}
		existingColor.Name = req.Name
	}

	if err := s.colorRepo.UpdateColor(existingColor); err != nil {
		return nil, fmt.Errorf("failed to update color: %w", err)
	}

	return existingColor, nil
}

func (s *colorService) DeleteColor(id uint) error {
	_, err := s.colorRepo.GetColorByID(id, nil)
	if err != nil {
		return fmt.Errorf("color not found: %w", err)
	}

	if err := s.colorRepo.DeleteColor(id); err != nil {
		return fmt.Errorf("failed to delete color: %w", err)
	}

	return nil
}

func (s *colorService) UndoDeletedColor(id uint) error {
	showDeleted := true
	exists, err := s.colorRepo.ColorExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("color not found: %w", err)
	}

	if err := s.colorRepo.UndoDeletedColor(id); err != nil {
		return fmt.Errorf("failed to undo deleted color: %w", err)
	}

	return nil
}
