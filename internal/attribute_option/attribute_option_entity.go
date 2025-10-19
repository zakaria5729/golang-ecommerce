package attribute_option

import (
	"github.com/easy-comerce/backend/internal/attribute_type"
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AttributeOptionEntity struct {
	base.BaseEntity
	AttributeTypeID     *uint                               `gorm:"column:attribute_type_id"`
	AttributeOptionName string                              `gorm:"column:attribute_option_name"`
	AttributeType       *attribute_type.AttributeTypeEntity `gorm:"foreignKey:AttributeTypeID;references:ID"`
}

func (AttributeOptionEntity) TableName() string {
	return constants.TableAttributeOption
}

func (ao *AttributeOptionEntity) Sanitize() {
	if ao.AttributeOptionName != "" {
		ao.AttributeOptionName = utils.Trim(ao.AttributeOptionName)
	}
}
