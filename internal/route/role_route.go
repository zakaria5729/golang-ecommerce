package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterRoleRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewRoleHandler()

	r.GET("/roles", h.GetAllRoles).Use(
		pm.RequirePermission(c.PermissionRoleRead),
	).Register()

	r.GET("/roles/{id}", h.GetRoleByID).Use(
		pm.RequirePermission(c.PermissionRoleRead),
	).Register()

	r.POST("/roles", h.CreateRole).Use(
		pm.RequirePermission(c.PermissionRoleCreate),
	).Register()

	r.PUT("/roles/{id}", h.UpdateRole).Use(
		pm.RequirePermission(c.PermissionRoleUpdate),
	).Register()

	r.DELETE("/roles/{id}", h.DeleteRole).Use(
		pm.RequirePermission(c.PermissionRoleDelete),
	).Register()

	r.POST("/roles/assign", h.AssignRoleToUser).Use(
		pm.RequirePermission(c.PermissionRoleAssign),
	).Register()

	r.POST("/roles/{id}/undo", h.UndoDeletedRole).Use(
		pm.RequirePermission(c.PermissionRoleUndoDelete),
	).Register()
}
