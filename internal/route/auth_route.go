package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterAuthRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	jwtSecret := permissionMiddleware.GetJWTSecret()
	logger.Logger.Info("Initializing auth routes with JWT secret", "jwtSecretLength", len(jwtSecret))
	authHandler := handler.NewAuthHandler(jwtSecret)

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
