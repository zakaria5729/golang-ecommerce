package attribute_type

import (
	"errors"
	"fmt"

	m "github.com/easy-comerce/backend/internal/attribute_type/model"
)

type AttributeTypeService interface {
	GetAllAttributeTypes(showDeleted *bool, sortBy, sortOrder string) ([]AttributeTypeEntity, error)
	GetAttributeTypeByID(id uint, showDeleted *bool) (*AttributeTypeEntity, error)
	CreateAttributeType(req *m.CreateAttributeTypeRequest) (*AttributeTypeEntity, error)
	UpdateAttributeType(id uint, req *m.UpdateAttributeTypeRequest) (*AttributeTypeEntity, error)
	DeleteAttributeType(id uint) error
	UndoDeletedAttributeType(id uint) error
}

type attributeTypeService struct {
	repo AttributeTypeRepository
}

func NewAttributeTypeService(repo AttributeTypeRepository) AttributeTypeService {
	return &attributeTypeService{
		repo: repo,
	}
}

func (s *attributeTypeService) GetAllAttributeTypes(showDeleted *bool, sortBy, sortOrder string) ([]AttributeTypeEntity, error) {
	attributeTypes, err := s.repo.GetAllAttributeTypes(showDeleted, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attribute types: %w", err)
	}

	return attributeTypes, nil
}

func (s *attributeTypeService) GetAttributeTypeByID(id uint, showDeleted *bool) (*AttributeTypeEntity, error) {
	attributeType, err := s.repo.GetAttributeTypeByID(id, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("attribute type not found: %w", err)
	}

	return attributeType, nil
}

func (s *attributeTypeService) CreateAttributeType(req *m.CreateAttributeTypeRequest) (*AttributeTypeEntity, error) {
	req.Sanitize()

	exists, err := s.repo.AttributeTypeExistsByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check attribute type existence: %w", err)
	}
	if exists {
		return nil, errors.New("attribute type with this name already exists")
	}

	attributeType := &AttributeTypeEntity{
		Name: req.Name,
	}

	createdAttributeType, err := s.repo.CreateAttributeType(attributeType)
	if err != nil {
		return nil, fmt.Errorf("failed to create attribute type: %w", err)
	}

	return createdAttributeType, nil
}

func (s *attributeTypeService) UpdateAttributeType(id uint, req *m.UpdateAttributeTypeRequest) (*AttributeTypeEntity, error) {
	req.Sanitize()

	existingAttributeType, err := s.repo.GetAttributeTypeByID(id, nil)
	if err != nil {
		return nil, fmt.Errorf("attribute type not found: %w", err)
	}

	if req.Name != "" && req.Name != existingAttributeType.Name {
		exists, err := s.repo.AttributeTypeExistsByName(req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check attribute type name: %w", err)
		}
		if exists {
			return nil, errors.New("attribute type with this name already exists")
		}
		existingAttributeType.Name = req.Name
	}

	if err := s.repo.UpdateAttributeType(existingAttributeType); err != nil {
		return nil, fmt.Errorf("failed to update attribute type: %w", err)
	}

	return existingAttributeType, nil
}

func (s *attributeTypeService) DeleteAttributeType(id uint) error {
	_, err := s.repo.GetAttributeTypeByID(id, nil)
	if err != nil {
		return fmt.Errorf("attribute type not found: %w", err)
	}

	if err := s.repo.DeleteAttributeType(id); err != nil {
		return fmt.Errorf("failed to delete attribute type: %w", err)
	}

	return nil
}

func (s *attributeTypeService) UndoDeletedAttributeType(id uint) error {
	showDeleted := true
	exists, err := s.repo.AttributeTypeExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("attribute type not found: %w", err)
	}

	if err := s.repo.UndoDeletedAttributeType(id); err != nil {
		return fmt.Errorf("failed to undo deleted attribute type: %w", err)
	}

	return nil
}
