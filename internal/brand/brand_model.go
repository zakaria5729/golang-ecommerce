package brand

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Brand struct {
	models.BaseModel
	Name        string  `gorm:"not null; column:name"`
	Description *string `gorm:"column:description"`
}

func (Brand) TableName() string {
	return constants.TableBrand
}

func (b *Brand) Sanitize() {
	if b.Name != "" {
		b.Name = utils.Trim(b.Name)
	}
	if b.Description != nil && *b.Description != "" {
		sanitized := utils.Trim(*b.Description)
		b.Description = &sanitized
	}
}
