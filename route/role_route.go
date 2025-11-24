package route

import (
	"github.com/easy-comerce/backend/db"
	p "github.com/easy-comerce/backend/internal/permission"
	"github.com/easy-comerce/backend/internal/role"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterRoleRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	roleRepo := role.NewRoleRepository(db.GetDB())
	permissionRepo := p.NewPermissionRepository(db.GetDB())
	permissionService := p.NewPermissionService(permissionRepo)
	service := role.NewRoleService(roleRepo, permissionService)
	h := role.NewRoleHandler(service)

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

	r.POST("/roles/assign-to-user", h.AssignRoleToUser).Use(
		pm.RequirePermission(c.PermissionRoleAssign),
	).Register()

	r.POST("/roles/append-permissions", h.AppendPermissionsToRole).Use(
		pm.RequirePermission(c.PermissionRoleToAddPermissions),
	).Register()

	r.POST("/roles/{id}/undo", h.UndoDeletedRole).Use(
		pm.RequirePermission(c.PermissionRoleUndoDelete),
	).Register()
}
