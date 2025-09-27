package review

import (
	"fmt"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
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

func (uc *ReviewUseCase) GetAllReviewsPaginated(showDeleted *bool, includeStr string, productIDFilter string, userIDFilter string, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	productID, _ := utils.ParseUint(productIDFilter)
	userID, _ := utils.ParseUint(userIDFilter)
	rating, _ := utils.ParseInt(ratingFilter)

	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)

	reviews, total, err := uc.repo.GetAllReviewsPaginated(showDeleted, include, productID, userID, rating, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews paginated", "method", "GetAllReviewsPaginated", "error", err, "include", include, "productID", productID, "userID", userID, "rating", rating, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (uc *ReviewUseCase) GetReviewByID(id uint, showDeleted *bool, includeStr string) (*Review, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	review, err := uc.repo.GetReviewByID(id, showDeleted, include)

	if err != nil {
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("failed to fetch review: %w", err)
	}

	return review, nil
}

func (uc *ReviewUseCase) CreateReview(userID uint, productIDStr string, ratingStr string, comment string) (*Review, error) {
	productID, err := utils.ParseUint(productIDStr)
	if err != nil || productID == nil || *productID == 0 {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	rating, err := utils.ParseInt(ratingStr)
	if err != nil || rating == nil || *rating < constants.MinReviewRating || *rating > constants.MaxReviewRating {
		return nil, fmt.Errorf("invalid rating: %w", err)
	}

	exists, err := uc.repo.CheckUserReviewExists(*productID, userID)
	if err != nil {
		logger.Logger.Error("Failed to check user review exists", "method", "CreateReview", "error", err, "productID", *productID, "userID", userID)
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user has already reviewed this product")
	}

	comment = utils.Trim(comment)
	review := &Review{
		ProductID: *productID,
		UserID:    userID,
		Rating:    *rating,
		Comment:   &comment,
	}

	err = uc.repo.CreateReview(review)
	if err != nil {
		logger.Logger.Error("Failed to create review", "method", "CreateReview", "error", err, "review", review)
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (uc *ReviewUseCase) UpdateReview(id uint, userID uint, ratingStr string, comment string) error {
	review := &Review{}

	if ratingStr != "" {
		rating, err := utils.ParseInt(ratingStr)
		if err != nil {
			return fmt.Errorf("invalid rating: %w", err)
		}

		if *rating < constants.MinReviewRating || *rating > constants.MaxReviewRating {
			return fmt.Errorf("rating must be between %d and %d", constants.MinReviewRating, constants.MaxReviewRating)
		}

		review.Rating = *rating
	}

	if comment != "" {
		review.Comment = &comment
	}

	if review.Rating == 0 && review.Comment == nil {
		return fmt.Errorf("no fields to update")
	}

	err := uc.repo.UpdateReview(id, userID, review)
	if err != nil {
		logger.Logger.Error("Failed to update review", "method", "UpdateReview", "error", err, "id", id, "userID", userID, "review", review)
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

func (uc *ReviewUseCase) UndoDeletedReview(idStr string, userID uint) error {
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil || *id == 0 {
		return fmt.Errorf("invalid review ID: %w", err)
	}

	err = uc.repo.UndoDeletedReview(*id, userID)
	if err != nil {
		logger.Logger.Error("Failed to undo deleted review", "method", "UndoDeletedReview", "error", err, "id", *id, "userID", userID)
		return fmt.Errorf("failed to undo deleted review: %w", err)
	}

	return nil
}

func (uc *ReviewUseCase) DeleteReview(idStr string, userID uint) error {
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil || *id == 0 {
		return fmt.Errorf("invalid review ID: %w", err)
	}

	err = uc.repo.DeleteReview(*id, userID)
	if err != nil {
		logger.Logger.Error("Failed to delete review", "method", "DeleteReview", "error", err, "id", *id, "userID", userID)
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

func (uc *ReviewUseCase) GetReviewsByProduct(productID uint, showDeleted *bool, includeStr string, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	rating, _ := utils.ParseInt(ratingFilter)
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	reviews, total, err := uc.repo.GetReviewsByProduct(productID, showDeleted, include, rating, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by product", "method", "GetReviewsByProduct", "error", err, "productID", productID, "include", include, "rating", rating, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (uc *ReviewUseCase) GetReviewsByUser(userID uint, showDeleted *bool, productIDStr string, includeStr string, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	rating, _ := utils.ParseInt(ratingFilter)
	productID, _ := utils.ParseUint(productIDStr)
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	var ratingPtr *int
	if rating != nil {
		ratingPtr = rating
	}

	reviews, total, err := uc.repo.GetReviewsByUser(userID, showDeleted, productID, include, ratingPtr, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by user", "method", "GetReviewsByUser", "error", err, "userID", userID, "productID", productID, "include", include, "rating", ratingPtr, "sortBy", sortBy, "sortOrder", sortOrder)
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
