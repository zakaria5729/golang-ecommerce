package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterReviewRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
	handler := handler.NewReviewHandler()

	mux.HandleFunc(constants.GET+" /v1/products/{id}/rating-stats", handler.GetProductRatingStats)
	mux.HandleFunc(constants.GET+" /v1/reviews/paginated", handler.GetAllReviewsPaginated)
	mux.HandleFunc(constants.GET+" /v1/reviews/id/{id}", handler.GetReviewByID)

	mux.Handle(constants.GET+" /v1/users/{id}/reviews", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetReviewsByUser)))

	mux.Handle(constants.POST+" /v1/reviews", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.CreateReview)))

	mux.Handle(constants.PUT+" /v1/reviews/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.UpdateReview)))

	mux.Handle(constants.DELETE+" /v1/reviews/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.DeleteReview)))
}
