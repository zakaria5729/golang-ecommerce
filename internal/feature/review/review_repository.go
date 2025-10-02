package review

import (
	"fmt"

	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/constants"
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

func (r *ReviewRepository) GetAllReviewsPaginated(showDeleted *bool, productID *uint, userID *uint, ratingFrom *int, ratingTo *int, page, pageSize int, sortBy, sortOrder string) ([]Review, int64, error) {
	var reviews []Review
	var total int64
	query := r.db.Model(&Review{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if productID != nil {
		query = query.Where(c.ReviewProductID+" = ?", *productID)
	}

	if userID != nil {
		query = query.Where(c.ReviewUserID+" = ?", *userID)
	}

	if ratingFrom != nil {
		query = query.Where(c.ReviewRating+" >= ?", *ratingFrom)
	}

	if ratingTo != nil {
		query = query.Where(c.ReviewRating+" <= ?", *ratingTo)
	}

	err := query.Count(&total).Error
	if err != nil {
		logger.Logger.Error("Failed to count reviews", "method", "GetAllReviewsPaginated", "error", err, "productID", productID, "userID", userID, "ratingFrom", ratingFrom, "ratingTo", ratingTo)
		return nil, 0, err
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err = query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&reviews).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews paginated", "method", "GetAllReviewsPaginated", "error", err, "productID", productID, "userID", userID, "ratingFrom", ratingFrom, "ratingTo", ratingTo, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) GetReviewByID(id uint, showDeleted *bool) (*Review, error) {
	var review Review
	query := r.db.Model(&Review{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Where(c.FieldID+" = ?", id).Take(&review).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id)
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

func (r *ReviewRepository) UpdateReview(id uint, userID *uint, review *Review) error {
	query := r.db.Model(&Review{})

	if userID != nil {
		query = query.Where(c.FieldID+" = ? AND "+c.ReviewUserID+" = ?", id, userID)
	} else {
		query = query.Where(c.FieldID + " = ?")
	}

	result := query.Updates(review)
	if result.Error != nil {
		logger.Logger.Error("Failed to update review", "method", "UpdateReview", "error", result.Error, "id", id, "userID", userID, "review", review)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("review not found or not owned by user")
	}

	return nil
}

func (r *ReviewRepository) UndoDeletedReview(id uint) error {
	result := r.db.Unscoped().Model(&Review{}).Select(c.FieldCreatedAt).Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil)
	if result.Error != nil {
		logger.Logger.Error("Failed to undo deleted review", "method", "UndoDeletedReview", "error", result.Error, "id", id)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("review not found")
	}

	return nil
}

func (r *ReviewRepository) DeleteReview(id uint, userID *uint) error {
	query := r.db

	if userID != nil {
		query = query.Where(c.FieldID+" = ? AND "+c.ReviewUserID+" = ?", id, userID)
	} else {
		query = query.Where(c.FieldID+" = ?", id)
	}

	result := query.Delete(&Review{})
	if result.Error != nil {
		logger.Logger.Error("Failed to delete review", "method", "DeleteReview", "error", result.Error, "id", id, "userID", userID)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("review not found or not owned by user")
	}

	return nil
}

func (r *ReviewRepository) GetReviewsByProduct(productID uint, showDeleted *bool, rating *int, page, pageSize int, sortBy, sortOrder string) ([]Review, int64, error) {
	var reviews []Review
	var total int64
	query := r.db.Model(&Review{}).Where(c.ReviewProductID+" = ?", productID)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if rating != nil {
		query = query.Where(c.ReviewRating+" = ?", *rating)
	}

	err := query.Count(&total).Error
	if err != nil {
		logger.Logger.Error("Failed to count reviews", "method", "GetReviewsByProduct", "error", err, "productID", productID, "rating", rating)
		return nil, 0, err
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err = query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&reviews).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by product paginated", "method", "GetReviewsByProduct", "error", err, "productID", productID, "rating", rating, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) GetReviewsByUser(userID uint, showDeleted *bool, productID *uint, rating *int, page, pageSize int, sortBy, sortOrder string) ([]Review, int64, error) {
	var reviews []Review
	var total int64
	query := r.db.Model(&Review{}).Where(c.ReviewUserID+" = ?", userID)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if productID != nil {
		query = query.Where(c.ReviewProductID+" = ?", *productID)
	}

	if rating != nil {
		query = query.Where(c.ReviewRating+" = ?", *rating)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Count(&total).Error
	if err != nil {
		logger.Logger.Error("Failed to count reviews", "method", "GetReviewsByUser", "error", err, "productID", productID, "userID", userID, "rating", rating)
		return nil, 0, err
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err = query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&reviews).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch reviews by user paginated", "method", "GetReviewsByUser", "error", err, "productID", productID, "userID", userID, "rating", rating, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) GetAverageRating(productID uint, showDeleted *bool) (float64, error) {
	var avgRating float64
	query := r.db.Model(&Review{}).Where(c.ReviewProductID+" = ?", productID)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select("AVG(rating)").Scan(&avgRating).Error
	if err != nil {
		logger.Logger.Error("Failed to get average rating", "method", "GetAverageRating", "error", err, "productID", productID)
		return 0, err
	}

	return avgRating, nil
}

func (r *ReviewRepository) GetRatingCounts(productID uint, showDeleted *bool) (map[int]int, error) {
	var results []struct {
		Rating int
		Count  int
	}

	query := r.db.Model(&Review{})
	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.
		Select(c.ReviewRating+", COUNT(*) as count").
		Where(c.ReviewProductID+" = ?", productID).
		Group(c.ReviewRating).
		Order(c.ReviewRating).
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
	var review Review
	err := r.db.Model(&Review{}).
		Where(c.ReviewProductID+" = ? AND "+c.ReviewUserID+" = ?", productID, userID).
		Select(c.FieldID).
		Take(&review).Error

	if err != nil {
		logger.Logger.Error("Failed to check user review exists", "method", "CheckUserReviewExists", "error", err, "productID", productID, "userID", userID)
		return false, err
	}

	return review.ID != 0, nil
}
