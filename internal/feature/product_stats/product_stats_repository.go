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

func (r *ProductStatsRepository) GetProductStatsByID(id uint) (*ProductStats, error) {
	var history ProductStats

	if err := r.db.Model(&ProductStats{}).Where(constants.FieldID+" = ?", id).First(&history).Error; err != nil {
		logger.Logger.Error("Failed to fetch product stats by ID", "method", "GetProductStatsByID", "error", err, "id", id)
		return nil, err
	}

	return &history, nil
}

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
