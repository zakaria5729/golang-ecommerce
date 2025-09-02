package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
)

const (
	versionV1 = "/v1"
)

func RegisterRoutes(mux *http.ServeMux) {
	registerCategoryRoutes(mux)
}

func registerCategoryRoutes(mux *http.ServeMux) {
	handler := handler.NewCategoryHandler()

	mux.HandleFunc("GET "+versionV1+"/categories", handler.GetAllCategories)
	mux.HandleFunc("GET "+versionV1+"/categories/paginated", handler.GetAllCategoriesPaginated)
	mux.HandleFunc("GET "+versionV1+"/categories/{id}", handler.GetCategoryByID)
	mux.HandleFunc("POST "+versionV1+"/categories", handler.CreateCategory)
	mux.HandleFunc("PUT "+versionV1+"/categories/{id}", handler.UpdateCategory)
	mux.HandleFunc("DELETE "+versionV1+"/categories/{id}", handler.DeleteCategory)
	mux.HandleFunc("PATCH "+versionV1+"/categories/{id}/toggle", handler.ToggleCategoryStatus)
}
