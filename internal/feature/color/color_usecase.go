package color

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ColorUseCase struct {
	colorRepo *ColorRepository
}

func NewColorUseCase() *ColorUseCase {
	return &ColorUseCase{
		colorRepo: NewColorRepository(),
	}
}

func (uc *ColorUseCase) GetAllColors(includeStr string, showDeleted *bool, sortBy, sortOrder string) ([]Color, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	colors, err := uc.colorRepo.GetAllColors(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch colors", "method", "GetAllColors", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch colors: %w", err)
	}

	return colors, nil
}

func (uc *ColorUseCase) GetColorByID(id uint, includeStr string, showDeleted *bool) (*Color, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	color, err := uc.colorRepo.GetColorByID(id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Color not found", "method", "GetColorByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("color not found: %w", err)
	}

	return color, nil
}

func (uc *ColorUseCase) CreateColor(req *CreateColorRequest) (*Color, error) {
	req.Sanitize()

	exists, err := uc.colorRepo.ColorExistsByName(req.Name)
	if err != nil {
		logger.Logger.Error("Failed to check if color exists", "method", "CreateColor", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to check color existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Color already exists", "method", "CreateColor", "name", req.Name)
		return nil, errors.New("color with this name already exists")
	}

	color := &Color{
		Name: req.Name,
	}

	if color, err := uc.colorRepo.CreateColor(color); err != nil {
		logger.Logger.Error("Failed to create color", "method", "CreateColor", "error", err, "color", color)
		return nil, fmt.Errorf("failed to create color: %w", err)
	}

	return color, nil
}

func (uc *ColorUseCase) UpdateColor(id uint, req *UpdateColorRequest) (*Color, error) {
	req.Sanitize()

	existingColor, err := uc.colorRepo.GetColorByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Color not found", "method", "UpdateColor", "error", err, "id", id)
		return nil, fmt.Errorf("color not found: %w", err)
	}

	if req.Name != "" && req.Name != existingColor.Name {
		exists, err := uc.colorRepo.ColorExistsByName(req.Name, id)
		if err != nil {
			logger.Logger.Error("Failed to check if color name exists", "method", "UpdateColor", "error", err, "id", id, "name", req.Name)
			return nil, fmt.Errorf("failed to check color name: %w", err)
		}
		if exists {
			logger.Logger.Error("Color name already exists", "method", "UpdateColor", "id", id, "name", req.Name)
			return nil, errors.New("color with this name already exists")
		}
		existingColor.Name = req.Name
	}

	if err := uc.colorRepo.UpdateColor(existingColor); err != nil {
		logger.Logger.Error("Failed to update color", "method", "UpdateColor", "error", err, "color", existingColor)
		return nil, fmt.Errorf("failed to update color: %w", err)
	}

	return existingColor, nil
}

func (uc *ColorUseCase) DeleteColor(id uint) error {
	_, err := uc.colorRepo.GetColorByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Color not found", "method", "DeleteColor", "error", err, "id", id)
		return fmt.Errorf("color not found: %w", err)
	}

	if err := uc.colorRepo.DeleteColor(id); err != nil {
		logger.Logger.Error("Failed to delete color", "method", "DeleteColor", "error", err, "id", id)
		return fmt.Errorf("failed to delete color: %w", err)
	}

	return nil
}

func (uc *ColorUseCase) UndoDeletedColor(id uint) error {
	showDeleted := true
	exists, err := uc.colorRepo.ColorExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Color not found", "method", "UndoDeletedColor", "error", err, "id", id)
		return fmt.Errorf("color not found: %w", err)
	}

	if err := uc.colorRepo.UndoDeletedColor(id); err != nil {
		logger.Logger.Error("Failed to undo deleted color", "method", "UndoDeletedColor", "error", err, "id", id)
		return fmt.Errorf("failed to undo deleted color: %w", err)
	}

	return nil
}
