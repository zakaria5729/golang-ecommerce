package review

import (
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ReviewUseCase struct {
	repo *ReviewRepository
}

func NewReviewUseCase() *ReviewUseCase {
	return &ReviewUseCase{
		repo: NewReviewRepository(),
	}
}

func (uc *ReviewUseCase) GetAllReviews(includeStr string, productIDFilter string, userIDFilter string, ratingFilter string, sortBy, sortOrder string) ([]Review, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	productID, _ := utils.ParseUint(productIDFilter)
	userID, _ := utils.ParseUint(userIDFilter)
	rating, _ := utils.ParseInt(ratingFilter)

	reviews, err := uc.repo.GetAllReviews(include, productID, userID, rating, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews", "method", "GetAllReviews", "error", err, "include", include, "productID", productID, "userID", userID, "rating", rating, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return reviews, nil
}

func (uc *ReviewUseCase) GetAllReviewsPaginated(includeStr string, productIDFilter string, userIDFilter string, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	include := utils.ParseCommaSeparatedString(includeStr)
	productID, _ := utils.ParseUint(productIDFilter)
	userID, _ := utils.ParseUint(userIDFilter)
	rating, _ := utils.ParseInt(ratingFilter)

	reviews, total, err := uc.repo.GetAllReviewsPaginated(include, productID, userID, rating, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews paginated", "method", "GetAllReviewsPaginated", "error", err, "include", include, "productID", productID, "userID", userID, "rating", rating, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	var reviewPtrs []*Review
	for i := range reviews {
		reviewPtrs = append(reviewPtrs, &reviews[i])
	}

	return utils.BuildPaginatedResponse(reviewPtrs, int(total), page, pageSize), nil
}

func (uc *ReviewUseCase) GetReviewByID(idStr string, includeStr string) (*Review, error) {
	id, err := utils.ParseUint(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid review ID: %w", err)
	}

	include := utils.ParseCommaSeparatedString(includeStr)
	review, err := uc.repo.GetReviewByID(*id, include)

	if err != nil {
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("failed to fetch review: %w", err)
	}

	return review, nil
}

func (uc *ReviewUseCase) CreateReview(userID uint, productIDStr string, ratingStr string, comment string) (*Review, error) {
	productID, err := utils.ParseUint(productIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	rating, err := utils.ParseInt(ratingStr)
	if err != nil {
		return nil, fmt.Errorf("invalid rating: %w", err)
	}

	if *rating < MinRating || *rating > MaxRating {
		return nil, fmt.Errorf("rating must be between %d and %d", MinRating, MaxRating)
	}

	exists, err := uc.repo.CheckUserReviewExists(*productID, userID)
	if err != nil {
		logger.Logger.Error("Failed to check user review exists", "method", "CreateReview", "error", err, "productID", *productID, "userID", userID)
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user has already reviewed this product")
	}

	review := &Review{
		ProductID: *productID,
		UserID:    userID,
		Rating:    *rating,
		Comment:   utils.ParseStringPtr(comment),
	}

	err = uc.repo.CreateReview(review)
	if err != nil {
		logger.Logger.Error("Failed to create review", "method", "CreateReview", "error", err, "review", review)
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (uc *ReviewUseCase) UpdateReview(idStr string, userID uint, ratingStr string, comment string) (*Review, error) {
	id, err := utils.ParseUint(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid review ID: %w", err)
	}

	updates := make(map[string]any)

	if ratingStr != "" {
		rating, err := utils.ParseInt(ratingStr)
		if err != nil {
			return nil, fmt.Errorf("invalid rating: %w", err)
		}

		if *rating < MinRating || *rating > MaxRating {
			return nil, fmt.Errorf("rating must be between %d and %d", MinRating, MaxRating)
		}

		updates[ReviewRating] = *rating
	}

	if comment != "" {
		updates[ReviewComment] = utils.Trim(comment)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	err = uc.repo.UpdateReview(*id, userID, updates)
	if err != nil {
		logger.Logger.Error("Failed to update review", "method", "UpdateReview", "error", err, "id", *id, "userID", userID, "updates", updates)
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	review, err := uc.repo.GetReviewByID(*id, []string{})
	if err != nil {
		logger.Logger.Error("Failed to fetch updated review", "method", "UpdateReview", "error", err, "id", *id)
		return nil, fmt.Errorf("failed to fetch updated review: %w", err)
	}

	return review, nil
}

func (uc *ReviewUseCase) DeleteReview(idStr string, userID uint) error {
	id, err := utils.ParseUint(idStr)
	if err != nil {
		return fmt.Errorf("invalid review ID: %w", err)
	}

	err = uc.repo.DeleteReview(*id, userID)
	if err != nil {
		logger.Logger.Error("Failed to delete review", "method", "DeleteReview", "error", err, "id", *id, "userID", userID)
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

func (uc *ReviewUseCase) GetReviewsByProduct(productIDStr string, includeStr string, ratingFilter string, sortBy, sortOrder string) ([]Review, error) {
	productID, err := utils.ParseUint(productIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	include := utils.ParseCommaSeparatedString(includeStr)
	rating, _ := utils.ParseInt(ratingFilter)

	var ratingPtr *int
	if rating != nil {
		ratingPtr = rating
	}

	reviews, err := uc.repo.GetReviewsByProduct(*productID, include, ratingPtr, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by product", "method", "GetReviewsByProduct", "error", err, "productID", *productID, "include", include, "rating", ratingPtr, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return reviews, nil
}

func (uc *ReviewUseCase) GetReviewsByUser(userIDStr string, includeStr string, ratingFilter string, sortBy, sortOrder string) ([]Review, error) {
	userID, err := utils.ParseUint(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	include := utils.ParseCommaSeparatedString(includeStr)
	rating, _ := utils.ParseInt(ratingFilter)

	var ratingPtr *int
	if rating != nil {
		ratingPtr = rating
	}

	reviews, err := uc.repo.GetReviewsByUser(*userID, include, ratingPtr, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by user", "method", "GetReviewsByUser", "error", err, "userID", *userID, "include", include, "rating", ratingPtr, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return reviews, nil
}

func (uc *ReviewUseCase) GetProductRatingStats(productIDStr string) (map[string]interface{}, error) {
	productID, err := utils.ParseUint(productIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	avgRating, err := uc.repo.GetAverageRating(*productID)
	if err != nil {
		logger.Logger.Error("Failed to get average rating", "method", "GetProductRatingStats", "error", err, "productID", *productID)
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}

	ratingCounts, err := uc.repo.GetRatingCounts(*productID)
	if err != nil {
		logger.Logger.Error("Failed to get rating counts", "method", "GetProductRatingStats", "error", err, "productID", *productID)
		return nil, fmt.Errorf("failed to get rating counts: %w", err)
	}

	return map[string]interface{}{
		"average_rating": avgRating,
		"rating_counts":  ratingCounts,
	}, nil
}

func (uc *ReviewUseCase) ValidateReviewInput(ratingStr string, comment string) []validator.ValidationError {
	var errors []validator.ValidationError

	if ratingStr == "" {
		errors = append(errors, validator.ValidationError{
			Field:   "rating",
			Message: "rating is required",
		})
	} else {
		rating, err := utils.ParseInt(ratingStr)
		if err != nil {
			errors = append(errors, validator.ValidationError{
				Field:   "rating",
				Message: "rating must be a valid number",
			})
		} else if *rating < MinRating || *rating > MaxRating {
			errors = append(errors, validator.ValidationError{
				Field:   "rating",
				Message: fmt.Sprintf("rating must be between %d and %d", MinRating, MaxRating),
			})
		}
	}

	if comment != "" && len(comment) > 1000 {
		errors = append(errors, validator.ValidationError{
			Field:   "comment",
			Message: "comment must be less than 1000 characters",
		})
	}

	return errors
}
