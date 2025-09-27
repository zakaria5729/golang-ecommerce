package attribute_option

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/attribute_type"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AttributeOptionUseCase struct {
	attributeOptionRepo *AttributeOptionRepository
	attributeTypeRepo   *attribute_type.AttributeTypeRepository
}

func NewAttributeOptionUseCase() *AttributeOptionUseCase {
	return &AttributeOptionUseCase{
		attributeOptionRepo: NewAttributeOptionRepository(),
		attributeTypeRepo:   attribute_type.NewAttributeTypeRepository(),
	}
}

func (uc *AttributeOptionUseCase) GetAllAttributeOptions(includeStr string, showDeleted *bool, attributeTypeIDStr string, sortBy, sortOrder string) ([]AttributeOption, error) {
	include := utils.ParseCommaSeparatedString(includeStr)

	var attributeTypeID *uint
	if attributeTypeIDStr != "" {
		if id, err := strconv.ParseUint(attributeTypeIDStr, 10, 32); err == nil {
			uintID := uint(id)
			attributeTypeID = &uintID
		}
	}

	attributeOptions, err := uc.attributeOptionRepo.GetAllAttributeOptions(include, showDeleted, attributeTypeID, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch attribute options", "method", "GetAllAttributeOptions", "error", err, "include", include, "attributeTypeID", attributeTypeID, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch attribute options: %w", err)
	}

	return attributeOptions, nil
}

func (uc *AttributeOptionUseCase) GetAttributeOptionByID(id uint, includeStr string, showDeleted *bool) (*AttributeOption, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	attributeOption, err := uc.attributeOptionRepo.GetAttributeOptionByID(id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Attribute option not found", "method", "GetAttributeOptionByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("attribute option not found: %w", err)
	}

	return attributeOption, nil
}

func (uc *AttributeOptionUseCase) CreateAttributeOption(req *CreateAttributeOptionRequest) (*AttributeOption, error) {
	req.Sanitize()

	// Check if attribute type exists
	exists, err := uc.attributeTypeRepo.AttributeTypeExists(req.AttributeTypeID, nil)
	if err != nil {
		logger.Logger.Error("Failed to check if attribute type exists", "method", "CreateAttributeOption", "error", err, "attributeTypeID", req.AttributeTypeID)
		return nil, fmt.Errorf("failed to check attribute type existence: %w", err)
	}
	if !exists {
		logger.Logger.Error("Attribute type not found", "method", "CreateAttributeOption", "attributeTypeID", req.AttributeTypeID)
		return nil, errors.New("attribute type not found")
	}

	// Check if attribute option name already exists for this type
	exists, err = uc.attributeOptionRepo.AttributeOptionExistsByNameAndType(req.AttributeOptionName, req.AttributeTypeID)
	if err != nil {
		logger.Logger.Error("Failed to check if attribute option exists", "method", "CreateAttributeOption", "error", err, "attributeOptionName", req.AttributeOptionName, "attributeTypeID", req.AttributeTypeID)
		return nil, fmt.Errorf("failed to check attribute option existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Attribute option already exists", "method", "CreateAttributeOption", "attributeOptionName", req.AttributeOptionName, "attributeTypeID", req.AttributeTypeID)
		return nil, errors.New("attribute option with this name already exists for this attribute type")
	}

	attributeOption := &AttributeOption{
		AttributeTypeID:     &req.AttributeTypeID,
		AttributeOptionName: req.AttributeOptionName,
	}

	if attributeOption, err := uc.attributeOptionRepo.CreateAttributeOption(attributeOption); err != nil {
		logger.Logger.Error("Failed to create attribute option", "method", "CreateAttributeOption", "error", err, "attributeOption", attributeOption)
		return nil, fmt.Errorf("failed to create attribute option: %w", err)
	}

	return attributeOption, nil
}

func (uc *AttributeOptionUseCase) UpdateAttributeOption(id uint, req *UpdateAttributeOptionRequest) (*AttributeOption, error) {
	req.Sanitize()

	existingAttributeOption, err := uc.attributeOptionRepo.GetAttributeOptionByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Attribute option not found", "method", "UpdateAttributeOption", "error", err, "id", id)
		return nil, fmt.Errorf("attribute option not found: %w", err)
	}

	// Check if attribute type exists (if provided)
	if req.AttributeTypeID != nil && *req.AttributeTypeID != *existingAttributeOption.AttributeTypeID {
		exists, err := uc.attributeTypeRepo.AttributeTypeExists(*req.AttributeTypeID, nil)
		if err != nil {
			logger.Logger.Error("Failed to check if attribute type exists", "method", "UpdateAttributeOption", "error", err, "attributeTypeID", *req.AttributeTypeID)
			return nil, fmt.Errorf("failed to check attribute type existence: %w", err)
		}
		if !exists {
			logger.Logger.Error("Attribute type not found", "method", "UpdateAttributeOption", "attributeTypeID", *req.AttributeTypeID)
			return nil, errors.New("attribute type not found")
		}
		existingAttributeOption.AttributeTypeID = req.AttributeTypeID
	}

	// Check if attribute option name already exists for this type
	if req.AttributeOptionName != "" && req.AttributeOptionName != existingAttributeOption.AttributeOptionName {
		attributeTypeID := *existingAttributeOption.AttributeTypeID
		if req.AttributeTypeID != nil {
			attributeTypeID = *req.AttributeTypeID
		}

		exists, err := uc.attributeOptionRepo.AttributeOptionExistsByNameAndType(req.AttributeOptionName, attributeTypeID, id)
		if err != nil {
			logger.Logger.Error("Failed to check if attribute option name exists", "method", "UpdateAttributeOption", "error", err, "id", id, "attributeOptionName", req.AttributeOptionName, "attributeTypeID", attributeTypeID)
			return nil, fmt.Errorf("failed to check attribute option name: %w", err)
		}
		if exists {
			logger.Logger.Error("Attribute option name already exists", "method", "UpdateAttributeOption", "id", id, "attributeOptionName", req.AttributeOptionName, "attributeTypeID", attributeTypeID)
			return nil, errors.New("attribute option with this name already exists for this attribute type")
		}
		existingAttributeOption.AttributeOptionName = req.AttributeOptionName
	}

	if err := uc.attributeOptionRepo.UpdateAttributeOption(existingAttributeOption); err != nil {
		logger.Logger.Error("Failed to update attribute option", "method", "UpdateAttributeOption", "error", err, "attributeOption", existingAttributeOption)
		return nil, fmt.Errorf("failed to update attribute option: %w", err)
	}

	return existingAttributeOption, nil
}

func (uc *AttributeOptionUseCase) DeleteAttributeOption(id uint) error {
	_, err := uc.attributeOptionRepo.GetAttributeOptionByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Attribute option not found", "method", "DeleteAttributeOption", "error", err, "id", id)
		return fmt.Errorf("attribute option not found: %w", err)
	}

	if err := uc.attributeOptionRepo.DeleteAttributeOption(id); err != nil {
		logger.Logger.Error("Failed to delete attribute option", "method", "DeleteAttributeOption", "error", err, "id", id)
		return fmt.Errorf("failed to delete attribute option: %w", err)
	}

	return nil
}

func (uc *AttributeOptionUseCase) UndoDeletedAttributeOption(id uint) error {
	showDeleted := true
	exists, err := uc.attributeOptionRepo.AttributeOptionExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Attribute option not found", "method", "UndoDeletedAttributeOption", "error", err, "id", id)
		return fmt.Errorf("attribute option not found: %w", err)
	}

	if err := uc.attributeOptionRepo.UndoDeletedAttributeOption(id); err != nil {
		logger.Logger.Error("Failed to undo deleted attribute option", "method", "UndoDeletedAttributeOption", "error", err, "id", id)
		return fmt.Errorf("failed to undo deleted attribute option: %w", err)
	}

	return nil
}
