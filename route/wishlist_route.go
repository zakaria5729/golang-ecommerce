package route

import (
	"github.com/easy-comerce/backend/db"
	w "github.com/easy-comerce/backend/internal/wishlist"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterWishlistRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := w.NewWishlistRepository(db.GetDB())
	service := w.NewWishlistService(repo)
	h := w.NewWishlistHandler(service)

	r.GET("/wishlists/count/me", h.GetWishlistCountByUser).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.GET("/wishlists/me/paginated", h.GetAllWishlistsPaginatedByUser).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.POST("/wishlists/product/{product_id}", h.AddToWishlistsByUser).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.DELETE("/wishlists/product/{product_id}", h.RemoveFromWishlistByUser).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.DELETE("/wishlists/clear", h.ClearUserWishlist).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.GET("/wishlists/paginated", h.GetAllWishlistsPaginated).Use(
		pm.RequirePermission(c.PermissionWishlistRead),
	).Register()

	r.DELETE("/wishlists/delete", h.DeleteWishlistById).Use(
		pm.RequirePermission(c.PermissionWishlistDelete),
	).Register()

	r.POST("/wishlists/delete/undo", h.UndoDeleteWishlistById).Use(
		pm.RequirePermission(c.PermissionWishlistUndoDelete),
	).Register()
}
