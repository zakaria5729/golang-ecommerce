package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/wishlist"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type WishlistHandler struct {
	useCase *wishlist.WishlistUseCase
}

func NewWishlistHandler() *WishlistHandler {
	return &WishlistHandler{
		useCase: wishlist.NewWishlistUseCase(),
	}
}

func (h *WishlistHandler) GetAllWishlists(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	productIDFilter := q.Get(wishlist.WishlistProductID)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	wishlists, err := h.useCase.GetAllWishlists(userID, includeStr, productIDFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch wishlists", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, wishlists)
}

func (h *WishlistHandler) GetAllWishlistsPaginated(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	productIDFilter := q.Get(wishlist.WishlistProductID)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetAllWishlistsPaginated(userID, includeStr, pageStr, pageSizeStr, productIDFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch wishlists", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *WishlistHandler) GetWishlistByID(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	wishlist, err := h.useCase.GetWishlistByID(*id, userID, include)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendErrorJSON(w, "Wishlist not found", http.StatusNotFound)
			return
		}
		response.SendErrorJSON(w, "Failed to fetch wishlist", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, wishlist)
}

func (h *WishlistHandler) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	var req struct {
		ProductID string `json:"product_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorJSON(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	wishlist, err := h.useCase.CreateWishlist(userID, req.ProductID)
	if err != nil {
		if strings.Contains(err.Error(), "already in wishlist") {
			response.SendErrorJSON(w, "Product is already in wishlist", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "invalid product ID") {
			response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		response.SendErrorJSON(w, "Failed to create wishlist", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, wishlist, http.StatusCreated)
}

func (h *WishlistHandler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteWishlist(*id, userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendErrorJSON(w, "Wishlist not found", http.StatusNotFound)
			return
		}
		response.SendErrorJSON(w, "Failed to delete wishlist", http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Wishlist deleted successfully")
}

func (h *WishlistHandler) DeleteWishlistByProduct(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	productID := r.PathValue("product_id")
	if productID == "" {
		response.SendErrorJSON(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteWishlistByProduct(userID, productID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendErrorJSON(w, "Product not found in wishlist", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "invalid product ID") {
			response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		response.SendErrorJSON(w, "Failed to delete wishlist", http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Product removed from wishlist successfully")
}

func (h *WishlistHandler) ClearWishlist(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	if err := h.useCase.ClearUserWishlist(userID); err != nil {
		response.SendErrorJSON(w, "Failed to clear wishlist", http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Wishlist cleared successfully")
}

func (h *WishlistHandler) GetWishlistCount(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	count, err := h.useCase.GetWishlistCount(userID)
	if err != nil {
		response.SendErrorJSON(w, "Failed to get wishlist count", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]int64{"count": count})
}

func (h *WishlistHandler) GetWishlistByProduct(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	productID := r.PathValue("product_id")
	if productID == "" {
		response.SendErrorJSON(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	wishlist, err := h.useCase.GetWishlistByProduct(userID, productID, include)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendErrorJSON(w, "Product not found in wishlist", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "invalid product ID") {
			response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		response.SendErrorJSON(w, "Failed to fetch wishlist", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, wishlist)
}

func (h *WishlistHandler) IsProductInWishlist(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()
	productID := r.PathValue("product_id")
	if productID == "" {
		response.SendErrorJSON(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	isInWishlist, err := h.useCase.IsProductInWishlist(userID, productID)
	if err != nil {
		if strings.Contains(err.Error(), "invalid product ID") {
			response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		response.SendErrorJSON(w, "Failed to check wishlist", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]bool{"is_in_wishlist": isInWishlist})
}

func (h *WishlistHandler) getUserID() uint {
	// Hardcoded user ID for now - replace with actual user authentication
	return 1
}
