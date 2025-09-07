package auth

import (
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Permission struct {
	models.BaseModel
	Name        string  `json:"name" gorm:"uniqueIndex;not null;column:name"`
	Description *string `json:"description,omitempty" gorm:"column:description"`
}

const (
	PermissionName = "name"
)

func (p *Permission) Sanitize() {
	if p.Name != "" {
		p.Name = utils.Trim(p.Name)
	}
	if p.Description != nil && *p.Description != "" {
		sanitized := utils.Trim(*p.Description)
		p.Description = &sanitized
	}
}
