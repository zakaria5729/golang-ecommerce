package role

import (
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateRoleRequest struct {
	Description   *string `json:"description,omitempty"`
	PermissionIDs []uint  `json:"permission_ids"`
	RoleName      string  `json:"role_name"`
	RoleType      string  `json:"role_type"`
}

type UpdateRoleRequest struct {
	Description   *string `json:"description,omitempty"`
	PermissionIDs []uint  `json:"permission_ids,omitempty"`
	RoleName      string  `json:"role_name"`
	RoleType      string  `json:"role_type"`
}

type AssignRoleRequest struct {
	RoleId uint `json:"role_id"`
	UserID uint `json:"user_id"`
}

func (r *CreateRoleRequest) Sanitize() {
	if r.RoleName != "" {
		r.RoleName = utils.Trim(r.RoleName)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}

func (r *UpdateRoleRequest) Sanitize() {
	if r.RoleName != "" {
		r.RoleName = utils.Trim(r.RoleName)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}
