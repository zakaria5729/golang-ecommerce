package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAuthRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewAuthHandler(pm.GetJWTSecret())

	r.GET("/health", h.HealthCheck).Register()

	r.POST("/auth/login", h.Login).Register()

	r.POST("/auth/register", h.Register).Register()

	r.POST("/auth/logout/user/{user_id}", h.Logout).Register()

	r.POST("/auth/refresh-token", h.RefreshToken).Register()

	r.POST("/auth/reset-password", h.ResetPassword).Register()

	r.POST("/auth/forgot-password", h.ForgotPassword).Register()
}
