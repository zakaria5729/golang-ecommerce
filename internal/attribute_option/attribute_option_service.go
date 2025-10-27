package attribute_option

import (
	"errors"
	"fmt"
	"strconv"

	m "github.com/easy-comerce/backend/internal/attribute_option/model"
	"github.com/easy-comerce/backend/internal/attribute_type"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AttributeOptionService interface {
	GetAllAttributeOptions(includeStr string, showDeleted *bool, attributeTypeIDStr string, sortBy, sortOrder string) ([]AttributeOptionEntity, error)
	GetAttributeOptionByID(id uint, includeStr string, showDeleted *bool) (*AttributeOptionEntity, error)
	CreateAttributeOption(req *m.CreateAttributeOptionRequest) (*AttributeOptionEntity, error)
	UpdateAttributeOption(id uint, req *m.UpdateAttributeOptionRequest) (*AttributeOptionEntity, error)
	DeleteAttributeOption(id uint) error
	UndoDeletedAttributeOption(id uint) error
}

type attributeOptionService struct {
	optionRepo AttributeOptionRepository
	typeRepo   attribute_type.AttributeTypeRepository
}

func NewAttributeOptionService(
	optionRepo AttributeOptionRepository,
	typeRepo attribute_type.AttributeTypeRepository,
) AttributeOptionService {
	return &attributeOptionService{
		optionRepo: optionRepo,
		typeRepo:   typeRepo,
	}
}

func (s *attributeOptionService) GetAllAttributeOptions(includeStr string, showDeleted *bool, attributeTypeIDStr string, sortBy, sortOrder string) ([]AttributeOptionEntity, error) {
	include := utils.ParseCommaSeparatedString(includeStr)

	var attributeTypeID *uint
	if attributeTypeIDStr != "" {
		if id, err := strconv.ParseUint(attributeTypeIDStr, 10, 32); err == nil {
			uintID := uint(id)
			attributeTypeID = &uintID
		}
	}

	attributeOptions, err := s.optionRepo.GetAllAttributeOptions(include, showDeleted, attributeTypeID, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attribute options: %w", err)
	}

	return attributeOptions, nil
}

func (s *attributeOptionService) GetAttributeOptionByID(id uint, includeStr string, showDeleted *bool) (*AttributeOptionEntity, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	attributeOption, err := s.optionRepo.GetAttributeOptionByID(id, include, showDeleted)

	if err != nil {
		return nil, fmt.Errorf("attribute option not found: %w", err)
	}

	return attributeOption, nil
}

func (s *attributeOptionService) CreateAttributeOption(req *m.CreateAttributeOptionRequest) (*AttributeOptionEntity, error) {
	req.Sanitize()

	exists, err := s.typeRepo.AttributeTypeExists(req.AttributeTypeID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check attribute type existence: %w", err)
	}
	if !exists {
		return nil, errors.New("attribute type not found")
	}

	exists, err = s.optionRepo.AttributeOptionExistsByNameAndType(req.AttributeOptionName, req.AttributeTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check attribute option existence: %w", err)
	}
	if exists {
		return nil, errors.New("attribute option with this name already exists for this attribute type")
	}

	attributeOption := &AttributeOptionEntity{
		AttributeTypeID:     &req.AttributeTypeID,
		AttributeOptionName: req.AttributeOptionName,
	}

	attributeOption, err = s.optionRepo.CreateAttributeOption(attributeOption)
	if err != nil {
		return nil, fmt.Errorf("failed to create attribute option: %w", err)
	}

	return attributeOption, nil
}

func (s *attributeOptionService) UpdateAttributeOption(id uint, req *m.UpdateAttributeOptionRequest) (*AttributeOptionEntity, error) {
	req.Sanitize()

	existingAttributeOption, err := s.optionRepo.GetAttributeOptionByID(id, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("attribute option not found: %w", err)
	}

	if req.AttributeTypeID != nil && *req.AttributeTypeID != *existingAttributeOption.AttributeTypeID {
		exists, err := s.typeRepo.AttributeTypeExists(*req.AttributeTypeID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check attribute type existence: %w", err)
		}
		if !exists {
			return nil, errors.New("attribute type not found")
		}
		existingAttributeOption.AttributeTypeID = req.AttributeTypeID
	}

	if req.AttributeOptionName != "" && req.AttributeOptionName != existingAttributeOption.AttributeOptionName {
		attributeTypeID := *existingAttributeOption.AttributeTypeID
		if req.AttributeTypeID != nil {
			attributeTypeID = *req.AttributeTypeID
		}

		exists, err := s.optionRepo.AttributeOptionExistsByNameAndType(req.AttributeOptionName, attributeTypeID)
		if err != nil {
			return nil, fmt.Errorf("failed to check attribute option name: %w", err)
		}
		if exists {
			return nil, errors.New("attribute option with this name already exists for this attribute type")
		}
		existingAttributeOption.AttributeOptionName = req.AttributeOptionName
	}

	if err := s.optionRepo.UpdateAttributeOption(existingAttributeOption); err != nil {
		return nil, fmt.Errorf("failed to update attribute option: %w", err)
	}

	return existingAttributeOption, nil
}

func (s *attributeOptionService) DeleteAttributeOption(id uint) error {
	_, err := s.optionRepo.GetAttributeOptionByID(id, nil, nil)
	if err != nil {
		return fmt.Errorf("attribute option not found: %w", err)
	}

	if err := s.optionRepo.DeleteAttributeOption(id); err != nil {
		return fmt.Errorf("failed to delete attribute option: %w", err)
	}

	return nil
}

func (s *attributeOptionService) UndoDeletedAttributeOption(id uint) error {
	showDeleted := true
	exists, err := s.optionRepo.AttributeOptionExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("attribute option not found: %w", err)
	}

	if err := s.optionRepo.UndoDeletedAttributeOption(id); err != nil {
		return fmt.Errorf("failed to undo deleted attribute option: %w", err)
	}

	return nil
}
