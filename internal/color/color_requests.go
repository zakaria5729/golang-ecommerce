package color

import (
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateColorRequest struct {
	Name string `json:"name"`
}

type UpdateColorRequest struct {
	Name string `json:"name"`
}

func (r *CreateColorRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}

func (r *UpdateColorRequest) Sanitize() {
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}
