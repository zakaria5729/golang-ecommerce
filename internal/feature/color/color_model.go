package color

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Color struct {
	models.BaseModel
	Name string `gorm:"not null; column:name"`
}

func (Color) TableName() string {
	return constants.TableColor
}

func (c *Color) Sanitize() {
	if c.Name != "" {
		c.Name = utils.Trim(c.Name)
	}
}
