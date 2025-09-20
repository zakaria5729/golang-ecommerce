package wishlist

import (
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type WishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository() *WishlistRepository {
	return &WishlistRepository{
		db: db.GetDB(),
	}
}

func (r *WishlistRepository) GetAllWishlistsPaginated(showDeleted *bool, userID *uint, productID *uint, include []string, page, pageSize int, sortBy, sortOrder string) ([]Wishlist, int, error) {
	var wishlists []Wishlist
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if userID != nil {
		query = query.Where(c.WishlistUserID+" = ?", *userID)
	}

	if productID != nil {
		query = query.Where(c.WishlistProductID+" = ?", *productID)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Model(&Wishlist{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count wishlists", "method", "GetAllWishlistsPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&wishlists).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch wishlists paginated", "method", "GetAllWishlistsPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "sortBy", sortOrder, "sortOrder", sortOrder)
	}

	return wishlists, int(total), err
}

func (r *WishlistRepository) AddToWishlistsByUser(wishlist *Wishlist) error {
	err := r.db.Create(wishlist).Error
	if err != nil {
		logger.Logger.Error("Failed to create wishlist", "method", "CreateWishlist", "error", err, "wishlist", wishlist)
	}
	return err
}

func (r *WishlistRepository) RemoveFromWishlistByUser(userID uint, productID uint) error {
	result := r.db.Where(c.WishlistUserID+" = ? AND "+c.WishlistProductID+" = ?", userID, productID).Delete(&Wishlist{})
	if result.Error != nil {
		logger.Logger.Error("Failed to delete wishlist by product", "method", "DeleteWishlistByProduct", "error", result.Error, "userID", userID, "productID", productID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("wishlist not found")
	}
	return nil
}

func (r *WishlistRepository) ClearUserWishlist(userID uint) error {
	result := r.db.Where(c.WishlistUserID+" = ?", userID).Delete(&Wishlist{})
	if result.Error != nil {
		logger.Logger.Error("Failed to clear user wishlist", "method", "ClearUserWishlist", "error", result.Error, "userID", userID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("wishlist not found")
	}
	return nil
}

func (r *WishlistRepository) DeleteWishlistById(id uint) error {
	result := r.db.Where(c.FieldID+" = ?", id).Delete(&Wishlist{})
	if result.Error != nil {
		logger.Logger.Error("Failed to delete wishlist by ID", "method", "DeleteWishlistById", "error", result.Error, "id", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("wishlist not found")
	}
	return nil
}

func (r *WishlistRepository) UndoDeleteWishlistById(id uint) error {
	result := r.db.Unscoped().Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil)
	if result.Error != nil {
		logger.Logger.Error("Failed to undo delete wishlist by ID", "method", "UndoDeleteWishlistById", "error", result.Error, "id", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("wishlist not found")
	}
	return nil
}

func (r *WishlistRepository) WishlistExists(userID uint, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Wishlist{}).Where(c.WishlistUserID+" = ? AND "+c.WishlistProductID+" = ?", userID, productID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if wishlist exists", "method", "WishlistExists", "error", err, "userID", userID, "productID", productID)
	}
	return count > 0, err
}

func (r *WishlistRepository) GetWishlistCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Wishlist{}).Where(c.WishlistUserID+" = ?", userID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to get wishlist count", "method", "GetWishlistCount", "error", err, "userID", userID)
	}
	return count, err
}

func (r *WishlistRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{c.FieldID, c.WishlistUserID, c.WishlistProductID, c.FieldCreatedAt, c.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
