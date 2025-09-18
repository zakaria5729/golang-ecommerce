package route

import (
	"net/http"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterRoutes(mux *http.ServeMux, cfg config.Config) {
	permissionMiddleware := middleware.NewPermissionMiddleware(cfg.JWTSecret)

	RegisterAuthRoute(mux, permissionMiddleware)
	RegisterUserRoute(mux, permissionMiddleware)
	RegisterRoleRoute(mux, permissionMiddleware)
	RegisterPermissionRoute(mux, permissionMiddleware)
	RegisterCategoryRoute(mux, permissionMiddleware)
	RegisterAddressRoute(mux, permissionMiddleware)
	RegisterBrowsingHistoryRoute(mux, permissionMiddleware)
	RegisterReviewRoute(mux, permissionMiddleware)
	RegisterWishlistRoute(mux, permissionMiddleware)
	RegisterFileRoute(mux, permissionMiddleware)
}
