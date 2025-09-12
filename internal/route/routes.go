package route

import (
	"net/http"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterRoutes(mux *http.ServeMux, cfg config.Config) {
	logger.Logger.Info("Initializing routes with JWT secret", "jwtSecretLength", len(cfg.JWTSecret))
	permissionMiddleware := middleware.NewAuthPermissionMiddleware(cfg.JWTSecret)

	// Register all route modules
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
