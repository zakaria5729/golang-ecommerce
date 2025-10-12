package route

import (
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAllRoutes(r *router.Router) {
	pm := m.NewPermissionMiddleware()
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
