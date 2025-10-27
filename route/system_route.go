package route

import (
	"github.com/easy-comerce/backend/internal/system"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterSystemRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	service := system.NewSystemService(config.GetConfig())
	h := system.NewSystemHandler(service)

	r.POST("/system/health", h.SystemHealthCheck).Register()

	r.GET("/system/logs", h.GetSystemLogFiles).Use(
		pm.RequirePermission(c.PermissionSystemLogRead),
	).Register()

	r.GET("/system/logs/download/{file_name}", h.DownloadSystemLogFile).Use(
		pm.RequirePermission(c.PermissionSystemLogDownload),
	).Register()

	r.DELETE("/system/logs/delete/{file_name}", h.DeleteSystemLogFile).Use(
		pm.RequirePermission(c.PermissionSystemLogDelete),
	).Register()

	r.GET("/system/social-flow", h.HandleSocialFlowTemp).Register()
	r.GET("/system/social-flow/callback", h.HandleSocialFlowCallbackTemp).Register()
}
