package wishlist

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type WishlistUseCase struct {
	repo *WishlistRepository
}

func NewWishlistUseCase() *WishlistUseCase {
	return &WishlistUseCase{
		repo: NewWishlistRepository(),
	}
}

func (uc *WishlistUseCase) GetAllWishlistsPaginated(showDeleted *bool, userID *uint, productID *uint, includeStr string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)

	wishlists, total, err := uc.repo.GetAllWishlistsPaginated(showDeleted, userID, productID, include, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch wishlists paginated", "method", "GetAllWishlistsPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch wishlists: %w", err)
	}

	return utils.BuildPaginatedResponse(wishlists, total, page, pageSize), nil
}

func (uc *WishlistUseCase) AddToWishlistsByUser(userID uint, productID uint) error {
	exists, err := uc.repo.WishlistExists(userID, productID)
	if err != nil {
		logger.Logger.Error("Failed to check wishlist exists", "method", "CreateWishlist", "error", err, "userID", userID, "productID", productID)
		return fmt.Errorf("failed to check existing wishlist: %w", err)
	}

	if exists {
		logger.Logger.Error("Product already in wishlist", "method", "CreateWishlist", "userID", userID, "productID", productID)
		return errors.New("product is already in wishlist")
	}

	wishlist := &Wishlist{
		UserID:    userID,
		ProductID: productID,
	}

	err = uc.repo.AddToWishlistsByUser(wishlist)
	if err != nil {
		logger.Logger.Error("Failed to create wishlist", "method", "CreateWishlist", "error", err, "wishlist", wishlist)
		return fmt.Errorf("failed to create wishlist: %w", err)
	}

	return nil
}

func (uc *WishlistUseCase) RemoveFromWishlistByUser(userID uint, productID uint) error {
	if err := uc.repo.RemoveFromWishlistByUser(userID, productID); err != nil {
		logger.Logger.Error("Failed to delete wishlist by product", "method", "DeleteWishlistByProduct", "error", err, "userID", userID, "productID", productID)
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

func (uc *WishlistUseCase) ClearUserWishlist(userID uint) error {
	if err := uc.repo.ClearUserWishlist(userID); err != nil {
		logger.Logger.Error("Failed to clear user wishlist", "method", "ClearUserWishlist", "error", err, "userID", userID)
		return fmt.Errorf("failed to clear wishlist: %w", err)
	}

	return nil
}

func (uc *WishlistUseCase) DeleteWishlistById(id uint) error {
	if err := uc.repo.DeleteWishlistById(id); err != nil {
		logger.Logger.Error("Failed to delete user wishlist", "method", "ClearUserWishlist", "error", err, "id", id)
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

func (uc *WishlistUseCase) UndoDeleteWishlistById(id uint) error {
	if err := uc.repo.UndoDeleteWishlistById(id); err != nil {
		logger.Logger.Error("Failed to delete user wishlist", "method", "ClearUserWishlist", "error", err, "id", id)
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

func (uc *WishlistUseCase) GetWishlistCount(userID uint) (int64, error) {
	count, err := uc.repo.GetWishlistCount(userID)
	if err != nil {
		logger.Logger.Error("Failed to get wishlist count", "method", "GetWishlistCount", "error", err, "userID", userID)
		return 0, fmt.Errorf("failed to get wishlist count: %w", err)
	}

	return count, nil
}
