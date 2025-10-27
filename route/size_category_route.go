package route

import (
	"github.com/easy-comerce/backend/db"
	sc "github.com/easy-comerce/backend/internal/size_category"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterSizeCategoryRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	repo := sc.NewSizeCategoryRepository(db.GetDB())
	service := sc.NewSizeCategoryService(repo)
	h := sc.NewSizeCategoryHandler(service)

	r.GET("/size-categories/public", h.GetAllSizeCategoriesPublic).Register()

	r.GET("/size-categories/{id}/public", h.GetSizeCategoryByIDPublic).Register()

	r.GET("/size-categories", h.GetAllSizeCategories).Use(
		pm.RequirePermission(c.PermissionSizeCategoryRead),
	).Register()

	r.GET("/size-categories/{id}", h.GetSizeCategoryByID).Use(
		pm.RequirePermission(c.PermissionSizeCategoryRead),
	).Register()

	r.POST("/size-categories", h.CreateSizeCategory).Use(
		pm.RequirePermission(c.PermissionSizeCategoryCreate),
	).Register()

	r.PUT("/size-categories/{id}", h.UpdateSizeCategory).Use(
		pm.RequirePermission(c.PermissionSizeCategoryUpdate),
	).Register()

	r.DELETE("/size-categories/{id}", h.DeleteSizeCategory).Use(
		pm.RequirePermission(c.PermissionSizeCategoryDelete),
	).Register()

	r.POST("/size-categories/{id}/undo", h.UndoDeletedSizeCategory).Use(
		pm.RequirePermission(c.PermissionSizeCategoryUndoDelete),
	).Register()
}
