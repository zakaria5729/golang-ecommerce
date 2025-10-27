package product_stats

import (
	"context"
	"fmt"
	"time"

	m "github.com/easy-comerce/backend/internal/product_stats/model"
	l "github.com/easy-comerce/backend/pkg/logger"
	r "github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ProductStatsService interface {
	GetAllProductStatsPaginated(pageStr string, pageSizeStr string, productIDFilter string, dateFromFilter string, dateToFilter string, sortBy, sortOrder string) (*r.PaginatedResponse, error)
	GetProductStatsByID(id uint) (*ProductStatsEntity, error)
	IncreaseProductStats(ctx context.Context, req *m.IncreaseProductStatsRequest) error
	IncreasePurchaseCountInProductStats(ctx context.Context, productID uint) error
}

type productStatsService struct {
	repo ProductStatsRepository
}

func NewProductStatsService(repo ProductStatsRepository) ProductStatsService {
	return &productStatsService{
		repo: repo,
	}
}

func (s *productStatsService) GetAllProductStatsPaginated(pageStr string, pageSizeStr string, productIDFilter string, dateFromFilter string, dateToFilter string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	productID, _ := utils.ParseUint(productIDFilter)
	dateFrom := parseTimeFilter(dateFromFilter)
	dateTo := parseTimeFilter(dateToFilter)

	histories, total, err := s.repo.GetAllProductStatsPaginated(page, pageSize, productID, dateFrom, dateTo, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product stats: %w", err)
	}

	return utils.BuildPaginatedResponse(histories, total, page, pageSize), nil
}

func (s *productStatsService) GetProductStatsByID(id uint) (*ProductStatsEntity, error) {
	history, err := s.repo.GetProductStatsByID(id)
	if err != nil {
		return nil, fmt.Errorf("product stats not found: %w", err)
	}

	return history, nil
}

func (s *productStatsService) IncreaseProductStats(ctx context.Context, req *m.IncreaseProductStatsRequest) error {
	productStats, err := s.repo.GetProductStatsByID(req.ProductID)
	if err != nil || productStats == nil {
		productStats = &ProductStatsEntity{
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

	if err := s.repo.IncreaseProductStats(ctx, productStats); err != nil {
		return fmt.Errorf("failed to increase product stats: %w", err)
	}

	return nil
}

func (s *productStatsService) IncreasePurchaseCountInProductStats(ctx context.Context, productID uint) error {
	productStats, err := s.repo.GetProductStatsByID(productID)
	if err != nil || productStats == nil {
		productStats = &ProductStatsEntity{
			ProductID: productID,
		}
	}

	productStats.PurchaseCount = productStats.PurchaseCount + 1
	if err := s.repo.IncreaseProductStats(ctx, productStats); err != nil {
		return fmt.Errorf("failed to increase product stats: %w", err)
	}

	return nil
}

func parseTimeFilter(timeStr string) *time.Time {
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

	l.Logger.Warn("Failed to parse time filter", "timeStr", timeStr)
	return nil
}
