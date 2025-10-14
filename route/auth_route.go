package route

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/auth"
	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/user"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAuthRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	db := db.GetDB()
	cfg := config.GetConfig()
	userRepo := user.NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)
	service := auth.NewAuthService(cfg.JWTSecret, userRepo, roleRepo)
	h := auth.NewAuthHandler(service)

	r.POST("/auth/login", h.Login).Register()

	r.POST("/auth/social-login", h.SocialLogin).Register()

	r.POST("/auth/register", h.Register).Register()

	r.POST("/auth/logout/user/{user_id}", h.Logout).Register()

	r.POST("/auth/forgot-password", h.ForgotPassword).Register()

	r.POST("/auth/reset-password", h.ResetPassword).Register()

	r.POST("/auth/refresh-token", h.RefreshToken).Register()

	r.POST("/auth/resend-verify-link", h.ResendVerifyLink).Register()

	r.GET("/auth/verify-account", h.VerifyAccount).Register()
}
