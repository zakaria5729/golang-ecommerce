package route

import (
	"github.com/easy-comerce/backend/internal/system"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterSystemRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	service := system.NewSystemService(config.GetConfig())
	h := system.NewSystemHandler(service)

	r.POST("/system/health", h.SystemHealthCheck).Register()

	r.GET("/system/social-flow", h.HandleSocialFlowTemp).Register()
	r.GET("/system/social-flow/callback", h.HandleSocialFlowCallbackTemp).Register()
}
