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
	)

	r.GET("/roles/{id}", h.GetRoleByID).Use(
		pm.RequirePermission(c.PermissionRoleRead),
	)

	r.POST("/roles", h.CreateRole).Use(
		pm.RequirePermission(c.PermissionRoleCreate),
	)

	r.PUT("/roles/{id}", h.UpdateRole).Use(
		pm.RequirePermission(c.PermissionRoleUpdate),
	)

	r.DELETE("/roles/{id}", h.DeleteRole).Use(
		pm.RequirePermission(c.PermissionRoleDelete),
	)

	r.POST("/roles/assign", h.AssignRoleToUser).Use(
		pm.RequirePermission(c.PermissionRoleAssign),
	)
}

// func RegisterRoleRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	roleHandler := handler.NewRoleHandler()

// 	mux.Handle(constants.GET+" /v1/roles", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionRoleRead),
// 	)(http.HandlerFunc(roleHandler.GetAllRoles)))

// 	mux.Handle(constants.GET+" /v1/roles/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionRoleRead),
// 	)(http.HandlerFunc(roleHandler.GetRoleByID)))

// 	mux.Handle(constants.POST+" /v1/roles", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionRoleCreate),
// 	)(http.HandlerFunc(roleHandler.CreateRole)))

// 	mux.Handle(constants.PUT+" /v1/roles/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionRoleUpdate),
// 	)(http.HandlerFunc(roleHandler.UpdateRole)))

// 	mux.Handle(constants.DELETE+" /v1/roles/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionRoleDelete),
// 	)(http.HandlerFunc(roleHandler.DeleteRole)))

// 	mux.Handle(constants.POST+" /v1/roles/assign", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionRoleAssign),
// 	)(http.HandlerFunc(roleHandler.AssignRoleToUser)))
// }
