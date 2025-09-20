package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAuthRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewAuthHandler(pm.GetJWTSecret())

	r.POST("/auth/login", h.Login)

	r.POST("/auth/register", h.Register)

	r.POST("/auth/logout/id/{id}", h.Logout)

	r.POST("/auth/refresh-token", h.RefreshToken)

	r.POST("/auth/reset-password", h.ResetPassword)

	r.POST("/auth/forgot-password", h.ForgotPassword)

	r.POST("/auth/change-password", h.ChangePassword).Use(
		pm.RequireAuthUser(),
	)
}

// func RegisterAuthRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	jwtSecret := permissionMiddleware.GetJWTSecret()
// 	authHandler := handler.NewAuthHandler(jwtSecret)

// 	mux.HandleFunc(constants.POST+" /v1/auth/login", authHandler.Login)
// 	mux.HandleFunc(constants.POST+" /v1/auth/register", authHandler.Register)
// 	mux.HandleFunc(constants.POST+" /v1/auth/forgot-password", authHandler.ForgotPassword)
// 	mux.HandleFunc(constants.POST+" /v1/auth/reset-password", authHandler.ResetPassword)
// 	mux.HandleFunc(constants.POST+" /v1/auth/refresh-token", authHandler.RefreshToken)
// 	mux.HandleFunc(constants.POST+" /v1/auth/logout/id/{id}", authHandler.Logout)

// 	mux.Handle(constants.POST+" /v1/auth/change-password", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuthUser(),
// 	)(http.HandlerFunc(authHandler.ChangePassword)))
// }
