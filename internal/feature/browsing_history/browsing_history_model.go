package browsing_history

import (
	"time"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
)

type BrowsingHistory struct {
	models.BaseModel
	ViewedAt  time.Time `gorm:"column:viewed_at"`
	UserID    uint      `gorm:"not null; column:user_id"`
	ProductID uint      `gorm:"not null; column:product_id"`
}

func (BrowsingHistory) TableName() string {
	return constants.TableBrowsingHistory
}
