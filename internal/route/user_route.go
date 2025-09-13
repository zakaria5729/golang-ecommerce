package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterUserRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	userHandler := handler.NewUserHandler()

	mux.Handle(constants.GET+" /v1/users/me", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
	)(http.HandlerFunc(userHandler.GetMe)))

	mux.Handle(constants.GET+" /v1/users/profile", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
	)(http.HandlerFunc(userHandler.GetProfile)))

	mux.Handle(constants.PUT+" /v1/users/profile", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
	)(http.HandlerFunc(userHandler.UpdateProfile)))

	mux.Handle(constants.GET+" /v1/users", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionUserRead),
	)(http.HandlerFunc(userHandler.GetAllUsers)))

	mux.Handle(constants.GET+" /v1/users/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionUserRead),
	)(http.HandlerFunc(userHandler.GetUserByID)))

	mux.Handle(constants.DELETE+" /v1/users/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(false, false),
		permissionMiddleware.RequirePermission(constants.PermissionUserDelete),
	)(http.HandlerFunc(userHandler.DeleteUser)))
}
