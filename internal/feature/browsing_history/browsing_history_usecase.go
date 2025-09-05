package browsing_history

import (
	"fmt"
	"time"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type BrowsingHistoryUseCase struct {
	repo *BrowsingHistoryRepository
}

func NewBrowsingHistoryUseCase() *BrowsingHistoryUseCase {
	return &BrowsingHistoryUseCase{
		repo: NewBrowsingHistoryRepository(),
	}
}

func (uc *BrowsingHistoryUseCase) GetAllBrowsingHistory(userID uint, includeStr string, productIDFilter string, dateFromFilter string, dateToFilter string, sortBy, sortOrder string) ([]BrowsingHistory, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	productID, _ := utils.ParseUint(productIDFilter)
	dateFrom := uc.parseTimeFilter(dateFromFilter)
	dateTo := uc.parseTimeFilter(dateToFilter)

	history, err := uc.repo.GetAllBrowsingHistory(userID, include, productID, dateFrom, dateTo, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch browsing history", "method", "GetAllBrowsingHistory", "error", err, "userID", userID, "include", include, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch browsing history: %w", err)
	}

	return history, nil
}

func (uc *BrowsingHistoryUseCase) GetAllBrowsingHistoryPaginated(userID uint, includeStr string, pageStr string, pageSizeStr string, productIDFilter string, dateFromFilter string, dateToFilter string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	include := utils.ParseCommaSeparatedString(includeStr)
	productID, _ := utils.ParseUint(productIDFilter)
	dateFrom := uc.parseTimeFilter(dateFromFilter)
	dateTo := uc.parseTimeFilter(dateToFilter)

	history, total, err := uc.repo.GetAllBrowsingHistoryPaginated(userID, include, page, pageSize, productID, dateFrom, dateTo, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch browsing history paginated", "method", "GetAllBrowsingHistoryPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch browsing history: %w", err)
	}

	var historyPtrs []*BrowsingHistory
	for i := range history {
		historyPtrs = append(historyPtrs, &history[i])
	}

	return utils.BuildPaginatedResponse(historyPtrs, total, page, pageSize), nil
}

func (uc *BrowsingHistoryUseCase) GetBrowsingHistoryByID(id uint, userID uint, includeStr string) (*BrowsingHistory, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	history, err := uc.repo.GetBrowsingHistoryByID(id, userID, include)
	if err != nil {
		logger.Logger.Error("Browsing history not found", "method", "GetBrowsingHistoryByID", "error", err, "id", id, "userID", userID, "include", include)
		return nil, fmt.Errorf("browsing history not found: %w", err)
	}

	return history, nil
}

func (uc *BrowsingHistoryUseCase) CreateBrowsingHistory(userID uint, req *BrowsingHistory) (*BrowsingHistory, error) {
	if err := uc.validateCreateRequest(req); err.HasErrors() {
		logger.Logger.Error("Validation failed", "method", "CreateBrowsingHistory", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	history := &BrowsingHistory{
		UserID:    userID,
		ProductID: req.ProductID,
		ViewedAt:  req.ViewedAt,
	}

	if history.ViewedAt.IsZero() {
		history.ViewedAt = time.Now()
	}

	if err := uc.repo.CreateBrowsingHistory(history); err != nil {
		logger.Logger.Error("Failed to create browsing history", "method", "CreateBrowsingHistory", "error", err, "history", history)
		return nil, fmt.Errorf("failed to create browsing history: %w", err)
	}

	return history, nil
}

func (uc *BrowsingHistoryUseCase) DeleteBrowsingHistory(id uint, userID uint) error {
	_, err := uc.repo.GetBrowsingHistoryByID(id, userID, nil)
	if err != nil {
		logger.Logger.Error("Browsing history not found", "method", "DeleteBrowsingHistory", "error", err, "id", id, "userID", userID)
		return fmt.Errorf("browsing history not found: %w", err)
	}

	if err := uc.repo.DeleteBrowsingHistory(id, userID); err != nil {
		logger.Logger.Error("Failed to delete browsing history", "method", "DeleteBrowsingHistory", "error", err, "id", id, "userID", userID)
		return fmt.Errorf("failed to delete browsing history: %w", err)
	}

	return nil
}

func (uc *BrowsingHistoryUseCase) ClearBrowsingHistory(userID uint) error {
	if err := uc.repo.ClearBrowsingHistory(userID); err != nil {
		logger.Logger.Error("Failed to clear browsing history", "method", "ClearBrowsingHistory", "error", err, "userID", userID)
		return fmt.Errorf("failed to clear browsing history: %w", err)
	}

	return nil
}

func (uc *BrowsingHistoryUseCase) GetRecentBrowsingHistory(userID uint, limit int) ([]BrowsingHistory, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	history, err := uc.repo.GetRecentBrowsingHistory(userID, limit)
	if err != nil {
		logger.Logger.Error("Failed to get recent browsing history", "method", "GetRecentBrowsingHistory", "error", err, "userID", userID, "limit", limit)
		return nil, fmt.Errorf("failed to get recent browsing history: %w", err)
	}

	return history, nil
}

func (uc *BrowsingHistoryUseCase) GetMostViewedProducts(userID uint, limit int) ([]uint, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	productIDs, err := uc.repo.GetMostViewedProducts(userID, limit)
	if err != nil {
		logger.Logger.Error("Failed to get most viewed products", "method", "GetMostViewedProducts", "error", err, "userID", userID, "limit", limit)
		return nil, fmt.Errorf("failed to get most viewed products: %w", err)
	}

	return productIDs, nil
}

func (uc *BrowsingHistoryUseCase) parseTimeFilter(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}

	// Try parsing different time formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05",       // ISO 8601 without timezone
		"2006-01-02 15:04:05",       // Standard format
		"2006-01-02",                // Date only
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t
		}
	}

	logger.Logger.Warn("Failed to parse time filter", "timeStr", timeStr)
	return nil
}

func (uc *BrowsingHistoryUseCase) validateCreateRequest(req *BrowsingHistory) validator.ValidationErrors {
	var errors validator.ValidationErrors

	errors = validator.MergeValidationErrors(
		errors,
		validator.ValidateRequired(fmt.Sprintf("%d", req.ProductID), "product_id"),
	)

	if req.ProductID == 0 {
		errors.AddError("product_id", "product_id must be a positive integer")
	}

	return errors
}
