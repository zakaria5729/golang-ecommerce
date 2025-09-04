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
	registerAddressRoutes(mux)
	registerBrowsingHistoryRoutes(mux)
}

func registerCategoryRoutes(mux *http.ServeMux) {
	handler := handler.NewCategoryHandler()

	mux.HandleFunc("GET "+versionV1+"/categories", handler.GetAllCategories)
	mux.HandleFunc("GET "+versionV1+"/categories/paginated", handler.GetAllCategoriesPaginated)
	mux.HandleFunc("GET "+versionV1+"/categories/id/{id}", handler.GetCategoryByID)
	mux.HandleFunc("POST "+versionV1+"/categories", handler.CreateCategory)
	mux.HandleFunc("PUT "+versionV1+"/categories/id/{id}", handler.UpdateCategory)
	mux.HandleFunc("DELETE "+versionV1+"/categories/id/{id}", handler.DeleteCategory)
	mux.HandleFunc("PATCH "+versionV1+"/categories/id/{id}/toggle", handler.ToggleCategoryStatus)
}

func registerAddressRoutes(mux *http.ServeMux) {
	handler := handler.NewAddressHandler()

	mux.HandleFunc("GET "+versionV1+"/addresses", handler.GetAllAddresses)
	mux.HandleFunc("GET "+versionV1+"/addresses/paginated", handler.GetAllAddressesPaginated)
	mux.HandleFunc("GET "+versionV1+"/addresses/id/{id}", handler.GetAddressByID)
	mux.HandleFunc("POST "+versionV1+"/addresses", handler.CreateAddress)
	mux.HandleFunc("PUT "+versionV1+"/addresses/id/{id}", handler.UpdateAddress)
	mux.HandleFunc("DELETE "+versionV1+"/addresses/id/{id}", handler.DeleteAddress)
	mux.HandleFunc("PATCH "+versionV1+"/addresses/id/{id}/set-default", handler.SetDefaultAddress)
	mux.HandleFunc("GET "+versionV1+"/addresses/default", handler.GetDefaultAddress)
}

func registerBrowsingHistoryRoutes(mux *http.ServeMux) {
	handler := handler.NewBrowsingHistoryHandler()

	mux.HandleFunc("GET "+versionV1+"/browsing-history", handler.GetAllBrowsingHistory)
	mux.HandleFunc("GET "+versionV1+"/browsing-history/paginated", handler.GetAllBrowsingHistoryPaginated)
	mux.HandleFunc("GET "+versionV1+"/browsing-history/id/{id}", handler.GetBrowsingHistoryByID)
	mux.HandleFunc("POST "+versionV1+"/browsing-history", handler.CreateBrowsingHistory)
	mux.HandleFunc("DELETE "+versionV1+"/browsing-history/id/{id}", handler.DeleteBrowsingHistory)
	mux.HandleFunc("DELETE "+versionV1+"/browsing-history/clear", handler.ClearBrowsingHistory)
	mux.HandleFunc("GET "+versionV1+"/browsing-history/recent", handler.GetRecentBrowsingHistory)
	mux.HandleFunc("GET "+versionV1+"/browsing-history/most-viewed", handler.GetMostViewedProducts)
}
