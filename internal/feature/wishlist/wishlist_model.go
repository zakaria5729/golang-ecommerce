package wishlist

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
)

type Wishlist struct {
	models.BaseModel
	UserID    uint `gorm:"not null; column:user_id"`
	ProductID uint `gorm:"not null; column:product_id"`
}

func (Wishlist) TableName() string {
	return constants.TableWishlist
}

const (
	WishlistUserID    = "user_id"
	WishlistProductID = "product_id"
)
