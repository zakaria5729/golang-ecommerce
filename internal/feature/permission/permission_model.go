package permission

import (
	"errors"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Permission struct {
	models.BaseModel
	Name        string  `json:"name" gorm:"uniqueIndex;not null;column:name"`
	Description *string `json:"description,omitempty" gorm:"column:description"`
}

func (Permission) TableName() string {
	return constants.TablePermission
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

func (p *Permission) IsValid() error {
	if p.Name == "" {
		return errors.New("permission name is required")
	}
	if len(p.Name) < 2 {
		return errors.New("permission name must be at least 2 characters long")
	}
	return nil
}
