package attribute_type

import (
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateAttributeTypeRequest struct {
	Name string `json:"name"`
}

type UpdateAttributeTypeRequest struct {
	Name string `json:"name"`
}

func (r *CreateAttributeTypeRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}

func (r *UpdateAttributeTypeRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}
