package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterWishlistRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewWishlistHandler()

	r.GET("/wishlists/count/me", h.GetWishlistCountByUser).Use(
		pm.RequireAuthUserStatus(),
	)

	r.GET("/wishlists/me/paginated", h.GetAllWishlistsPaginatedByUser).Use(
		pm.RequireAuthUserStatus(),
	)

	r.POST("/wishlists/product/{product_id}", h.AddToWishlistsByUser).Use(
		pm.RequireAuthUserStatus(),
	)

	r.DELETE("/wishlists/product/{product_id}", h.RemoveFromWishlistByUser).Use(
		pm.RequireAuthUserStatus(),
	)

	r.DELETE("/wishlists/clear", h.ClearUserWishlist).Use(
		pm.RequireAuthUserStatus(),
	)

	r.GET("/wishlists/paginated", h.GetAllWishlistsPaginated).Use(
		pm.RequirePermission(c.PermissionWishlistRead),
	)

	// r.GET("/wishlists/product/{product_id}", h.GetWishlistByProduct).Use(
	// 	pm.RequireAuthUserStatus(),
	// )

	// r.GET("/wishlists/product/{product_id}/check", h.IsProductInWishlist).Use(
	// 	pm.RequirePermission(c.PermissionWishlistRead),
	// )

	// r.DELETE("/wishlists/id/{id}", h.DeleteWishlist).Use(
	// 	pm.RequirePermission(c.PermissionWishlistDelete),
	// )

	// r.GET("/wishlists", h.GetAllWishlists).Use(
	// 	pm.RequirePermission(c.PermissionWishlistRead),
	// )
}

// func RegisterWishlistRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	handler := handler.NewWishlistHandler()

// 	mux.Handle(constants.GET+" /v1/wishlists", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
// 	)(http.HandlerFunc(handler.GetAllWishlists)))

// 	mux.Handle(constants.GET+" /v1/wishlists/paginated", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
// 	)(http.HandlerFunc(handler.GetAllWishlistsPaginated)))

// 	mux.Handle(constants.GET+" /v1/wishlists/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
// 	)(http.HandlerFunc(handler.GetWishlistByID)))

// 	mux.Handle(constants.POST+" /v1/wishlists", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistCreate),
// 	)(http.HandlerFunc(handler.CreateWishlist)))

// 	mux.Handle(constants.DELETE+" /v1/wishlists/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistDelete),
// 	)(http.HandlerFunc(handler.DeleteWishlist)))

// 	mux.Handle(constants.DELETE+" /v1/wishlists/clear", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistDelete),
// 	)(http.HandlerFunc(handler.ClearWishlist)))

// 	mux.Handle(constants.GET+" /v1/wishlists/count", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
// 	)(http.HandlerFunc(handler.GetWishlistCount)))

// 	mux.Handle(constants.GET+" /v1/wishlists/product/{product_id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
// 	)(http.HandlerFunc(handler.GetWishlistByProduct)))

// 	mux.Handle(constants.DELETE+" /v1/wishlists/product/{product_id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistDelete),
// 	)(http.HandlerFunc(handler.DeleteWishlistByProduct)))

// 	mux.Handle(constants.GET+" /v1/wishlists/product/{product_id}/check", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
// 	)(http.HandlerFunc(handler.IsProductInWishlist)))
// }
