package route

import (
	"net/http"

	m 	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAllRoutes(mux *http.ServeMux) {
	pm := m.NewPermissionMiddleware()
	r := router.New(mux)
	r.Use(
		m.RecoveryMiddleware,
		m.LoggingMiddleware,
		m.CORSMiddleware,
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
	RegisterBrandRoute(r, pm)
	RegisterColorRoute(r, pm)
	RegisterAttributeTypeRoute(r, pm)
	RegisterAttributeOptionRoute(r, pm)
	RegisterSizeCategoryRoute(r, pm)
	RegisterSizeOptionRoute(r, pm)
	RegisterNotificationRoute(r, pm)
}
