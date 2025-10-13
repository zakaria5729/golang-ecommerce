package size_category

import (
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateSizeCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateSizeCategoryRequest struct {
	Name string `json:"name"`
}

func (r *CreateSizeCategoryRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}

func (r *UpdateSizeCategoryRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}
