package size_option

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SizeOption struct {
	models.BaseModel
	Name           string `gorm:"not null;column:name"`
	SortOrder      *int   `gorm:"column:sort_order"`
	SizeCategoryID uint   `gorm:"not null;column:size_category_id"`
}

func (SizeOption) TableName() string {
	return constants.TableSizeOption
}

func (so *SizeOption) Sanitize() {
	if so.Name != "" {
		so.Name = utils.Trim(so.Name)
	}
}
