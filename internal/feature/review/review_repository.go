package review

import (
	"errors"
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		db: db.GetDB(),
	}
}

func (r *ReviewRepository) GetAllReviews(include []string, productID *uint, userID *uint, rating *int, sortBy, sortOrder string) ([]Review, error) {
	var reviews []Review

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if productID != nil {
		query = query.Where(ReviewProductID+" = ?", *productID)
	}

	if userID != nil {
		query = query.Where(ReviewUserID+" = ?", *userID)
	}

	if rating != nil {
		query = query.Where(ReviewRating+" = ?", *rating)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&reviews).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews", "method", "GetAllReviews", "error", err, "include", include, "productID", productID, "userID", userID, "rating", rating, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return reviews, err
}

func (r *ReviewRepository) GetAllReviewsPaginated(include []string, productID *uint, userID *uint, rating *int, page, pageSize int, sortBy, sortOrder string) ([]Review, int64, error) {
	var reviews []Review
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if productID != nil {
		query = query.Where(ReviewProductID+" = ?", *productID)
	}

	if userID != nil {
		query = query.Where(ReviewUserID+" = ?", *userID)
	}

	if rating != nil {
		query = query.Where(ReviewRating+" = ?", *rating)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&reviews).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews paginated", "method", "GetAllReviewsPaginated", "error", err, "include", include, "productID", productID, "userID", userID, "rating", rating, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	countQuery := r.db.Model(&Review{})
	if productID != nil {
		countQuery = countQuery.Where(ReviewProductID+" = ?", *productID)
	}
	if userID != nil {
		countQuery = countQuery.Where(ReviewUserID+" = ?", *userID)
	}
	if rating != nil {
		countQuery = countQuery.Where(ReviewRating+" = ?", *rating)
	}

	err = countQuery.Count(&total).Error
	if err != nil {
		logger.Logger.Error("Failed to count reviews", "method", "GetAllReviewsPaginated", "error", err, "productID", productID, "userID", userID, "rating", rating)
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) GetReviewByID(id uint, include []string) (*Review, error) {
	var review Review

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	err := query.Where(constants.FieldID+" = ?", id).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("review not found")
		}
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &review, nil
}

func (r *ReviewRepository) CreateReview(review *Review) error {
	review.Sanitize()

	err := r.db.Create(review).Error
	if err != nil {
		logger.Logger.Error("Failed to create review", "method", "CreateReview", "error", err, "review", review)
		return err
	}

	return nil
}

func (r *ReviewRepository) UpdateReview(id uint, userID uint, updates map[string]interface{}) error {
	query := r.db.Model(&Review{}).Where(constants.FieldID+" = ? AND "+ReviewUserID+" = ?", id, userID)

	result := query.Updates(updates)
	if result.Error != nil {
		logger.Logger.Error("Failed to update review", "method", "UpdateReview", "error", result.Error, "id", id, "userID", userID, "updates", updates)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("review not found or not owned by user")
	}

	return nil
}

func (r *ReviewRepository) DeleteReview(id uint, userID uint) error {
	result := r.db.Where(constants.FieldID+" = ? AND "+ReviewUserID+" = ?", id, userID).Delete(&Review{})
	if result.Error != nil {
		logger.Logger.Error("Failed to delete review", "method", "DeleteReview", "error", result.Error, "id", id, "userID", userID)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("review not found or not owned by user")
	}

	return nil
}

func (r *ReviewRepository) GetReviewsByProduct(productID uint, include []string, rating *int, sortBy, sortOrder string) ([]Review, error) {
	var reviews []Review

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", ")).Where(ReviewProductID+" = ?", productID)

	if rating != nil {
		query = query.Where(ReviewRating+" = ?", *rating)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&reviews).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by product", "method", "GetReviewsByProduct", "error", err, "productID", productID, "include", include, "rating", rating, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return reviews, err
}

func (r *ReviewRepository) GetReviewsByUser(userID uint, include []string, rating *int, sortBy, sortOrder string) ([]Review, error) {
	var reviews []Review

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", ")).Where(ReviewUserID+" = ?", userID)

	if rating != nil {
		query = query.Where(ReviewRating+" = ?", *rating)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&reviews).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by user", "method", "GetReviewsByUser", "error", err, "userID", userID, "include", include, "rating", rating, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return reviews, err
}

func (r *ReviewRepository) GetAverageRating(productID uint) (float64, error) {
	var avgRating float64

	err := r.db.Model(&Review{}).Where(ReviewProductID+" = ?", productID).Select("AVG(rating)").Scan(&avgRating).Error
	if err != nil {
		logger.Logger.Error("Failed to get average rating", "method", "GetAverageRating", "error", err, "productID", productID)
		return 0, err
	}

	return avgRating, nil
}

func (r *ReviewRepository) GetRatingCounts(productID uint) (map[int]int, error) {
	var results []struct {
		Rating int
		Count  int
	}

	err := r.db.Model(&Review{}).
		Select("rating, COUNT(*) as count").
		Where(ReviewProductID+" = ?", productID).
		Group("rating").
		Order("rating").
		Scan(&results).Error

	if err != nil {
		logger.Logger.Error("Failed to get rating counts", "method", "GetRatingCounts", "error", err, "productID", productID)
		return nil, err
	}

	ratingCounts := make(map[int]int)
	for _, result := range results {
		ratingCounts[result.Rating] = result.Count
	}

	return ratingCounts, nil
}

func (r *ReviewRepository) CheckUserReviewExists(productID, userID uint) (bool, error) {
	var count int64

	err := r.db.Model(&Review{}).Where(ReviewProductID+" = ? AND "+ReviewUserID+" = ?", productID, userID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check user review exists", "method", "CheckUserReviewExists", "error", err, "productID", productID, "userID", userID)
		return false, err
	}

	return count > 0, nil
}

func (r *ReviewRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, ReviewProductID, ReviewUserID, ReviewRating, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{ReviewComment}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
