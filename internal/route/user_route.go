package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterUserRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewUserHandler()

	r.GET("/users/profile", h.GetProfile).Use(
		pm.RequireAuthWithRolePermission(),
	)

	r.PUT("/users/profile", h.UpdateProfile).Use(
		pm.RequireAuthUserStatus(),
	)

	r.POST("/users", h.CreateUser).Use(
		pm.RequirePermission(constants.PermissionUserCreate),
	)

	r.PUT("/users/{id}", h.UpdateUser).Use(
		pm.RequirePermission(constants.PermissionUserUpdate),
	)

	r.GET("/users", h.GetAllUsersPaginated).Use(
		pm.RequirePermission(constants.PermissionUserRead),
	)

	r.GET("/users/{id}", h.GetUserByID).Use(
		pm.RequirePermission(constants.PermissionUserRead),
	)

	r.DELETE("/users/{id}", h.DeleteUser).Use(
		pm.RequirePermission(constants.PermissionUserDelete),
	)
}

// func RegisterUserRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	userHandler := handler.NewUserHandler()

// 	mux.Handle(constants.GET+" /v1/users/profile", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuthWithRolePermission(),
// 	)(http.HandlerFunc(userHandler.GetProfile)))

// 	mux.Handle(constants.PUT+" /v1/users/profile", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuthUserId(),
// 	)(http.HandlerFunc(userHandler.UpdateProfile)))

// 	mux.Handle(constants.POST+" /v1/users", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionUserCreate),
// 	)(http.HandlerFunc(userHandler.CreateUser)))

// 	mux.Handle(constants.PUT+" /v1/users/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionUserUpdate),
// 	)(http.HandlerFunc(userHandler.UpdateUser)))

// 	mux.Handle(constants.GET+" /v1/users", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionUserRead),
// 	)(http.HandlerFunc(userHandler.GetAllUsersPaginated)))

// 	mux.Handle(constants.GET+" /v1/users/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionUserRead),
// 	)(http.HandlerFunc(userHandler.GetUserByID)))

// 	mux.Handle(constants.DELETE+" /v1/users/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionUserDelete),
// 	)(http.HandlerFunc(userHandler.DeleteUser)))
// }
