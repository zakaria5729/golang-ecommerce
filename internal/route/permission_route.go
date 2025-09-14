package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterPermissionRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
	permissionHandler := handler.NewPermissionHandler()

	mux.Handle(constants.GET+" /v1/permissions", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionPermissionRead),
	)(http.HandlerFunc(permissionHandler.GetAllPermissions)))

	mux.Handle(constants.GET+" /v1/permissions/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionPermissionRead),
	)(http.HandlerFunc(permissionHandler.GetPermissionByID)))
}
