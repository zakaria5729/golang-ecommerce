package wishlist

import (
	"errors"
	"fmt"

	r "github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type WishlistService struct {
	repo *WishlistRepository
}

func NewWishlistService(repo *WishlistRepository) *WishlistService {
	return &WishlistService{
		repo: repo,
	}
}

func (s *WishlistService) GetAllWishlistsPaginated(showDeleted *bool, userID *uint, productID *uint, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	wishlists, total, err := s.repo.GetAllWishlistsPaginated(showDeleted, userID, productID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wishlists: %w", err)
	}

	return utils.BuildPaginatedResponse(wishlists, total, page, pageSize), nil
}

func (s *WishlistService) AddToWishlistsByUser(userID uint, productID uint) error {
	exists, err := s.repo.WishlistExists(userID, productID)
	if err != nil {
		return fmt.Errorf("failed to check existing wishlist: %w", err)
	}

	if exists {
		return errors.New("product is already in wishlist")
	}

	wishlist := &WishlistEntity{
		UserID:    userID,
		ProductID: productID,
	}

	err = s.repo.AddToWishlistsByUser(wishlist)
	if err != nil {
		return fmt.Errorf("failed to create wishlist: %w", err)
	}

	return nil
}

func (s *WishlistService) RemoveFromWishlistByUser(userID uint, productID uint) error {
	if err := s.repo.RemoveFromWishlistByUser(userID, productID); err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

func (s *WishlistService) ClearUserWishlist(userID uint) error {
	if err := s.repo.ClearUserWishlist(userID); err != nil {
		return fmt.Errorf("failed to clear wishlist: %w", err)
	}

	return nil
}

func (s *WishlistService) DeleteWishlistById(id uint) error {
	if err := s.repo.DeleteWishlistById(id); err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

func (s *WishlistService) UndoDeleteWishlistById(id uint) error {
	if err := s.repo.UndoDeleteWishlistById(id); err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

func (s *WishlistService) GetWishlistCount(userID uint) (int64, error) {
	count, err := s.repo.GetWishlistCount(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get wishlist count: %w", err)
	}

	return count, nil
}
