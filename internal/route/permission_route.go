package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterPermissionRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	permissionHandler := handler.NewPermissionHandler()

	// Permission management routes (admin permissions required)
	mux.Handle("GET /v1/permissions", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("permission.read"),
	)(http.HandlerFunc(permissionHandler.GetAllPermissions)))

	mux.Handle("GET /v1/permissions/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("permission.read"),
	)(http.HandlerFunc(permissionHandler.GetPermissionByID)))

	mux.Handle("POST /v1/permissions", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("permission.create"),
	)(http.HandlerFunc(permissionHandler.CreatePermission)))

	mux.Handle("PUT /v1/permissions/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("permission.update"),
	)(http.HandlerFunc(permissionHandler.UpdatePermission)))

	mux.Handle("DELETE /v1/permissions/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("permission.delete"),
	)(http.HandlerFunc(permissionHandler.DeletePermission)))
}
