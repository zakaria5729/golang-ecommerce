package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterProductStatsRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewProductStatsHandler()

	r.GET("/product-stats/paginated", h.GetAllProductStatsPaginated).Use(
		pm.RequirePermission(c.PermissionProductStatsRead),
	).Register()

	r.GET("/product-stats/id/{id}", h.GetProductStatsByID).Use(
		pm.RequirePermission(c.PermissionProductStatsRead),
	).Register()

	r.POST("/product-stats", h.IncreaseProductStats).Use(
		pm.RequireAuthUserStatus(),
	).Register()
}
