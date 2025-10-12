package route

import (
	fs "github.com/easy-comerce/backend/internal/feature/file_storage"
	"github.com/easy-comerce/backend/internal/handler"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterFileRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	storage, err := fs.NewObjectStorage()
	if err != nil {
		l.Logger.Error("Failed to create object storage client", "error", err)
		return
	}

	repo := fs.NewFileStorageRepository(storage)
	service := fs.NewFileStorageService(repo)
	h := handler.NewFileStorageHandler(service)

	r.POST("/files/upload", h.UploadFile).Use(
		pm.RequireAuthUserStatus(),
	).Register()
}
