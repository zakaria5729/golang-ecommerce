package product_stats

import (
	"fmt"
	"time"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ProductStatsUseCase struct {
	repo *ProductStatsRepository
}

func NewProductStatsUseCase() *ProductStatsUseCase {
	return &ProductStatsUseCase{
		repo: NewProductStatsRepository(),
	}
}

// **REQUIRED
func (uc *ProductStatsUseCase) GetAllProductStatsPaginated(pageStr string, pageSizeStr string, productIDFilter string, dateFromFilter string, dateToFilter string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	productID, _ := utils.ParseUint(productIDFilter)
	dateFrom := uc.parseTimeFilter(dateFromFilter)
	dateTo := uc.parseTimeFilter(dateToFilter)

	histories, total, err := uc.repo.GetAllProductStatsPaginated(page, pageSize, productID, dateFrom, dateTo, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch browsing history paginated", "method", "GetAllBrowsingHistoryPaginated", "error", err, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch browsing history: %w", err)
	}

	return utils.BuildPaginatedResponse(histories, total, page, pageSize), nil
}

// **REQUIRED
func (uc *ProductStatsUseCase) GetProductStatsByID(id uint) (*ProductStats, error) {
	history, err := uc.repo.GetProductStatsByID(id)
	if err != nil {
		logger.Logger.Error("Browsing history not found", "method", "GetBrowsingHistoryByID", "error", err, "id", id)
		return nil, fmt.Errorf("browsing history not found: %w", err)
	}

	return history, nil
}

// **REQUIRED
func (uc *ProductStatsUseCase) IncreaseProductStats(req *IncreaseProductStatsRequest) error {
	productStats, err := uc.repo.GetProductStatsByID(req.ProductID)
	if err != nil || productStats == nil {
		productStats = &ProductStats{
			ProductID: req.ProductID,
		}
	}

	if req.IncreaseViewCount != nil && *req.IncreaseViewCount {
		productStats.ViewCount = productStats.ViewCount + 1
	}
	if req.IncreaseAddToCartCount != nil && *req.IncreaseAddToCartCount {
		productStats.AddToCartCount = productStats.AddToCartCount + 1
	}
	if req.IncreaseRemoveFromCartCount != nil && *req.IncreaseRemoveFromCartCount {
		productStats.RemoveFromCartCount = productStats.RemoveFromCartCount + 1
	}
	if req.IncreaseWishlistCount != nil && *req.IncreaseWishlistCount {
		productStats.WishlistCount = productStats.WishlistCount + 1
	}
	if req.IncreaseRemoveFromWishlistCount != nil && *req.IncreaseRemoveFromWishlistCount {
		productStats.RemoveFromWishlistCount = productStats.RemoveFromWishlistCount + 1
	}

	if err := uc.repo.IncreaseProductStats(productStats); err != nil {
		logger.Logger.Error("Failed to increase product stats", "method", "IncreaseProductStats", "error", err)
		return fmt.Errorf("failed to increase product stats: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *ProductStatsUseCase) IncreasePurchaseCountInProductStats(productID uint) error {
	productStats, err := uc.repo.GetProductStatsByID(productID)
	if err != nil || productStats == nil {
		productStats = &ProductStats{
			ProductID: productID,
		}
	}

	productStats.PurchaseCount = productStats.PurchaseCount + 1
	if err := uc.repo.IncreaseProductStats(productStats); err != nil {
		logger.Logger.Error("Failed to increase product stats", "method", "IncreasePurchaseCountInProductStats", "error", err)
		return fmt.Errorf("failed to increase product stats: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *ProductStatsUseCase) parseTimeFilter(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}

	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02 15:04:05",       // Standard format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t
		}
	}

	logger.Logger.Warn("Failed to parse time filter", "timeStr", timeStr)
	return nil
}

// func (uc *BrowsingHistoryUseCase) DeleteBrowsingHistory(id uint, userID uint) error {
// 	_, err := uc.repo.GetBrowsingHistoryByID(id, userID, nil)
// 	if err != nil {
// 		logger.Logger.Error("Browsing history not found", "method", "DeleteBrowsingHistory", "error", err, "id", id, "userID", userID)
// 		return fmt.Errorf("browsing history not found: %w", err)
// 	}

// 	if err := uc.repo.DeleteBrowsingHistory(id, userID); err != nil {
// 		logger.Logger.Error("Failed to delete browsing history", "method", "DeleteBrowsingHistory", "error", err, "id", id, "userID", userID)
// 		return fmt.Errorf("failed to delete browsing history: %w", err)
// 	}

// 	return nil
// }

// func (uc *BrowsingHistoryUseCase) ClearBrowsingHistory(userID uint) error {
// 	if err := uc.repo.ClearBrowsingHistory(userID); err != nil {
// 		logger.Logger.Error("Failed to clear browsing history", "method", "ClearBrowsingHistory", "error", err, "userID", userID)
// 		return fmt.Errorf("failed to clear browsing history: %w", err)
// 	}

// 	return nil
// }

// func (uc *BrowsingHistoryUseCase) GetRecentBrowsingHistory(userID uint, limit int) ([]BrowsingHistory, error) {
// 	if limit <= 0 {
// 		limit = 10 // Default limit
// 	}

// 	history, err := uc.repo.GetRecentBrowsingHistory(userID, limit)
// 	if err != nil {
// 		logger.Logger.Error("Failed to get recent browsing history", "method", "GetRecentBrowsingHistory", "error", err, "userID", userID, "limit", limit)
// 		return nil, fmt.Errorf("failed to get recent browsing history: %w", err)
// 	}

// 	return history, nil
// }

// func (uc *BrowsingHistoryUseCase) GetMostViewedProducts(userID uint, limit int) ([]uint, error) {
// 	if limit <= 0 {
// 		limit = 10 // Default limit
// 	}

// 	productIDs, err := uc.repo.GetMostViewedProducts(userID, limit)
// 	if err != nil {
// 		logger.Logger.Error("Failed to get most viewed products", "method", "GetMostViewedProducts", "error", err, "userID", userID, "limit", limit)
// 		return nil, fmt.Errorf("failed to get most viewed products: %w", err)
// 	}

// 	return productIDs, nil
// }

// func (uc *BrowsingHistoryUseCase) GetAllBrowsingHistory(userID uint, includeStr string, productIDFilter string, dateFromFilter string, dateToFilter string, sortBy, sortOrder string) ([]BrowsingHistory, error) {
// 	include := utils.ParseCommaSeparatedString(includeStr)
// 	productID, _ := utils.ParseUint(productIDFilter)
// 	dateFrom := uc.parseTimeFilter(dateFromFilter)
// 	dateTo := uc.parseTimeFilter(dateToFilter)

// 	history, err := uc.repo.GetAllBrowsingHistory(userID, include, productID, dateFrom, dateTo, sortBy, sortOrder)
// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch browsing history", "method", "GetAllBrowsingHistory", "error", err, "userID", userID, "include", include, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
// 		return nil, fmt.Errorf("failed to fetch browsing history: %w", err)
// 	}

// 	return history, nil
// }

// func (uc *BrowsingHistoryUseCase) validateCreateRequest(req *BrowsingHistory) validator.ValidationErrors {
// 	var errors validator.ValidationErrors

// 	errors = validator.MergeValidationErrors(
// 		errors,
// 		validator.ValidateRequired(fmt.Sprintf("%d", req.ProductID), "product_id"),
// 	)

// 	if req.ProductID == 0 {
// 		errors.AddError("product_id", "product_id must be a positive integer")
// 	}

// 	return errors
// }
