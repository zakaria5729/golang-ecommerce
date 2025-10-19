package wishlist

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
)

type WishlistEntity struct {
	base.BaseEntity
	UserID    uint `gorm:"not null; column:user_id"`
	ProductID uint `gorm:"not null; column:product_id"`
}

func (WishlistEntity) TableName() string {
	return constants.TableWishlist
}
