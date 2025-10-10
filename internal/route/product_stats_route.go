package route

import (
	"github.com/easy-comerce/backend/db"
	ps "github.com/easy-comerce/backend/internal/feature/product_stats"
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterProductStatsRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := ps.NewProductStatsRepository(db.GetDB())
	service := ps.NewProductStatsService(repo)
	h := handler.NewProductStatsHandler(service)

	r.GET("/product-stats/paginated", h.GetAllProductStatsPaginated).Use(
		pm.RequirePermission(c.PermissionProductStatsRead),
	).Register()

	r.GET("/product-stats/{id}", h.GetProductStatsByID).Use(
		pm.RequirePermission(c.PermissionProductStatsRead),
	).Register()

	r.POST("/product-stats", h.IncreaseProductStats).Use(
		pm.RequireAuthUserStatus(),
	).Register()
}
