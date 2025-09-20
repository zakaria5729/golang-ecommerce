package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/wishlist"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type WishlistHandler struct {
	useCase *wishlist.WishlistUseCase
}

func NewWishlistHandler() *WishlistHandler {
	return &WishlistHandler{
		useCase: wishlist.NewWishlistUseCase(),
	}
}

// **REQUIRED
func (h *WishlistHandler) GetAllWishlistsPaginatedByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	productID, _ := utils.ParseUint(q.Get(c.WishlistProductID))
	includeStr := q.Get(c.Include)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	paginatedResponse, err := h.useCase.GetAllWishlistsPaginated(userID, productID, includeStr, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch wishlists", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

// **REQUIRED
func (h *WishlistHandler) AddToWishlistsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	productID, err := utils.ParseUint(r.PathValue(c.WishlistProductID))
	if err != nil || productID == nil || *productID == 0 {
		response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.useCase.AddToWishlistsByUser(*userID, *productID)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, "Product added to wishlist successfully")
}

// **REQUIRED
func (h *WishlistHandler) RemoveFromWishlistByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	productID, err := utils.ParseUint(r.PathValue(c.WishlistProductID))
	if err != nil || productID == nil || *productID == 0 {
		response.SendErrorJSON(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	if err := h.useCase.RemoveFromWishlistByUser(*userID, *productID); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Product removed from wishlist successfully")
}

// **REQUIRED
func (h *WishlistHandler) ClearUserWishlist(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if err := h.useCase.ClearUserWishlist(*userID); err != nil {
		response.SendErrorJSON(w, "Failed to clear wishlist", http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Wishlist cleared successfully")
}

// **REQUIRED
func (h *WishlistHandler) GetAllWishlistsPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(c.Include)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	userIDFilter, _ := utils.ParseUint(q.Get(c.WishlistUserID))
	productIDFilter, _ := utils.ParseUint(q.Get(c.WishlistProductID))

	paginatedResponse, err := h.useCase.GetAllWishlistsPaginated(userIDFilter, productIDFilter, includeStr, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch wishlists", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

// **REQUIRED
func (h *WishlistHandler) GetWishlistCountByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	count, err := h.useCase.GetWishlistCount(*userID)
	if err != nil {
		response.SendErrorJSON(w, "Failed to get wishlist count", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]int64{"count": count})
}

func (h *WishlistHandler) validateWishlistRequest(productID uint) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if productID == 0 {
		errors.AddError("product_id", "Product ID must be a positive integer")
	}

	return errors
}

// func (h *WishlistHandler) GetWishlistByID(w http.ResponseWriter, r *http.Request) {
// 	userID, err := middleware.GetUserIDFromContext(r)
// 	if err != nil || userID == nil || *userID == 0 {
// 		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
// 		return
// 	}

// 	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
// 	if err != nil || id == nil || *id == 0 {
// 		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
// 		return
// 	}

// 	include := r.URL.Query().Get(constants.Include)
// 	wishlist, err := h.useCase.GetWishlistByID(*id, *userID, include)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "not found") {
// 			response.SendErrorJSON(w, "Wishlist not found", http.StatusNotFound)
// 			return
// 		}
// 		response.SendErrorJSON(w, "Failed to fetch wishlist", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccessJSON(w, wishlist)
// }

// func (h *WishlistHandler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
// 	userID, err := middleware.GetUserIDFromContext(r)
// 	if err != nil || userID == nil || *userID == 0 {
// 		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
// 		return
// 	}

// 	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
// 	if err != nil || id == nil || *id == 0 {
// 		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.useCase.DeleteWishlist(*id, *userID); err != nil {
// 		if strings.Contains(err.Error(), "not found") {
// 			response.SendErrorJSON(w, "Wishlist not found", http.StatusNotFound)
// 			return
// 		}
// 		response.SendErrorJSON(w, "Failed to delete wishlist", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendDeleteJSON(w, "Wishlist deleted successfully")
// }

// func (h *WishlistHandler) GetWishlistByProduct(w http.ResponseWriter, r *http.Request) {
// 	userID, err := middleware.GetUserIDFromContext(r)
// 	if err != nil || userID == nil || *userID == 0 {
// 		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
// 		return
// 	}

// 	productID := r.PathValue("product_id")
// 	if productID == "" {
// 		response.SendErrorJSON(w, "Product ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	include := r.URL.Query().Get(constants.Include)
// 	wishlist, err := h.useCase.GetWishlistByProduct(*userID, productID, include)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "not found") {
// 			response.SendErrorJSON(w, "Product not found in wishlist", http.StatusNotFound)
// 			return
// 		}
// 		if strings.Contains(err.Error(), "invalid product ID") {
// 			response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
// 			return
// 		}
// 		response.SendErrorJSON(w, "Failed to fetch wishlist", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccessJSON(w, wishlist)
// }

// func (h *WishlistHandler) IsProductInWishlist(w http.ResponseWriter, r *http.Request) {
// 	userID, err := middleware.GetUserIDFromContext(r)
// 	if err != nil || userID == nil || *userID == 0 {
// 		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
// 		return
// 	}

// 	productID := r.PathValue("product_id")
// 	if productID == "" {
// 		response.SendErrorJSON(w, "Product ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	isInWishlist, err := h.useCase.IsProductInWishlist(*userID, productID)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "invalid product ID") {
// 			response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
// 			return
// 		}
// 		response.SendErrorJSON(w, "Failed to check wishlist", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccessJSON(w, map[string]bool{"is_in_wishlist": isInWishlist})
// }

// func (h *WishlistHandler) GetAllWishlists(w http.ResponseWriter, r *http.Request) {
// 	userID, err := middleware.GetUserIDFromContext(r)
// 	if err != nil || userID == nil || *userID == 0 {
// 		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
// 		return
// 	}

// 	q := r.URL.Query()
// 	includeStr := q.Get(c.Include)
// 	productIDFilter := q.Get(c.WishlistProductID)
// 	sortBy := q.Get(c.SortBy)
// 	sortOrder := q.Get(c.SortOrder)

// 	wishlists, err := h.useCase.GetAllWishlists(*userID, includeStr, productIDFilter, sortBy, sortOrder)
// 	if err != nil {
// 		response.SendErrorJSON(w, "Failed to fetch wishlists", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccessJSON(w, wishlists)
// }
