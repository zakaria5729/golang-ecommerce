package route

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/color"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterColorRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := color.NewColorRepository(db.GetDB())
	service := color.NewColorService(repo)
	h := color.NewColorHandler(service)

	r.GET("/colors/public", h.GetAllColorsPublic).Register()

	r.GET("/colors/{id}/public", h.GetColorByIDPublic).Register()

	r.GET("/colors", h.GetAllColors).Use(
		pm.RequirePermission(c.PermissionColorRead),
	).Register()

	r.GET("/colors/{id}", h.GetColorByID).Use(
		pm.RequirePermission(c.PermissionColorRead),
	).Register()

	r.POST("/colors", h.CreateColor).Use(
		pm.RequirePermission(c.PermissionColorCreate),
	).Register()

	r.PUT("/colors/{id}", h.UpdateColor).Use(
		pm.RequirePermission(c.PermissionColorUpdate),
	).Register()

	r.DELETE("/colors/{id}", h.DeleteColor).Use(
		pm.RequirePermission(c.PermissionColorDelete),
	).Register()

	r.POST("/colors/{id}/undo", h.UndoDeletedColor).Use(
		pm.RequirePermission(c.PermissionColorUndoDelete),
	).Register()
}
