package attribute_option

import (
	"github.com/easy-comerce/backend/internal/feature/attribute_type"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AttributeOption struct {
	models.BaseModel
	AttributeTypeID     *uint                         `gorm:"column:attribute_type_id"`
	AttributeOptionName string                        `gorm:"column:attribute_option_name"`
	AttributeType       *attribute_type.AttributeType `gorm:"foreignKey:AttributeTypeID;references:ID"`
}

func (AttributeOption) TableName() string {
	return constants.TableAttributeOption
}

func (ao *AttributeOption) Sanitize() {
	if ao.AttributeOptionName != "" {
		ao.AttributeOptionName = utils.Trim(ao.AttributeOptionName)
	}
}
