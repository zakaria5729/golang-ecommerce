package size_option

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SizeOptionUseCase struct {
	sizeOptionRepo *SizeOptionRepository
}

func NewSizeOptionUseCase() *SizeOptionUseCase {
	return &SizeOptionUseCase{
		sizeOptionRepo: NewSizeOptionRepository(),
	}
}

func (uc *SizeOptionUseCase) GetAllSizeOptions(includeStr string, showDeleted *bool, sizeCategoryIDStr string, sortBy, sortOrder string) ([]SizeOption, error) {
	include := utils.ParseCommaSeparatedString(includeStr)

	var sizeCategoryID *uint
	if sizeCategoryIDStr != "" {
		if id, err := strconv.ParseUint(sizeCategoryIDStr, 10, 32); err == nil {
			categoryID := uint(id)
			sizeCategoryID = &categoryID
		}
	}

	sizeOptions, err := uc.sizeOptionRepo.GetAllSizeOptions(include, showDeleted, sizeCategoryID, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch size options", "method", "GetAllSizeOptions", "error", err, "include", include, "sizeCategoryID", sizeCategoryID, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch size options: %w", err)
	}

	return sizeOptions, nil
}

func (uc *SizeOptionUseCase) GetSizeOptionByID(id uint, includeStr string, showDeleted *bool) (*SizeOption, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	sizeOption, err := uc.sizeOptionRepo.GetSizeOptionByID(id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Size option not found", "method", "GetSizeOptionByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("size option not found: %w", err)
	}

	return sizeOption, nil
}

func (uc *SizeOptionUseCase) CreateSizeOption(req *CreateSizeOptionRequest) (*SizeOption, error) {
	req.Sanitize()

	exists, err := uc.sizeOptionRepo.SizeOptionExistsByName(req.Name, req.SizeCategoryID)
	if err != nil {
		logger.Logger.Error("Failed to check if size option exists", "method", "CreateSizeOption", "error", err, "name", req.Name, "sizeCategoryID", req.SizeCategoryID)
		return nil, fmt.Errorf("failed to check size option existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Size option already exists", "method", "CreateSizeOption", "name", req.Name, "sizeCategoryID", req.SizeCategoryID)
		return nil, errors.New("size option with this name already exists in this category")
	}

	sizeOption := &SizeOption{
		Name:           req.Name,
		SortOrder:      req.SortOrder,
		SizeCategoryID: req.SizeCategoryID,
	}

	if sizeOption, err := uc.sizeOptionRepo.CreateSizeOption(sizeOption); err != nil {
		logger.Logger.Error("Failed to create size option", "method", "CreateSizeOption", "error", err, "sizeOption", sizeOption)
		return nil, fmt.Errorf("failed to create size option: %w", err)
	}

	return sizeOption, nil
}

func (uc *SizeOptionUseCase) UpdateSizeOption(id uint, req *UpdateSizeOptionRequest) (*SizeOption, error) {
	req.Sanitize()

	existingSizeOption, err := uc.sizeOptionRepo.GetSizeOptionByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Size option not found", "method", "UpdateSizeOption", "error", err, "id", id)
		return nil, fmt.Errorf("size option not found: %w", err)
	}

	if req.Name != "" && req.Name != existingSizeOption.Name {
		exists, err := uc.sizeOptionRepo.SizeOptionExistsByName(req.Name, existingSizeOption.SizeCategoryID, id)
		if err != nil {
			logger.Logger.Error("Failed to check if size option name exists", "method", "UpdateSizeOption", "error", err, "id", id, "name", req.Name, "sizeCategoryID", existingSizeOption.SizeCategoryID)
			return nil, fmt.Errorf("failed to check size option name: %w", err)
		}
		if exists {
			logger.Logger.Error("Size option name already exists", "method", "UpdateSizeOption", "id", id, "name", req.Name, "sizeCategoryID", existingSizeOption.SizeCategoryID)
			return nil, errors.New("size option with this name already exists in this category")
		}
		existingSizeOption.Name = req.Name
	}

	if req.SortOrder != nil {
		existingSizeOption.SortOrder = req.SortOrder
	}

	if err := uc.sizeOptionRepo.UpdateSizeOption(existingSizeOption); err != nil {
		logger.Logger.Error("Failed to update size option", "method", "UpdateSizeOption", "error", err, "sizeOption", existingSizeOption)
		return nil, fmt.Errorf("failed to update size option: %w", err)
	}

	return existingSizeOption, nil
}

func (uc *SizeOptionUseCase) DeleteSizeOption(id uint) error {
	_, err := uc.sizeOptionRepo.GetSizeOptionByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Size option not found", "method", "DeleteSizeOption", "error", err, "id", id)
		return fmt.Errorf("size option not found: %w", err)
	}

	if err := uc.sizeOptionRepo.DeleteSizeOption(id); err != nil {
		logger.Logger.Error("Failed to delete size option", "method", "DeleteSizeOption", "error", err, "id", id)
		return fmt.Errorf("failed to delete size option: %w", err)
	}

	return nil
}

func (uc *SizeOptionUseCase) UndoDeletedSizeOption(id uint) error {
	showDeleted := true
	exists, err := uc.sizeOptionRepo.SizeOptionExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Size option not found", "method", "UndoDeletedSizeOption", "error", err, "id", id)
		return fmt.Errorf("size option not found: %w", err)
	}

	if err := uc.sizeOptionRepo.UndoDeletedSizeOption(id); err != nil {
		logger.Logger.Error("Failed to undo deleted size option", "method", "UndoDeletedSizeOption", "error", err, "id", id)
		return fmt.Errorf("failed to undo deleted size option: %w", err)
	}

	return nil
}
