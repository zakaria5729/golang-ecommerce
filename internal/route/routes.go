package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
)

func RegisterRoutes(mux *http.ServeMux) {
	// Category routes
	categoryHandler := handler.NewCategoryHandler()
	mux.HandleFunc("/api/v1/categories", categoryHandler.GetAllCategories)
	mux.HandleFunc("/api/v1/categories/", categoryHandler.GetCategoryByID)
}
