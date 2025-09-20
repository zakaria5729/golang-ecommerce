package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterReviewRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewReviewHandler()

	r.GET("/reviews/id/{id}/public", h.GetReviewByIdPublic).Register()

	r.GET("/reviews/paginated/public", h.GetAllReviewsPaginatedPublic).Register()

	r.GET("/products/{id}/reviews/public", h.GetReviewsByProductPublic).Register()

	r.GET("/products/{id}/rating-stats/public", h.GetProductRatingStatsPublic).Register()

	r.GET("/reviews/id/{id}", h.GetReviewByID).Register()

	r.GET("/reviews/paginated", h.GetAllReviewsPaginated).Register()

	r.GET("/products/{id}/reviews", h.GetReviewsByProduct).Register()

	r.GET("/products/{id}/rating-stats", h.GetProductRatingStats).Register()

	r.GET("/users/{id}/reviews", h.GetReviewsByUser).Use(
		pm.RequirePermission(c.PermissionReviewRead),
	).Register()

	r.POST("/reviews", h.CreateReview).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.PUT("/reviews/id/{id}", h.UpdateReview).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.DELETE("/reviews/id/{id}", h.DeleteReview).Use(
		pm.RequirePermission(c.PermissionReviewDelete),
	).Register()

	r.POST("/reviews/id/{id}/undo", h.UndoDeletedReview).Use(
		pm.RequirePermission(c.PermissionReviewUndoDelete),
	).Register()
}
