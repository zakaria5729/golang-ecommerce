package attribute_type

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AttributeType struct {
	models.BaseModel
	Name string `gorm:"column:name"`
}

func (AttributeType) TableName() string {
	return constants.TableAttributeType
}

func (at *AttributeType) Sanitize() {
	if at.Name != "" {
		at.Name = utils.Trim(at.Name)
	}
}
