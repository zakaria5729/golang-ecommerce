package category

import (
	"fmt"
	"strings"
)

type CategoryUseCase struct {
	repo *CategoryRepository
}

func NewCategoryUseCase() *CategoryUseCase {
	return &CategoryUseCase{
		repo: NewCategoryRepository(),
	}
}

func (uc *CategoryUseCase) GetAllCategories(includeStr string, parentIDStr string, isActiveStr string) ([]Category, error) {
	// Parse include parameter
	var include []string
	if includeStr != "" {
		include = strings.Split(includeStr, ",")
		for i, field := range include {
			include[i] = strings.TrimSpace(field)
		}
	}

	// Parse parent_id parameter
	var parentID *uint
	if parentIDStr != "" && parentIDStr != "null" {
		if id, err := parseUint(parentIDStr); err == nil {
			parentID = &id
		}
	}

	// Parse is_active parameter
	var isActive *bool
	if isActiveStr != "" {
		active := isActiveStr == "true"
		isActive = &active
	}

	return uc.repo.GetAllCategories(include, parentID, isActive)
}

func (uc *CategoryUseCase) GetCategoryByID(id uint, includeStr string) (*Category, error) {
	// Parse include parameter
	var include []string
	if includeStr != "" {
		include = strings.Split(includeStr, ",")
		for i, field := range include {
			include[i] = strings.TrimSpace(field)
		}
	}

	return uc.repo.GetCategoryByID(id, include)
}

func (uc *CategoryUseCase) CreateCategory(category *Category) error {
	return uc.repo.CreateCategory(category)
}

func (uc *CategoryUseCase) UpdateCategory(category *Category) error {
	return uc.repo.UpdateCategory(category)
}

func (uc *CategoryUseCase) DeleteCategory(id uint) error {
	return uc.repo.DeleteCategory(id)
}

// Helper function to parse uint
func parseUint(s string) (uint, error) {
	var result uint
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
