package role

import (
	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Role struct {
	models.BaseModel
	Description *string                 `gorm:"column:description"`
	Permissions []permission.Permission `gorm:"many2many:role_permissions;"`
	RoleName    string                  `gorm:"not null; column:role_name"`
	RoleType    string                  `gorm:"not null; column:role_type"`
}

func (Role) TableName() string {
	return constants.TableRole
}

func (r *Role) Sanitize() {
	if r.RoleName != "" {
		r.RoleName = utils.Trim(r.RoleName)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}
