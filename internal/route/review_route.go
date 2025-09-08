package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterReviewRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	handler := handler.NewReviewHandler()

	// Public endpoints - anyone can read reviews
	mux.HandleFunc("GET /v1/products/{id}/reviews", handler.GetReviewsByProduct)
	mux.HandleFunc("GET /v1/products/{id}/rating-stats", handler.GetProductRatingStats)

	// Protected endpoints - require authentication
	mux.Handle("GET /v1/reviews", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAllReviews)))

	mux.Handle("GET /v1/reviews/paginated", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAllReviewsPaginated)))

	mux.Handle("GET /v1/reviews/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetReviewByID)))

	mux.Handle("GET /v1/users/{id}/reviews", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetReviewsByUser)))

	// Write operations - require authentication
	mux.Handle("POST /v1/reviews", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.CreateReview)))

	mux.Handle("PUT /v1/reviews/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.UpdateReview)))

	mux.Handle("DELETE /v1/reviews/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.DeleteReview)))
}
