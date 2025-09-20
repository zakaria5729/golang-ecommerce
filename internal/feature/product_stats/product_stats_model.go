package product_stats

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
)

type ProductStats struct {
	models.BaseModel
	ProductID               uint `gorm:"not null; column:product_id"`
	ViewCount               uint `gorm:"not null; default:0; column:view_count"`
	AddToCartCount          uint `gorm:"not null; default:0; column:add_to_cart_count"`
	RemoveFromCartCount     uint `gorm:"not null; default:0; column:remove_from_cart_count"`
	PurchaseCount           uint `gorm:"not null; default:0; column:purchase_count"`
	WishlistCount           uint `gorm:"not null; default:0; column:wishlist_count"`
	RemoveFromWishlistCount uint `gorm:"not null; default:0; column:remove_from_wishlist_count"`
}

func (ProductStats) TableName() string {
	return constants.TableProductStats
}
