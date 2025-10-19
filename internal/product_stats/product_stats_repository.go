package product_stats

import (
	"context"
	"time"

	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type ProductStatsRepository struct {
	db *gorm.DB
}

func NewProductStatsRepository(db *gorm.DB) *ProductStatsRepository {
	return &ProductStatsRepository{
		db: db,
	}
}

func (r *ProductStatsRepository) GetAllProductStatsPaginated(page, pageSize int, productID *uint, dateFrom *time.Time, dateTo *time.Time, sortBy, sortOrder string) ([]ProductStatsEntity, int, error) {
	var history []ProductStatsEntity
	var total int64

	query := r.db.Model(&ProductStatsEntity{})

	if productID != nil {
		query = query.Where(c.ProductStatsProductID+" = ?", *productID)
	}

	if dateFrom != nil {
		query = query.Where(c.FieldCreatedAt+" >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where(c.FieldCreatedAt+" <= ?", *dateTo)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &[]string{c.ProductStatsViewCount}); orderClause != "" {
		query = query.Order(orderClause)
	} else {
		query = query.Order(c.ProductStatsViewCount + " " + c.SortOrderDesc)
	}

	if err := query.Model(&ProductStatsEntity{}).Count(&total).Error; err != nil {
		l.Logger.Error("Failed to count product stats", "method", "GetAllProductStatsPaginated", "error", err, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&history).Error
	if err != nil {
		l.Logger.Error("Failed to fetch product stats paginated", "method", "GetAllProductStatsPaginated", "error", err, "page", page, "pageSize", pageSize, "productID", productID, "dateFrom", dateFrom, "dateTo", dateTo, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return history, int(total), err
}

func (r *ProductStatsRepository) GetProductStatsByID(id uint) (*ProductStatsEntity, error) {
	var history ProductStatsEntity

	if err := r.db.Model(&ProductStatsEntity{}).Where(c.FieldID+" = ?", id).First(&history).Error; err != nil {
		l.Logger.Error("Failed to fetch product stats by ID", "method", "GetProductStatsByID", "error", err, "id", id)
		return nil, err
	}

	return &history, nil
}

func (r *ProductStatsRepository) IncreaseProductStats(ctx context.Context, productStats *ProductStatsEntity) error {
	userID, _ := cu.GetUserIDFromContext(ctx)

	if productStats.ID == 0 {
		productStats.CreatedBy = userID
		err := r.db.Model(&ProductStatsEntity{}).Create(productStats).Error

		if err != nil {
			l.Logger.Error("Failed to create product stats", "method", "IncreaseProductStats", "error", err)
		}
		return err
	}

	productStats.UpdatedBy = userID
	err := r.db.Model(&ProductStatsEntity{}).Where(c.FieldID+" = ?", productStats.ID).Updates(productStats).Error
	if err != nil {
		l.Logger.Error("Failed to update product stats", "method", "IncreaseProductStats", "error", err)
	}
	return err
}
