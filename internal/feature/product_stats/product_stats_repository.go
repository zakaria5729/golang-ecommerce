package product_stats

import (
	"time"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type ProductStatsRepository struct {
	db *gorm.DB
}

func NewProductStatsRepository() *ProductStatsRepository {
	return &ProductStatsRepository{
		db: db.GetDB(),
	}
}

// **REQUIRED
func (r *ProductStatsRepository) GetAllProductStatsPaginated(page, pageSize int, productID *uint, dateFrom *time.Time, dateTo *time.Time, sortBy, sortOrder string) ([]ProductStats, int, error) {
	var history []ProductStats
	var total int64

	query := r.db.Model(&ProductStats{})

	if productID != nil {
		query = query.Where(c.ProductStatsProductID+" = ?", *productID)
	}

	if dateFrom != nil {
		query = query.Where(c.ProductStatsViewCount+" >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where(c.ProductStatsViewCount+" <= ?", *dateTo)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &[]string{c.ProductStatsViewCount}); orderClause != "" {
		query = query.Order(orderClause)
	} else {
		query = query.Order(c.ProductStatsViewCount + " " + constants.SortOrderDesc)
	}

	if err := query.Model(&ProductStats{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count product stats", "method", "GetAllProductStatsPaginated", "error", err, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&history).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch product stats paginated", "method", "GetAllProductStatsPaginated", "error", err, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return history, int(total), err
}

// **REQUIRED
func (r *ProductStatsRepository) GetProductStatsByID(id uint) (*ProductStats, error) {
	var history ProductStats

	if err := r.db.Model(&ProductStats{}).Where(constants.FieldID+" = ?", id).First(&history).Error; err != nil {
		logger.Logger.Error("Failed to fetch product stats by ID", "method", "GetProductStatsByID", "error", err, "id", id)
		return nil, err
	}

	return &history, nil
}

// **REQUIRED
func (r *ProductStatsRepository) IncreaseProductStats(productStats *ProductStats) error {
	if productStats.ID == 0 {
		err := r.db.Create(productStats).Error
		if err != nil {
			logger.Logger.Error("Failed to create product stats", "method", "IncreaseProductStats", "error", err)
		}
		return err
	}

	err := r.db.Model(&ProductStats{}).Where(constants.FieldID+" = ?", productStats.ID).Updates(productStats).Error
	if err != nil {
		logger.Logger.Error("Failed to update product stats", "method", "IncreaseProductStats", "error", err)
	}
	return err
}

// func (r *ProductStatsRepository) DeleteBrowsingHistory(id uint, userID uint) error {
// 	err := r.db.Where(constants.FieldID+" = ? AND "+ProductStatsProductID+" = ?", id, userID).Delete(&ProductStats{}).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to delete browsing history", "method", "DeleteBrowsingHistory", "error", err, "id", id, "userID", userID)
// 	}
// 	return err
// }

// func (r *ProductStatsRepository) GetMostViewedProducts(userID uint, limit int) ([]uint, error) {
// 	var productIDs []uint

// 	query := r.db.Model(&ProductStats{}).
// 		Select(ProductStatsProductID).
// 		Where(ProductStatsProductID+" = ?", userID).
// 		Group(ProductStatsProductID).
// 		Order("COUNT(*) " + constants.SortOrderDesc).
// 		Limit(limit)

// 	err := query.Pluck(ProductStatsProductID, &productIDs).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch most viewed products", "method", "GetMostViewedProducts", "error", err, "userID", userID, "limit", limit)
// 	}
// 	return productIDs, err
// }

// func (r *ProductStatsRepository) BrowsingHistoryExists(id uint, userID uint) (bool, error) {
// 	var count int64
// 	err := r.db.Model(&ProductStats{}).Where(constants.FieldID+" = ? AND "+c.ProductStatsProductID+" = ?", id, userID).Count(&count).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to check if browsing history exists", "method", "BrowsingHistoryExists", "error", err, "id", id, "userID", userID)
// 	}
// 	return count > 0, err
// }

// func (r *ProductStatsRepository) getSelectableFields(include []string) []string {
// 	defaultFields := []string{constants.FieldID, c.ProductStatsProductID, c.ProductStatsViewCount, c.ProductStatsAddToCartCount, c.ProductStatsRemoveFromCartCount, c.ProductStatsPurchaseCount, c.ProductStatsWishlistCount, c.ProductStatsRemoveFromWishlistCount, constants.FieldCreatedAt, constants.FieldUpdatedAt}
// 	optionalFields := []string{}
// 	return utils.BuildSelectFields(defaultFields, optionalFields, include)
// }

// func (r *BrowsingHistoryRepository) GetAllBrowsingHistory(userID uint, include []string, productID *uint, dateFrom *time.Time, dateTo *time.Time, sortBy, sortOrder string) ([]BrowsingHistory, error) {
// 	var history []BrowsingHistory

// 	selectFields := r.getSelectableFields(include)
// 	query := r.db.Select(strings.Join(selectFields, ", "))
// 	query = query.Where(BrowsingHistoryUserID+" = ?", userID)

// 	if productID != nil {
// 		query = query.Where(BrowsingHistoryProductID+" = ?", *productID)
// 	}

// 	if dateFrom != nil {
// 		query = query.Where(BrowsingHistoryViewedAt+" >= ?", *dateFrom)
// 	}

// 	if dateTo != nil {
// 		query = query.Where(BrowsingHistoryViewedAt+" <= ?", *dateTo)
// 	}

// 	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &[]string{BrowsingHistoryViewedAt}); orderClause != "" {
// 		query = query.Order(orderClause)
// 	} else {
// 		query = query.Order(BrowsingHistoryViewedAt + " " + constants.SortOrderDesc)
// 	}

// 	err := query.Find(&history).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch browsing history", "method", "GetAllBrowsingHistory", "error", err, "userID", userID, "include", include, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
// 	}
// 	return history, err
// }

// func (r *BrowsingHistoryRepository) ClearBrowsingHistory(userID uint) error {
// 	err := r.db.Where(BrowsingHistoryUserID+" = ?", userID).Delete(&BrowsingHistory{}).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to clear browsing history", "method", "ClearBrowsingHistory", "error", err, "userID", userID)
// 	}
// 	return err
// }

// func (r *BrowsingHistoryRepository) GetRecentBrowsingHistory(userID uint, limit int) ([]BrowsingHistory, error) {
// 	var history []BrowsingHistory

// 	query := r.db.Where(BrowsingHistoryUserID+" = ?", userID)
// 	query = query.Order(BrowsingHistoryViewedAt + " " + constants.SortOrderDesc)
// 	query = query.Limit(limit)

// 	err := query.Find(&history).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch recent browsing history", "method", "GetRecentBrowsingHistory", "error", err, "userID", userID, "limit", limit)
// 	}
// 	return history, err
// }
