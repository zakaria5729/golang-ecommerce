package route

import (
	"net/http"

	"github.com/easy-comerce/backend/pkg/config"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAllRoutes(mux *http.ServeMux, cfg config.Config) {
	pm := m.NewPermissionMiddleware(cfg.JWTSecret)
	r := router.New(mux)
	r.Use(
		m.RecoveryMiddleware,
		m.CORSMiddleware,
		m.LoggingMiddleware,
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
