package route

import (
	"net/http"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAllRoutes(mux *http.ServeMux, cfg config.Config) {
	pm := middleware.NewPermissionMiddleware(cfg.JWTSecret)

	r := router.New(mux)
	r.Use(
		middleware.RecoveryMiddleware,
		middleware.CORSMiddleware,
	)

	RegisterAuthRoute(r, pm)
	RegisterUserRoute(r, pm)
	RegisterRoleRoute(r, pm)
	RegisterPermissionRoute(r, pm)
	RegisterCategoryRoute(r, pm)
	RegisterAddressRoute(r, pm)
	RegisterProductStatsRoute(r, pm)
	RegisterReviewRoute(r, pm)
	RegisterWishlistRoute(r, pm)
	RegisterFileRoute(r, pm)
}
