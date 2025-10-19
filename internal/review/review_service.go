package review

import (
	"context"
	"fmt"

	m "github.com/easy-comerce/backend/internal/review/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	r "github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ReviewService struct {
	repo *ReviewRepository
}

func NewReviewService(repo *ReviewRepository) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

func (s *ReviewService) GetAllReviewsPaginated(showDeleted *bool, productIDFilter string, userIDFilter string, ratingFromFilter string, ratingToFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	productID, _ := utils.ParseUint(productIDFilter)
	userID, _ := utils.ParseUint(userIDFilter)
	ratingFrom, _ := utils.ParseInt(ratingFromFilter)
	ratingTo, _ := utils.ParseInt(ratingToFilter)

	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	reviews, total, err := s.repo.GetAllReviewsPaginated(showDeleted, productID, userID, ratingFrom, ratingTo, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (s *ReviewService) GetReviewByID(id uint, showDeleted *bool) (*ReviewEntity, error) {
	review, err := s.repo.GetReviewByID(id, showDeleted)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch review: %w", err)
	}

	return review, nil
}

func (s *ReviewService) CreateReview(userID uint, productID uint, rating int, comment string) (*ReviewEntity, error) {
	if rating < c.MinReviewRating || rating > c.MaxReviewRating {
		return nil, fmt.Errorf("invalid rating")
	}

	exists, _ := s.repo.CheckUserReviewExists(productID, userID)
	if exists {
		return nil, fmt.Errorf("user has already reviewed this product")
	}

	comment = utils.Trim(comment)
	review := &ReviewEntity{
		ProductID: productID,
		UserID:    userID,
		Rating:    rating,
		Comment:   &comment,
	}
	review.CreatedBy = &userID

	err := s.repo.CreateReview(review)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, id uint, userID *uint, rating int, comment string) error {
	review := &ReviewEntity{}

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

	review.UpdatedBy, _ = cu.GetUserIDFromContext(ctx)
	err := s.repo.UpdateReview(id, userID, review)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

func (s *ReviewService) UndoDeletedReview(ctx context.Context, idStr string) error {
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil || *id == 0 {
		return fmt.Errorf("invalid review ID: %w", err)
	}

	err = s.repo.UndoDeletedReview(ctx, *id)
	if err != nil {
		return fmt.Errorf("failed to undo deleted review: %w", err)
	}

	return nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, idStr string, userID *uint) error {
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil || *id == 0 {
		return fmt.Errorf("invalid review ID: %w", err)
	}

	err = s.repo.DeleteReview(ctx, *id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

func (s *ReviewService) GetReviewsByProduct(productID uint, showDeleted *bool, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	rating, _ := utils.ParseInt(ratingFilter)
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	reviews, total, err := s.repo.GetReviewsByProduct(productID, showDeleted, rating, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (s *ReviewService) GetReviewsByUser(userID uint, showDeleted *bool, productIDStr string, ratingFilter string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	rating, _ := utils.ParseInt(ratingFilter)
	productID, _ := utils.ParseUint(productIDStr)
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	var ratingPtr *int
	if rating != nil {
		ratingPtr = rating
	}

	reviews, total, err := s.repo.GetReviewsByUser(userID, showDeleted, productID, ratingPtr, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return utils.BuildPaginatedResponse(reviews, int(total), page, pageSize), nil
}

func (s *ReviewService) GetProductRatingStats(productID uint, showDeleted *bool) (*m.ProductRatingStatsResponse, error) {
	avgRating, err := s.repo.GetAverageRating(productID, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}

	ratingCounts, err := s.repo.GetRatingCounts(productID, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating counts: %w", err)
	}

	return &m.ProductRatingStatsResponse{
		AverageRating: avgRating,
		RatingCounts:  ratingCounts,
	}, nil
}
