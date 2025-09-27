package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAttributeOptionRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewAttributeOptionHandler()

	r.GET("/attribute-options/public", h.GetAllAttributeOptionsPublic).Register()

	r.GET("/attribute-options/{id}/public", h.GetAttributeOptionByIDPublic).Register()

	r.GET("/attribute-options", h.GetAllAttributeOptions).Use(
		pm.RequirePermission(c.PermissionAttributeOptionRead),
	).Register()

	r.GET("/attribute-options/{id}", h.GetAttributeOptionByID).Use(
		pm.RequirePermission(c.PermissionAttributeOptionRead),
	).Register()

	r.POST("/attribute-options", h.CreateAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionCreate),
	).Register()

	r.PUT("/attribute-options/{id}", h.UpdateAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionUpdate),
	).Register()

	r.DELETE("/attribute-options/{id}", h.DeleteAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionDelete),
	).Register()

	r.POST("/attribute-options/{id}/undo", h.UndoDeletedAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionUndoDelete),
	).Register()
}
