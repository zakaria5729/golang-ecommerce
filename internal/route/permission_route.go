package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterPermissionRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	permissionHandler := handler.NewPermissionHandler()

	mux.Handle(constants.GET+" /v1/permissions", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionPermissionRead),
	)(http.HandlerFunc(permissionHandler.GetAllPermissions)))

	mux.Handle(constants.GET+" /v1/permissions/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionPermissionRead),
	)(http.HandlerFunc(permissionHandler.GetPermissionByID)))

	mux.Handle(constants.POST+" /v1/permissions", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionPermissionCreate),
	)(http.HandlerFunc(permissionHandler.CreatePermission)))

	mux.Handle(constants.PUT+" /v1/permissions/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionPermissionUpdate),
	)(http.HandlerFunc(permissionHandler.UpdatePermission)))

	mux.Handle(constants.DELETE+" /v1/permissions/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionPermissionDelete),
	)(http.HandlerFunc(permissionHandler.DeletePermission)))
}
