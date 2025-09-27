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

func (uc *ProductStatsUseCase) GetProductStatsByID(id uint) (*ProductStats, error) {
	history, err := uc.repo.GetProductStatsByID(id)
	if err != nil {
		logger.Logger.Error("Browsing history not found", "method", "GetBrowsingHistoryByID", "error", err, "id", id)
		return nil, fmt.Errorf("browsing history not found: %w", err)
	}

	return history, nil
}

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

func (uc *ProductStatsUseCase) parseTimeFilter(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}

	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02 15:04:05",       // Standard format
		"2006-01-02",                // Date format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t
		}
	}

	logger.Logger.Warn("Failed to parse time filter", "timeStr", timeStr)
	return nil
}
