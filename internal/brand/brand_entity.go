package brand

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type BrandEntity struct {
	base.AuditEntity
	Name        string  `gorm:"not null; column:name"`
	Description *string `gorm:"column:description"`
}

func (BrandEntity) TableName() string {
	return constants.TableBrand
}

func (b *BrandEntity) Sanitize() {
	if b.Name != "" {
		b.Name = utils.Trim(b.Name)
	}
	if b.Description != nil && *b.Description != "" {
		sanitized := utils.Trim(*b.Description)
		b.Description = &sanitized
	}
}
