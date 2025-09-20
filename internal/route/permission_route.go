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
	)

	r.GET("/permissions/{id}", h.GetPermissionByID).Use(
		pm.RequirePermission(c.PermissionPermissionRead),
	)
}

// func RegisterPermissionRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	permissionHandler := handler.NewPermissionHandler()

// 	mux.Handle(constants.GET+" /v1/permissions", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionPermissionRead),
// 	)(http.HandlerFunc(permissionHandler.GetAllPermissions)))

// 	mux.Handle(constants.GET+" /v1/permissions/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionPermissionRead),
// 	)(http.HandlerFunc(permissionHandler.GetPermissionByID)))
// }
