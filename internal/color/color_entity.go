package color

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ColorEntity struct {
	base.BaseEntity
	Name string `gorm:"not null; column:name"`
}

func (ColorEntity) TableName() string {
	return constants.TableColor
}

func (c *ColorEntity) Sanitize() {
	if c.Name != "" {
		c.Name = utils.Trim(c.Name)
	}
}
