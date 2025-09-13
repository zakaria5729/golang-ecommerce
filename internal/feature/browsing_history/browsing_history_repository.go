package browsing_history

import (
	"strings"
	"time"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type BrowsingHistoryRepository struct {
	db *gorm.DB
}

func NewBrowsingHistoryRepository() *BrowsingHistoryRepository {
	return &BrowsingHistoryRepository{
		db: db.GetDB(),
	}
}

func (r *BrowsingHistoryRepository) GetAllBrowsingHistory(userID uint, include []string, productID *uint, dateFrom *time.Time, dateTo *time.Time, sortBy, sortOrder string) ([]BrowsingHistory, error) {
	var history []BrowsingHistory

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(BrowsingHistoryUserID+" = ?", userID)

	if productID != nil {
		query = query.Where(BrowsingHistoryProductID+" = ?", *productID)
	}

	if dateFrom != nil {
		query = query.Where(BrowsingHistoryViewedAt+" >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where(BrowsingHistoryViewedAt+" <= ?", *dateTo)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &[]string{BrowsingHistoryViewedAt}); orderClause != "" {
		query = query.Order(orderClause)
	} else {
		query = query.Order(BrowsingHistoryViewedAt + " " + constants.SortOrderDesc)
	}

	err := query.Find(&history).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch browsing history", "method", "GetAllBrowsingHistory", "error", err, "userID", userID, "include", include, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return history, err
}

func (r *BrowsingHistoryRepository) GetAllBrowsingHistoryPaginated(userID uint, include []string, page, pageSize int, productID *uint, dateFrom *time.Time, dateTo *time.Time, sortBy, sortOrder string) ([]BrowsingHistory, int, error) {
	var history []BrowsingHistory
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(BrowsingHistoryUserID+" = ?", userID)

	if productID != nil {
		query = query.Where(BrowsingHistoryProductID+" = ?", *productID)
	}

	if dateFrom != nil {
		query = query.Where(BrowsingHistoryViewedAt+" >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where(BrowsingHistoryViewedAt+" <= ?", *dateTo)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &[]string{BrowsingHistoryViewedAt}); orderClause != "" {
		query = query.Order(orderClause)
	} else {
		query = query.Order(BrowsingHistoryViewedAt + " " + constants.SortOrderDesc)
	}

	if err := query.Model(&BrowsingHistory{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count browsing history", "method", "GetAllBrowsingHistoryPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&history).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch browsing history paginated", "method", "GetAllBrowsingHistoryPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return history, int(total), err
}

func (r *BrowsingHistoryRepository) GetBrowsingHistoryByID(id uint, userID uint, include []string) (*BrowsingHistory, error) {
	var history BrowsingHistory

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(BrowsingHistoryUserID+" = ?", userID)

	if err := query.Where(constants.FieldID+" = ?", id).First(&history).Error; err != nil {
		logger.Logger.Error("Failed to fetch browsing history by ID", "method", "GetBrowsingHistoryByID", "error", err, "id", id, "userID", userID, "include", include)
		return nil, err
	}

	return &history, nil
}

func (r *BrowsingHistoryRepository) CreateBrowsingHistory(history *BrowsingHistory) error {
	if history.ViewedAt.IsZero() {
		history.ViewedAt = timeutil.NowUTC()
	}

	err := r.db.Create(history).Error
	if err != nil {
		logger.Logger.Error("Failed to create browsing history", "method", "CreateBrowsingHistory", "error", err, "history", history)
	}
	return err
}

func (r *BrowsingHistoryRepository) DeleteBrowsingHistory(id uint, userID uint) error {
	err := r.db.Where(constants.FieldID+" = ? AND "+BrowsingHistoryUserID+" = ?", id, userID).Delete(&BrowsingHistory{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete browsing history", "method", "DeleteBrowsingHistory", "error", err, "id", id, "userID", userID)
	}
	return err
}

func (r *BrowsingHistoryRepository) ClearBrowsingHistory(userID uint) error {
	err := r.db.Where(BrowsingHistoryUserID+" = ?", userID).Delete(&BrowsingHistory{}).Error
	if err != nil {
		logger.Logger.Error("Failed to clear browsing history", "method", "ClearBrowsingHistory", "error", err, "userID", userID)
	}
	return err
}

func (r *BrowsingHistoryRepository) GetRecentBrowsingHistory(userID uint, limit int) ([]BrowsingHistory, error) {
	var history []BrowsingHistory

	query := r.db.Where(BrowsingHistoryUserID+" = ?", userID)
	query = query.Order(BrowsingHistoryViewedAt + " " + constants.SortOrderDesc)
	query = query.Limit(limit)

	err := query.Find(&history).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch recent browsing history", "method", "GetRecentBrowsingHistory", "error", err, "userID", userID, "limit", limit)
	}
	return history, err
}

func (r *BrowsingHistoryRepository) GetMostViewedProducts(userID uint, limit int) ([]uint, error) {
	var productIDs []uint

	query := r.db.Model(&BrowsingHistory{}).
		Select(BrowsingHistoryProductID).
		Where(BrowsingHistoryUserID+" = ?", userID).
		Group(BrowsingHistoryProductID).
		Order("COUNT(*) " + constants.SortOrderDesc).
		Limit(limit)

	err := query.Pluck(BrowsingHistoryProductID, &productIDs).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch most viewed products", "method", "GetMostViewedProducts", "error", err, "userID", userID, "limit", limit)
	}
	return productIDs, err
}

func (r *BrowsingHistoryRepository) BrowsingHistoryExists(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&BrowsingHistory{}).Where(constants.FieldID+" = ? AND "+BrowsingHistoryUserID+" = ?", id, userID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if browsing history exists", "method", "BrowsingHistoryExists", "error", err, "id", id, "userID", userID)
	}
	return count > 0, err
}

func (r *BrowsingHistoryRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, BrowsingHistoryUserID, BrowsingHistoryProductID, BrowsingHistoryViewedAt, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
