package browsing_history

import (
	"time"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
)

type BrowsingHistory struct {
	models.BaseModel
	UserID    uint      `json:"user_id" gorm:"not null; column:user_id"`
	ProductID uint      `json:"product_id" gorm:"not null; column:product_id"`
	ViewedAt  time.Time `json:"viewed_at" gorm:"column:viewed_at"`
}

func (BrowsingHistory) TableName() string {
	return constants.TableBrowsingHistory
}

const (
	BrowsingHistoryUserID    = "user_id"
	BrowsingHistoryProductID = "product_id"
	BrowsingHistoryViewedAt  = "viewed_at"
)
