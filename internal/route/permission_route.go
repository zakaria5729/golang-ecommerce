package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterPermissionRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewPermissionHandler()

	r.GET("/permissions", h.GetAllPermissions).Use(
		pm.RequirePermission(c.PermissionPermissionRead),
	).Register()

	r.GET("/permissions/{id}", h.GetPermissionByID).Use(
		pm.RequirePermission(c.PermissionPermissionRead),
	).Register()
}
