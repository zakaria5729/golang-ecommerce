package route

import (
	"github.com/easy-comerce/backend/db"
	p "github.com/easy-comerce/backend/internal/permission"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterPermissionRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	repo := p.NewPermissionRepository(db.GetDB())
	service := p.NewPermissionService(repo)
	h := p.NewPermissionHandler(service)

	r.GET("/permissions", h.GetAllPermissions).Use(
		pm.RequirePermission(c.PermissionPermissionRead),
	).Register()

	r.GET("/permissions-group", h.GetAllPermissionsGroup).Use(
		pm.RequirePermission(c.PermissionPermissionRead),
	).Register()

	r.GET("/permissions/{id}", h.GetPermissionByID).Use(
		pm.RequirePermission(c.PermissionPermissionRead),
	).Register()
}
