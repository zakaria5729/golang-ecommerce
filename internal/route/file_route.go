package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterFileRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	fileHandler := handler.NewFileHandler()

	// Public endpoint - no authentication required
	mux.HandleFunc("GET /v1/files/{key}/url", fileHandler.GetFileURL)

	// Protected endpoints - require authentication
	mux.Handle("POST /v1/files/upload", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(fileHandler.UploadFile)))

	mux.Handle("DELETE /v1/files/{key}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(fileHandler.DeleteFile)))
}
