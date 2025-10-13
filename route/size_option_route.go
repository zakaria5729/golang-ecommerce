package route

import (
	"github.com/easy-comerce/backend/db"
	so "github.com/easy-comerce/backend/internal/size_option"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterSizeOptionRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := so.NewSizeOptionRepository(db.GetDB())
	service := so.NewSizeOptionService(repo)
	h := so.NewSizeOptionHandler(service)

	r.GET("/size-options/public", h.GetAllSizeOptionsPublic).Register()

	r.GET("/size-options/{id}/public", h.GetSizeOptionByIDPublic).Register()

	r.GET("/size-options", h.GetAllSizeOptions).Use(
		pm.RequirePermission(c.PermissionSizeOptionRead),
	).Register()

	r.GET("/size-options/{id}", h.GetSizeOptionByID).Use(
		pm.RequirePermission(c.PermissionSizeOptionRead),
	).Register()

	r.POST("/size-options", h.CreateSizeOption).Use(
		pm.RequirePermission(c.PermissionSizeOptionCreate),
	).Register()

	r.PUT("/size-options/{id}", h.UpdateSizeOption).Use(
		pm.RequirePermission(c.PermissionSizeOptionUpdate),
	).Register()

	r.DELETE("/size-options/{id}", h.DeleteSizeOption).Use(
		pm.RequirePermission(c.PermissionSizeOptionDelete),
	).Register()

	r.POST("/size-options/{id}/undo", h.UndoDeletedSizeOption).Use(
		pm.RequirePermission(c.PermissionSizeOptionUndoDelete),
	).Register()
}
