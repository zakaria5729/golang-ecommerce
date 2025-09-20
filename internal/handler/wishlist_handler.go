package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/wishlist"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
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

	paginatedResponse, err := h.useCase.GetAllWishlistsPaginated(nil, userID, productID, includeStr, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch wishlists", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

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

func (h *WishlistHandler) DeleteWishlistById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteWishlistById(*id); err != nil {
		response.SendErrorJSON(w, "Failed to delete wishlist", http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Wishlist delete successfully")
}

func (h *WishlistHandler) UndoDeleteWishlistById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.UndoDeleteWishlistById(*id); err != nil {
		response.SendErrorJSON(w, "Failed to delete wishlist", http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Wishlist delete successfully")
}

func (h *WishlistHandler) GetAllWishlistsPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(c.Include)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	userIDFilter, _ := utils.ParseUint(q.Get(c.WishlistUserID))
	productIDFilter, _ := utils.ParseUint(q.Get(c.WishlistProductID))
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	paginatedResponse, err := h.useCase.GetAllWishlistsPaginated(showDeleted, userIDFilter, productIDFilter, includeStr, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch wishlists", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

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
