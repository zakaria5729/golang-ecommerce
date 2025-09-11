package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterCategoryRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	handler := handler.NewCategoryHandler()

	mux.HandleFunc("GET /v1/categories", handler.GetAllCategories)
	mux.HandleFunc("GET /v1/categories/paginated", handler.GetAllCategoriesPaginated)
	mux.HandleFunc("GET /v1/categories/id/{id}", handler.GetCategoryByID)

	mux.Handle("POST /v1/categories", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionCategoryCreate),
	)(http.HandlerFunc(handler.CreateCategory)))

	mux.Handle("PUT /v1/categories/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionCategoryUpdate),
	)(http.HandlerFunc(handler.UpdateCategory)))

	mux.Handle("DELETE /v1/categories/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionCategoryDelete),
	)(http.HandlerFunc(handler.DeleteCategory)))

	mux.Handle("PATCH /v1/categories/id/{id}/toggle", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionCategoryUpdate),
	)(http.HandlerFunc(handler.ToggleCategoryStatus)))
}
