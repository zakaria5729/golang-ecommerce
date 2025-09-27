package attribute_type

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AttributeTypeUseCase struct {
	attributeTypeRepo *AttributeTypeRepository
}

func NewAttributeTypeUseCase() *AttributeTypeUseCase {
	return &AttributeTypeUseCase{
		attributeTypeRepo: NewAttributeTypeRepository(),
	}
}

func (uc *AttributeTypeUseCase) GetAllAttributeTypes(includeStr string, showDeleted *bool, sortBy, sortOrder string) ([]AttributeType, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	attributeTypes, err := uc.attributeTypeRepo.GetAllAttributeTypes(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch attribute types", "method", "GetAllAttributeTypes", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch attribute types: %w", err)
	}

	return attributeTypes, nil
}

func (uc *AttributeTypeUseCase) GetAttributeTypeByID(id uint, includeStr string, showDeleted *bool) (*AttributeType, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	attributeType, err := uc.attributeTypeRepo.GetAttributeTypeByID(id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Attribute type not found", "method", "GetAttributeTypeByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("attribute type not found: %w", err)
	}

	return attributeType, nil
}

func (uc *AttributeTypeUseCase) CreateAttributeType(req *CreateAttributeTypeRequest) (*AttributeType, error) {
	req.Sanitize()

	exists, err := uc.attributeTypeRepo.AttributeTypeExistsByName(req.Name)
	if err != nil {
		logger.Logger.Error("Failed to check if attribute type exists", "method", "CreateAttributeType", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to check attribute type existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Attribute type already exists", "method", "CreateAttributeType", "name", req.Name)
		return nil, errors.New("attribute type with this name already exists")
	}

	attributeType := &AttributeType{
		Name: req.Name,
	}

	if attributeType, err := uc.attributeTypeRepo.CreateAttributeType(attributeType); err != nil {
		logger.Logger.Error("Failed to create attribute type", "method", "CreateAttributeType", "error", err, "attributeType", attributeType)
		return nil, fmt.Errorf("failed to create attribute type: %w", err)
	}

	return attributeType, nil
}

func (uc *AttributeTypeUseCase) UpdateAttributeType(id uint, req *UpdateAttributeTypeRequest) (*AttributeType, error) {
	req.Sanitize()

	existingAttributeType, err := uc.attributeTypeRepo.GetAttributeTypeByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Attribute type not found", "method", "UpdateAttributeType", "error", err, "id", id)
		return nil, fmt.Errorf("attribute type not found: %w", err)
	}

	if req.Name != "" && req.Name != existingAttributeType.Name {
		exists, err := uc.attributeTypeRepo.AttributeTypeExistsByName(req.Name, id)
		if err != nil {
			logger.Logger.Error("Failed to check if attribute type name exists", "method", "UpdateAttributeType", "error", err, "id", id, "name", req.Name)
			return nil, fmt.Errorf("failed to check attribute type name: %w", err)
		}
		if exists {
			logger.Logger.Error("Attribute type name already exists", "method", "UpdateAttributeType", "id", id, "name", req.Name)
			return nil, errors.New("attribute type with this name already exists")
		}
		existingAttributeType.Name = req.Name
	}

	if err := uc.attributeTypeRepo.UpdateAttributeType(existingAttributeType); err != nil {
		logger.Logger.Error("Failed to update attribute type", "method", "UpdateAttributeType", "error", err, "attributeType", existingAttributeType)
		return nil, fmt.Errorf("failed to update attribute type: %w", err)
	}

	return existingAttributeType, nil
}

func (uc *AttributeTypeUseCase) DeleteAttributeType(id uint) error {
	_, err := uc.attributeTypeRepo.GetAttributeTypeByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Attribute type not found", "method", "DeleteAttributeType", "error", err, "id", id)
		return fmt.Errorf("attribute type not found: %w", err)
	}

	if err := uc.attributeTypeRepo.DeleteAttributeType(id); err != nil {
		logger.Logger.Error("Failed to delete attribute type", "method", "DeleteAttributeType", "error", err, "id", id)
		return fmt.Errorf("failed to delete attribute type: %w", err)
	}

	return nil
}

func (uc *AttributeTypeUseCase) UndoDeletedAttributeType(id uint) error {
	showDeleted := true
	exists, err := uc.attributeTypeRepo.AttributeTypeExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Attribute type not found", "method", "UndoDeletedAttributeType", "error", err, "id", id)
		return fmt.Errorf("attribute type not found: %w", err)
	}

	if err := uc.attributeTypeRepo.UndoDeletedAttributeType(id); err != nil {
		logger.Logger.Error("Failed to undo deleted attribute type", "method", "UndoDeletedAttributeType", "error", err, "id", id)
		return fmt.Errorf("failed to undo deleted attribute type: %w", err)
	}

	return nil
}
