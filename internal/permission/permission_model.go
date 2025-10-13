package permission

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Permission struct {
	models.BaseModel
	Description *string `json:"description,omitempty" gorm:"column:description"`
	GroupName   *string `json:"group_name,omitempty" gorm:"column:group_name"`
	Name        string  `json:"name" gorm:"not null; column:name"`
}

func (Permission) TableName() string {
	return constants.TablePermission
}

func (p *Permission) Sanitize() {
	if p.Name != "" {
		p.Name = utils.Trim(p.Name)
	}
	if p.Description != nil && *p.Description != "" {
		sanitized := utils.Trim(*p.Description)
		p.Description = &sanitized
	}
	if p.GroupName != nil && *p.GroupName != "" {
		sanitized := utils.Trim(*p.GroupName)
		p.GroupName = &sanitized
	}
}
