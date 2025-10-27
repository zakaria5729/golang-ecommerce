package route

import (
	"github.com/easy-comerce/backend/db"
	b "github.com/easy-comerce/backend/internal/brand"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterBrandRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	repo := b.NewBrandRepository(db.GetDB())
	service := b.NewBrandService(repo)
	h := b.NewBrandHandler(service)

	r.GET("/brands/public", h.GetAllBrandsPublic).Register()

	r.GET("/brands/{id}/public", h.GetBrandByIDPublic).Register()

	r.GET("/brands", h.GetAllBrands).Use(
		pm.RequirePermission(c.PermissionBrandRead),
	).Register()

	r.GET("/brands/{id}", h.GetBrandByID).Use(
		pm.RequirePermission(c.PermissionBrandRead),
	).Register()

	r.POST("/brands", h.CreateBrand).Use(
		pm.RequirePermission(c.PermissionBrandCreate),
	).Register()

	r.PUT("/brands/{id}", h.UpdateBrand).Use(
		pm.RequirePermission(c.PermissionBrandUpdate),
	).Register()

	r.DELETE("/brands/{id}", h.DeleteBrand).Use(
		pm.RequirePermission(c.PermissionBrandDelete),
	).Register()

	r.POST("/brands/{id}/undo", h.UndoDeletedBrand).Use(
		pm.RequirePermission(c.PermissionBrandUndoDelete),
	).Register()
}
