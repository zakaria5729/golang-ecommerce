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

// **REQUIRED
func (uc *WishlistUseCase) GetAllWishlistsPaginated(userID *uint, productID *uint, includeStr string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)

	wishlists, total, err := uc.repo.GetAllWishlistsPaginated(userID, productID, include, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch wishlists paginated", "method", "GetAllWishlistsPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "productID", productID, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch wishlists: %w", err)
	}

	return utils.BuildPaginatedResponse(wishlists, total, page, pageSize), nil
}

// **REQUIRED
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

// **REQUIRED
func (uc *WishlistUseCase) RemoveFromWishlistByUser(userID uint, productID uint) error {
	if err := uc.repo.DeleteWishlistByProduct(userID, productID); err != nil {
		logger.Logger.Error("Failed to delete wishlist by product", "method", "DeleteWishlistByProduct", "error", err, "userID", userID, "productID", productID)
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *WishlistUseCase) ClearUserWishlist(userID uint) error {
	if err := uc.repo.ClearUserWishlist(userID); err != nil {
		logger.Logger.Error("Failed to clear user wishlist", "method", "ClearUserWishlist", "error", err, "userID", userID)
		return fmt.Errorf("failed to clear wishlist: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *WishlistUseCase) GetWishlistCount(userID uint) (int64, error) {
	count, err := uc.repo.GetWishlistCount(userID)
	if err != nil {
		logger.Logger.Error("Failed to get wishlist count", "method", "GetWishlistCount", "error", err, "userID", userID)
		return 0, fmt.Errorf("failed to get wishlist count: %w", err)
	}

	return count, nil
}

// func (uc *WishlistUseCase) GetWishlistByID(id uint, userID uint, includeStr string) (*Wishlist, error) {
// 	include := utils.ParseCommaSeparatedString(includeStr)
// 	wishlist, err := uc.repo.GetWishlistByID(id, userID, include)
// 	if err != nil {
// 		logger.Logger.Error("Wishlist not found", "method", "GetWishlistByID", "error", err, "id", id, "userID", userID, "include", include)
// 		return nil, fmt.Errorf("wishlist not found: %w", err)
// 	}

// 	return wishlist, nil
// }

// func (uc *WishlistUseCase) DeleteWishlist(id uint, userID uint) error {
// 	if err := uc.repo.DeleteWishlist(id, userID); err != nil {
// 		logger.Logger.Error("Failed to delete wishlist", "method", "DeleteWishlist", "error", err, "id", id, "userID", userID)
// 		return fmt.Errorf("failed to delete wishlist: %w", err)
// 	}

// 	return nil
// }

// func (uc *WishlistUseCase) GetWishlistByProduct(userID uint, productIDStr string, includeStr string) (*Wishlist, error) {
// 	productID, err := utils.ParseUint(productIDStr)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid product ID: %w", err)
// 	}

// 	if *productID == 0 {
// 		logger.Logger.Error("Invalid product ID", "method", "GetWishlistByProduct", "productID", *productID)
// 		return nil, errors.New("product ID must be a positive integer")
// 	}

// 	include := utils.ParseCommaSeparatedString(includeStr)
// 	wishlist, err := uc.repo.GetWishlistByProduct(userID, *productID, include)
// 	if err != nil {
// 		logger.Logger.Error("Wishlist not found", "method", "GetWishlistByProduct", "error", err, "userID", userID, "productID", *productID, "include", include)
// 		return nil, fmt.Errorf("wishlist not found: %w", err)
// 	}

// 	return wishlist, nil
// }

// func (uc *WishlistUseCase) IsProductInWishlist(userID uint, productIDStr string) (bool, error) {
// 	productID, err := utils.ParseUint(productIDStr)
// 	if err != nil {
// 		return false, fmt.Errorf("invalid product ID: %w", err)
// 	}

// 	if *productID == 0 {
// 		logger.Logger.Error("Invalid product ID", "method", "IsProductInWishlist", "productID", *productID)
// 		return false, errors.New("product ID must be a positive integer")
// 	}

// 	exists, err := uc.repo.WishlistExists(userID, *productID)
// 	if err != nil {
// 		logger.Logger.Error("Failed to check if product is in wishlist", "method", "IsProductInWishlist", "error", err, "userID", userID, "productID", *productID)
// 		return false, fmt.Errorf("failed to check wishlist: %w", err)
// 	}

// 	return exists, nil
// }

// func (uc *WishlistUseCase) ValidateWishlistInput(req *Wishlist) validator.ValidationErrors {
// 	var errors validator.ValidationErrors

// 	errors = validator.MergeValidationErrors(
// 		errors,
// 		validator.ValidateRequired(fmt.Sprintf("%d", req.UserID), "user_id"),
// 		validator.ValidateRequired(fmt.Sprintf("%d", req.ProductID), "product_id"),
// 	)

// 	if req.UserID == 0 {
// 		errors.AddError("user_id", "user_id must be a positive integer")
// 	}

// 	if req.ProductID == 0 {
// 		errors.AddError("product_id", "product_id must be a positive integer")
// 	}

// 	return errors
// }

// func (uc *WishlistUseCase) GetAllWishlists(userID uint, includeStr string, productIDFilter string, sortBy, sortOrder string) ([]Wishlist, error) {
// 	include := utils.ParseCommaSeparatedString(includeStr)
// 	productID, _ := utils.ParseUint(productIDFilter)

// 	wishlists, err := uc.repo.GetAllWishlists(userID, include, productID, sortBy, sortOrder)
// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch wishlists", "method", "GetAllWishlists", "error", err, "userID", userID, "include", include, "productID", productID, "sortBy", sortBy, "sortOrder", sortOrder)
// 		return nil, fmt.Errorf("failed to fetch wishlists: %w", err)
// 	}

// 	return wishlists, nil
// }
