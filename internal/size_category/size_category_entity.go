package size_category

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SizeCategoryEntity struct {
	base.AuditEntity
	Name string `gorm:"not null; column:name"`
}

func (SizeCategoryEntity) TableName() string {
	return constants.TableSizeCategory
}

func (sc *SizeCategoryEntity) Sanitize() {
	if sc.Name != "" {
		sc.Name = utils.Trim(sc.Name)
	}
}
