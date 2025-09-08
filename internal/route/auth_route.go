package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterAuthRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	authHandler := handler.NewAuthHandler(permissionMiddleware.GetJWTSecret())

	// Public auth routes (no authentication required)
	mux.HandleFunc("POST /v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /v1/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("POST /v1/auth/refresh-token", authHandler.RefreshToken)

	// Protected auth routes (authentication required)
	mux.Handle("POST /v1/auth/change-password", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(authHandler.ChangePassword)))

	mux.Handle("POST /v1/auth/logout", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(authHandler.Logout)))
}
