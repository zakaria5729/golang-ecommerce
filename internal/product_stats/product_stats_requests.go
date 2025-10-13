package product_stats

type IncreaseProductStatsRequest struct {
	ProductID                       uint  `json:"product_id"`
	IncreaseViewCount               *bool `json:"view_count"`
	IncreaseAddToCartCount          *bool `json:"add_to_cart_count"`
	IncreaseRemoveFromCartCount     *bool `json:"remove_from_cart_count"`
	IncreaseWishlistCount           *bool `json:"wishlist_count"`
	IncreaseRemoveFromWishlistCount *bool `json:"remove_from_wishlist_count"`
}
