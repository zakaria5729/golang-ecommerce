package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterFileRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewFileStorageHandler()

	r.POST("/files/upload", h.UploadFile).Use(
		pm.RequireAuthUserStatus(),
	).Register()
}
