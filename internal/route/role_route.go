package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterRoleRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	roleHandler := handler.NewRoleHandler()

	// Role management routes (admin permissions required)
	mux.Handle("GET /v1/roles", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("role.read"),
	)(http.HandlerFunc(roleHandler.GetAllRoles)))

	mux.Handle("GET /v1/roles/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("role.read"),
	)(http.HandlerFunc(roleHandler.GetRoleByID)))

	mux.Handle("POST /v1/roles", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("role.create"),
	)(http.HandlerFunc(roleHandler.CreateRole)))

	mux.Handle("PUT /v1/roles/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("role.update"),
	)(http.HandlerFunc(roleHandler.UpdateRole)))

	mux.Handle("DELETE /v1/roles/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("role.delete"),
	)(http.HandlerFunc(roleHandler.DeleteRole)))

	// Role assignment routes (admin permissions required)
	mux.Handle("POST /v1/roles/assign", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("role.assign"),
	)(http.HandlerFunc(roleHandler.AssignRoleToUser)))
}
