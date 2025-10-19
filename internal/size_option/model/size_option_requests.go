package model

import (
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateSizeOptionRequest struct {
	Name           string `json:"name"`
	SortOrder      *int   `json:"sort_order"`
	SizeCategoryID uint   `json:"size_category_id"`
}

type UpdateSizeOptionRequest struct {
	Name      string `json:"name"`
	SortOrder *int   `json:"sort_order"`
}

func (r *CreateSizeOptionRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}

func (r *UpdateSizeOptionRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}
