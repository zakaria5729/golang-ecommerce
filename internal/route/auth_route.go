package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterAuthRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
	jwtSecret := permissionMiddleware.GetJWTSecret()
	authHandler := handler.NewAuthHandler(jwtSecret)

	mux.HandleFunc(constants.POST+" /v1/auth/login", authHandler.Login)
	mux.HandleFunc(constants.POST+" /v1/auth/register", authHandler.Register)
	mux.HandleFunc(constants.POST+" /v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc(constants.POST+" /v1/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc(constants.POST+" /v1/auth/refresh-token", authHandler.RefreshToken)
	mux.HandleFunc(constants.POST+" /v1/auth/logout/id/{id}", authHandler.Logout)

	mux.Handle(constants.POST+" /v1/auth/change-password", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(authHandler.ChangePassword)))
}
