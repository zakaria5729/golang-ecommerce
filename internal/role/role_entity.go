package role

import (
	"github.com/easy-comerce/backend/internal/permission"
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type RoleEntity struct {
	base.AuditEntity
	Description *string                       `json:"description" gorm:"column:description"`
	RoleName    string                        `json:"role_name" gorm:"not null; column:role_name"`
	RoleType    string                        `json:"role_type" gorm:"not null; column:role_type"`
	Permissions []permission.PermissionEntity `json:"permissions,omitempty" gorm:"many2many:role_permissions;joinForeignKey:role_id;joinReferences:permission_id"`
}

func (RoleEntity) TableName() string {
	return constants.TableRole
}

func (r *RoleEntity) Sanitize() {
	if r.RoleName != "" {
		r.RoleName = utils.Trim(r.RoleName)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}
