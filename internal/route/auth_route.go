package route

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAuthRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	db := db.GetDB()
	userRepo := user.NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)
	service := auth.NewAuthService(pm.GetJWTSecret(), userRepo, roleRepo)
	h := handler.NewAuthHandler(service)

	r.POST("/app-health", h.AppHealthCheck).Register()

	r.POST("/auth/login", h.Login).Register()

	r.POST("/auth/social-login", h.SocialLogin).Register()

	r.POST("/auth/register", h.Register).Register()

	r.POST("/auth/logout/user/{user_id}", h.Logout).Register()

	r.POST("/auth/refresh-token", h.RefreshToken).Register()

	r.POST("/auth/forgot-password", h.ForgotPassword).Register()

	r.POST("/auth/reset-password", h.ResetPassword).Register()

	r.GET("/auth/social-flow", auth.HandleSocialFlowTemp).Register()
	r.GET("/auth/social-flow/callback", auth.HandleSocialFlowCallbackTemp).Register()
}
