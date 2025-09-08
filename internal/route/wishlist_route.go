package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterWishlistRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	handler := handler.NewWishlistHandler()

	mux.Handle("GET /v1/wishlists", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
	)(http.HandlerFunc(handler.GetAllWishlists)))

	mux.Handle("GET /v1/wishlists/paginated", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
	)(http.HandlerFunc(handler.GetAllWishlistsPaginated)))

	mux.Handle("GET /v1/wishlists/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
	)(http.HandlerFunc(handler.GetWishlistByID)))

	mux.Handle("POST /v1/wishlists", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistCreate),
	)(http.HandlerFunc(handler.CreateWishlist)))

	mux.Handle("DELETE /v1/wishlists/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistDelete),
	)(http.HandlerFunc(handler.DeleteWishlist)))

	mux.Handle("DELETE /v1/wishlists/clear", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistDelete),
	)(http.HandlerFunc(handler.ClearWishlist)))

	mux.Handle("GET /v1/wishlists/count", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
	)(http.HandlerFunc(handler.GetWishlistCount)))

	mux.Handle("GET /v1/wishlists/product/{product_id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
	)(http.HandlerFunc(handler.GetWishlistByProduct)))

	mux.Handle("DELETE /v1/wishlists/product/{product_id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistDelete),
	)(http.HandlerFunc(handler.DeleteWishlistByProduct)))

	mux.Handle("GET /v1/wishlists/product/{product_id}/check", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionWishlistRead),
	)(http.HandlerFunc(handler.IsProductInWishlist)))
}
