package brand

import (
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateBrandRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type UpdateBrandRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func (r *CreateBrandRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}

func (r *UpdateBrandRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}
