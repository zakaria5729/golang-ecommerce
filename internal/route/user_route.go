package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterUserRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	userHandler := handler.NewUserHandler()

	// User profile routes (authentication required)
	mux.Handle("GET /v1/users/me", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(userHandler.GetMe)))

	mux.Handle("GET /v1/users/profile", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(userHandler.GetProfile)))

	mux.Handle("PUT /v1/users/profile", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(userHandler.UpdateProfile)))

	// User management routes (admin permissions required)
	mux.Handle("GET /v1/users", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("user.read"),
	)(http.HandlerFunc(userHandler.GetAllUsers)))

	mux.Handle("GET /v1/users/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("user.read"),
	)(http.HandlerFunc(userHandler.GetUserByID)))

	mux.Handle("DELETE /v1/users/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
		permissionMiddleware.RequirePermission("user.delete"),
	)(http.HandlerFunc(userHandler.DeleteUser)))
}
