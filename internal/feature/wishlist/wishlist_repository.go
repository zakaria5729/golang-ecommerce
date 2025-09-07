package wishlist

import (
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
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

func (r *WishlistRepository) GetAllWishlists(userID uint, include []string, productID *uint, sortBy, sortOrder string) ([]Wishlist, error) {
	var wishlists []Wishlist

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(WishlistUserID+" = ?", userID)

	if productID != nil {
		query = query.Where(WishlistProductID+" = ?", *productID)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&wishlists).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch wishlists", "method", "GetAllWishlists", "error", err, "userID", userID, "include", include, "productID", productID, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return wishlists, err
}

func (r *WishlistRepository) GetAllWishlistsPaginated(userID uint, include []string, page, pageSize int, productID *uint, sortBy, sortOrder string) ([]Wishlist, int, error) {
	var wishlists []Wishlist
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(WishlistUserID+" = ?", userID)

	if productID != nil {
		query = query.Where(WishlistProductID+" = ?", *productID)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Model(&Wishlist{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count wishlists", "method", "GetAllWishlistsPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&wishlists).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch wishlists paginated", "method", "GetAllWishlistsPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "sortBy", sortOrder, "sortOrder", sortOrder)
	}

	return wishlists, int(total), err
}

func (r *WishlistRepository) GetWishlistByID(id uint, userID uint, include []string) (*Wishlist, error) {
	var wishlist Wishlist

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(WishlistUserID+" = ?", userID)

	if err := query.Where(constants.FieldID+" = ?", id).First(&wishlist).Error; err != nil {
		logger.Logger.Error("Failed to fetch wishlist by ID", "method", "GetWishlistByID", "error", err, "id", id, "userID", userID, "include", include)
		return nil, err
	}

	return &wishlist, nil
}

func (r *WishlistRepository) CreateWishlist(wishlist *Wishlist) error {
	err := r.db.Create(wishlist).Error
	if err != nil {
		logger.Logger.Error("Failed to create wishlist", "method", "CreateWishlist", "error", err, "wishlist", wishlist)
	}
	return err
}

func (r *WishlistRepository) DeleteWishlist(id uint, userID uint) error {
	result := r.db.Where(constants.FieldID+" = ? AND "+WishlistUserID+" = ?", id, userID).Delete(&Wishlist{})
	if result.Error != nil {
		logger.Logger.Error("Failed to delete wishlist", "method", "DeleteWishlist", "error", result.Error, "id", id, "userID", userID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("wishlist not found")
	}
	return nil
}

func (r *WishlistRepository) DeleteWishlistByProduct(userID uint, productID uint) error {
	result := r.db.Where(WishlistUserID+" = ? AND "+WishlistProductID+" = ?", userID, productID).Delete(&Wishlist{})
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
	result := r.db.Where(WishlistUserID+" = ?", userID).Delete(&Wishlist{})
	if result.Error != nil {
		logger.Logger.Error("Failed to clear user wishlist", "method", "ClearUserWishlist", "error", result.Error, "userID", userID)
		return result.Error
	}
	return nil
}

func (r *WishlistRepository) WishlistExists(userID uint, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Wishlist{}).Where(WishlistUserID+" = ? AND "+WishlistProductID+" = ?", userID, productID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if wishlist exists", "method", "WishlistExists", "error", err, "userID", userID, "productID", productID)
	}
	return count > 0, err
}

func (r *WishlistRepository) GetWishlistCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Wishlist{}).Where(WishlistUserID+" = ?", userID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to get wishlist count", "method", "GetWishlistCount", "error", err, "userID", userID)
	}
	return count, err
}

func (r *WishlistRepository) GetWishlistByProduct(userID uint, productID uint, include []string) (*Wishlist, error) {
	var wishlist Wishlist

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(WishlistUserID+" = ? AND "+WishlistProductID+" = ?", userID, productID)

	if err := query.First(&wishlist).Error; err != nil {
		logger.Logger.Error("Failed to fetch wishlist by product", "method", "GetWishlistByProduct", "error", err, "userID", userID, "productID", productID, "include", include)
		return nil, err
	}

	return &wishlist, nil
}

func (r *WishlistRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, WishlistUserID, WishlistProductID, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
