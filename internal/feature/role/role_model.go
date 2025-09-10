package role

import (
	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Role struct {
	models.BaseModel
	RoleName    string                  `json:"role_name" gorm:"uniqueIndex:roles_role_name_key;not null;column:role_name"`
	RoleType    string                  `json:"role_type" gorm:"not null;column:role_type"`
	Description *string                 `json:"description,omitempty" gorm:"column:description"`
	Permissions []permission.Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

func (Role) TableName() string {
	return constants.TableRole
}

const (
	RoleRoleName = "role_name"
	RoleRoleType = "role_type"
)

func (r *Role) Sanitize() {
	if r.RoleName != "" {
		r.RoleName = utils.Trim(r.RoleName)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}

func (r *Role) IsValidRoleType() bool {
	switch r.RoleType {
	case constants.RoleTypeSuperAdmin, constants.RoleTypeAdmin, constants.RoleTypeManager, constants.RoleTypeSeller, constants.RoleTypeUser:
		return true
	default:
		return false
	}
}

func (r *Role) CanCreateRoles() bool {
	return r.RoleType == constants.RoleTypeSuperAdmin || r.RoleType == constants.RoleTypeAdmin
}

func (r *Role) CanManageUsers() bool {
	return r.RoleType == constants.RoleTypeSuperAdmin || r.RoleType == constants.RoleTypeAdmin || r.RoleType == constants.RoleTypeManager
}

func (r *Role) CanManageContent() bool {
	return r.RoleType == constants.RoleTypeSuperAdmin || r.RoleType == constants.RoleTypeAdmin || r.RoleType == constants.RoleTypeManager || r.RoleType == constants.RoleTypeSeller
}
