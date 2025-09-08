package role

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type CreateRoleRequest struct {
	RoleName    string   `json:"role_name"`
	RoleType    string   `json:"role_type"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type AssignRoleRequest struct {
	UserID uint     `json:"user_id"`
	Roles  []string `json:"roles"`
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

func (r *CreateRoleRequest) IsValidRoleType() bool {
	switch r.RoleType {
	case constants.RoleTypeAdmin, constants.RoleTypeManager, constants.RoleTypeSeller, constants.RoleTypeUser:
		return true
	default:
		return false
	}
}

func (r *CreateRoleRequest) CanCreateRoles() bool {
	return r.RoleType == constants.RoleTypeSuperAdmin || r.RoleType == constants.RoleTypeAdmin
}
