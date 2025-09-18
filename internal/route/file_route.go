package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterFileRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
	fileHandler := handler.NewFileHandler()

	mux.Handle(constants.POST+" /v1/files/upload", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(fileHandler.UploadFile)))
}
