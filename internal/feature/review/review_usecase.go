package review

import (
	"context"
	"fmt"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ReviewUseCase struct {
	repo *ReviewRepository
}

func NewReviewUseCase() *ReviewUseCase {
	return &ReviewUseCase{
		repo: NewReviewRepository(),
	}
}

func (uc *ReviewUseCase) GetAllReviewsPaginated(showDeleted *bool, productIDFilter string, userIDFilter string, ratingFromFilter string, ratingToFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	productID, _ := utils.ParseUint(productIDFilter)
	userID, _ := utils.ParseUint(userIDFilter)
	ratingFrom, _ := utils.ParseInt(ratingFromFilter)
	ratingTo, _ := utils.ParseInt(ratingToFilter)

	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	reviews, total, err := uc.repo.GetAllReviewsPaginated(showDeleted, productID, userID, ratingFrom, ratingTo, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews paginated", "method", "GetAllReviewsPaginated", "error", err, "productID", productID, "userID", userID, "ratingFrom", ratingFrom, "ratingTo", ratingTo, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (uc *ReviewUseCase) GetReviewByID(id uint, showDeleted *bool) (*Review, error) {
	review, err := uc.repo.GetReviewByID(id, showDeleted)

	if err != nil {
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id)
		return nil, fmt.Errorf("failed to fetch review: %w", err)
	}

	return review, nil
}

func (uc *ReviewUseCase) CreateReview(userID uint, productID uint, rating int, comment string) (*Review, error) {
	if rating < c.MinReviewRating || rating > c.MaxReviewRating {
		return nil, fmt.Errorf("invalid rating")
	}

	exists, _ := uc.repo.CheckUserReviewExists(productID, userID)
	if exists {
		return nil, fmt.Errorf("user has already reviewed this product")
	}

	comment = utils.Trim(comment)
	review := &Review{
		ProductID: productID,
		UserID:    userID,
		Rating:    rating,
		Comment:   &comment,
	}
	review.CreatedBy = &userID

	err := uc.repo.CreateReview(review)
	if err != nil {
		logger.Logger.Error("Failed to create review", "method", "CreateReview", "error", err, "review", review)
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (uc *ReviewUseCase) UpdateReview(ctx context.Context, id uint, userID *uint, rating int, comment string) error {
	review := &Review{}

	if rating != 0 {
		if rating < c.MinReviewRating || rating > c.MaxReviewRating {
			return fmt.Errorf("rating must be between %d and %d", c.MinReviewRating, c.MaxReviewRating)
		}

		if rating < c.MinReviewRating || rating > c.MaxReviewRating {
			return fmt.Errorf("rating must be between %d and %d", c.MinReviewRating, c.MaxReviewRating)
		}

		review.Rating = rating
	}

	if comment != "" {
		review.Comment = &comment
	}

	if review.Rating == 0 && review.Comment == nil {
		return fmt.Errorf("no fields to update")
	}

	review.UpdatedBy = m.GetUserIdOnlyFromContext(ctx)
	err := uc.repo.UpdateReview(id, userID, review)
	if err != nil {
		logger.Logger.Error("Failed to update review", "method", "UpdateReview", "error", err, "id", id, "userID", userID, "review", review)
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

func (uc *ReviewUseCase) UndoDeletedReview(ctx context.Context, idStr string) error {
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil || *id == 0 {
		return fmt.Errorf("invalid review ID: %w", err)
	}

	err = uc.repo.UndoDeletedReview(ctx, *id)
	if err != nil {
		logger.Logger.Error("Failed to undo deleted review", "method", "UndoDeletedReview", "error", err, "id", *id)
		return fmt.Errorf("failed to undo deleted review: %w", err)
	}

	return nil
}

func (uc *ReviewUseCase) DeleteReview(ctx context.Context, idStr string, userID *uint) error {
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil || *id == 0 {
		return fmt.Errorf("invalid review ID: %w", err)
	}

	err = uc.repo.DeleteReview(ctx, *id, userID)
	if err != nil {
		logger.Logger.Error("Failed to delete review", "method", "DeleteReview", "error", err, "id", *id, "userID", userID)
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

func (uc *ReviewUseCase) GetReviewsByProduct(productID uint, showDeleted *bool, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	rating, _ := utils.ParseInt(ratingFilter)
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	reviews, total, err := uc.repo.GetReviewsByProduct(productID, showDeleted, rating, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by product", "method", "GetReviewsByProduct", "error", err, "productID", productID, "rating", rating, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (uc *ReviewUseCase) GetReviewsByUser(userID uint, showDeleted *bool, productIDStr string, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	rating, _ := utils.ParseInt(ratingFilter)
	productID, _ := utils.ParseUint(productIDStr)
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	var ratingPtr *int
	if rating != nil {
		ratingPtr = rating
	}

	reviews, total, err := uc.repo.GetReviewsByUser(userID, showDeleted, productID, ratingPtr, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by user", "method", "GetReviewsByUser", "error", err, "userID", userID, "productID", productID, "rating", ratingPtr, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (uc *ReviewUseCase) GetProductRatingStats(productID uint, showDeleted *bool) (*ProductRatingStatsResponse, error) {
	avgRating, err := uc.repo.GetAverageRating(productID, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to get average rating", "method", "GetProductRatingStats", "error", err, "productID", productID)
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}

	ratingCounts, err := uc.repo.GetRatingCounts(productID, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to get rating counts", "method", "GetProductRatingStats", "error", err, "productID", productID)
		return nil, fmt.Errorf("failed to get rating counts: %w", err)
	}

	return &ProductRatingStatsResponse{
		AverageRating: avgRating,
		RatingCounts:  ratingCounts,
	}, nil
}
