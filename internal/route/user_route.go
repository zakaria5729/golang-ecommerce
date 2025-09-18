package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterUserRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
	userHandler := handler.NewUserHandler()

	mux.Handle(constants.GET+" /v1/users/profile", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuthWithRolePermission(),
	)(http.HandlerFunc(userHandler.GetProfile)))

	mux.Handle(constants.PUT+" /v1/users/profile", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(userHandler.UpdateProfile)))

	mux.Handle(constants.POST+" /v1/users", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionUserCreate),
	)(http.HandlerFunc(userHandler.CreateUser)))

	mux.Handle(constants.PUT+" /v1/users/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionUserUpdate),
	)(http.HandlerFunc(userHandler.UpdateUser)))

	mux.Handle(constants.GET+" /v1/users", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionUserRead),
	)(http.HandlerFunc(userHandler.GetAllUsersPaginated)))

	mux.Handle(constants.GET+" /v1/users/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionUserRead),
	)(http.HandlerFunc(userHandler.GetUserByID)))

	mux.Handle(constants.DELETE+" /v1/users/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionUserDelete),
	)(http.HandlerFunc(userHandler.DeleteUser)))
}
