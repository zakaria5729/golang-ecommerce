package size_option

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SizeOptionEntity struct {
	base.BaseEntity
	Name           string `gorm:"not null; column:name"`
	SortOrder      *int   `gorm:"column:sort_order"`
	SizeCategoryID uint   `gorm:"not null; column:size_category_id"`
}

func (SizeOptionEntity) TableName() string {
	return constants.TableSizeOption
}

func (so *SizeOptionEntity) Sanitize() {
	if so.Name != "" {
		so.Name = utils.Trim(so.Name)
	}
}
