package attribute_type

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AttributeTypeEntity struct {
	base.AuditEntity
	Name string `gorm:"column:name"`
}

func (AttributeTypeEntity) TableName() string {
	return constants.TableAttributeType
}

func (at *AttributeTypeEntity) Sanitize() {
	if at.Name != "" {
		at.Name = utils.Trim(at.Name)
	}
}
