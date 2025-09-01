package category

import (
	"github.com/easy-comerce/backend/pkg/utils"
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
	parentID, _ := utils.ParseUint(parentIDStr)
	isActive := utils.ParseBoolPtr(isActiveStr)
	return uc.repo.GetAllCategories(utils.ParseCommaSeparatedString(includeStr), parentID, isActive)
}

func (uc *CategoryUseCase) GetCategoryByID(id uint, includeStr string) (*Category, error) {
	return uc.repo.GetCategoryByID(id, utils.ParseCommaSeparatedString(includeStr))
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
