package attribute_option

import (
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateAttributeOptionRequest struct {
	AttributeTypeID     uint   `json:"attribute_type_id"`
	AttributeOptionName string `json:"attribute_option_name"`
}

type UpdateAttributeOptionRequest struct {
	AttributeTypeID     *uint  `json:"attribute_type_id,omitempty"`
	AttributeOptionName string `json:"attribute_option_name"`
}

func (r *CreateAttributeOptionRequest) Sanitize() {
	if r.AttributeOptionName != "" {
		r.AttributeOptionName = utils.Trim(r.AttributeOptionName)
	}
}

func (r *UpdateAttributeOptionRequest) Sanitize() {
	if r.AttributeOptionName != "" {
		r.AttributeOptionName = utils.Trim(r.AttributeOptionName)
	}
}
