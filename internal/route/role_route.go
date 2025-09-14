package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterRoleRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
	roleHandler := handler.NewRoleHandler()

	mux.Handle(constants.GET+" /v1/roles", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionRoleRead),
	)(http.HandlerFunc(roleHandler.GetAllRoles)))

	mux.Handle(constants.GET+" /v1/roles/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionRoleRead),
	)(http.HandlerFunc(roleHandler.GetRoleByID)))

	mux.Handle(constants.POST+" /v1/roles", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionRoleCreate),
	)(http.HandlerFunc(roleHandler.CreateRole)))

	mux.Handle(constants.PUT+" /v1/roles/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionRoleUpdate),
	)(http.HandlerFunc(roleHandler.UpdateRole)))

	mux.Handle(constants.DELETE+" /v1/roles/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionRoleDelete),
	)(http.HandlerFunc(roleHandler.DeleteRole)))

	mux.Handle(constants.POST+" /v1/roles/assign", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionRoleAssign),
	)(http.HandlerFunc(roleHandler.AssignRoleToUser)))
}
