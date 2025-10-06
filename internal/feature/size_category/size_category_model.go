package size_category

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SizeCategory struct {
	models.BaseModel
	Name string `gorm:"not null; column:name"`
}

func (SizeCategory) TableName() string {
	return constants.TableSizeCategory
}

func (sc *SizeCategory) Sanitize() {
	if sc.Name != "" {
		sc.Name = utils.Trim(sc.Name)
	}
}
