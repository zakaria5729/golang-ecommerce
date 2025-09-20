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
	)
}

// func RegisterFileRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	fileHandler := handler.NewFileHandler()

// 	mux.Handle(constants.POST+" /v1/files/upload", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuthUserId(),
// 	)(http.HandlerFunc(fileHandler.UploadFile)))
// }
