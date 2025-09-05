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
	registerReviewRoutes(mux)
	registerWishlistRoutes(mux)
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

func registerReviewRoutes(mux *http.ServeMux) {
	handler := handler.NewReviewHandler()

	mux.HandleFunc("GET "+versionV1+"/reviews", handler.GetAllReviews)
	mux.HandleFunc("GET "+versionV1+"/reviews/paginated", handler.GetAllReviewsPaginated)
	mux.HandleFunc("GET "+versionV1+"/reviews/id/{id}", handler.GetReviewByID)
	mux.HandleFunc("POST "+versionV1+"/reviews", handler.CreateReview)
	mux.HandleFunc("PUT "+versionV1+"/reviews/id/{id}", handler.UpdateReview)
	mux.HandleFunc("DELETE "+versionV1+"/reviews/id/{id}", handler.DeleteReview)
	mux.HandleFunc("GET "+versionV1+"/products/{id}/reviews", handler.GetReviewsByProduct)
	mux.HandleFunc("GET "+versionV1+"/users/{id}/reviews", handler.GetReviewsByUser)
	mux.HandleFunc("GET "+versionV1+"/products/{id}/rating-stats", handler.GetProductRatingStats)
}

func registerWishlistRoutes(mux *http.ServeMux) {
	handler := handler.NewWishlistHandler()

	mux.HandleFunc("GET "+versionV1+"/wishlists", handler.GetAllWishlists)
	mux.HandleFunc("GET "+versionV1+"/wishlists/paginated", handler.GetAllWishlistsPaginated)
	mux.HandleFunc("GET "+versionV1+"/wishlists/id/{id}", handler.GetWishlistByID)
	mux.HandleFunc("POST "+versionV1+"/wishlists", handler.CreateWishlist)
	mux.HandleFunc("DELETE "+versionV1+"/wishlists/id/{id}", handler.DeleteWishlist)
	mux.HandleFunc("DELETE "+versionV1+"/wishlists/clear", handler.ClearWishlist)
	mux.HandleFunc("GET "+versionV1+"/wishlists/count", handler.GetWishlistCount)
	mux.HandleFunc("GET "+versionV1+"/wishlists/product/{product_id}", handler.GetWishlistByProduct)
	mux.HandleFunc("DELETE "+versionV1+"/wishlists/product/{product_id}", handler.DeleteWishlistByProduct)
	mux.HandleFunc("GET "+versionV1+"/wishlists/product/{product_id}/check", handler.IsProductInWishlist)
}
