package wishlist

type CreateWishlistRequest struct {
	ProductID uint `json:"product_id" validate:"required,min=1"`
}
