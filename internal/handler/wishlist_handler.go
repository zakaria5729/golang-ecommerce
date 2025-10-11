package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/wishlist"
	c "github.com/easy-comerce/backend/pkg/constants"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type WishlistHandler struct {
	service *wishlist.WishlistService
}

func NewWishlistHandler(service *wishlist.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		service: service,
	}
}

func (h *WishlistHandler) GetAllWishlistsPaginatedByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	productID, _ := utils.ParseUint(q.Get(c.WishlistProductID))
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	paginatedResponse, err := h.service.GetAllWishlistsPaginated(nil, userID, productID, pageStr, pageSizeStr, sortBy, sortOrder)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *WishlistHandler) AddToWishlistsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	productID, err := utils.ParseUint(r.PathValue(c.WishlistProductID))
	if err != nil || productID == nil || *productID == 0 {
		response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.AddToWishlistsByUser(*userID, *productID)
	response.SendResponse(w, "Product added to wishlist successfully", err, http.StatusInternalServerError)
}

func (h *WishlistHandler) RemoveFromWishlistByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	productID, err := utils.ParseUint(r.PathValue(c.WishlistProductID))
	if err != nil || productID == nil || *productID == 0 {
		response.SendErrorJSON(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	err = h.service.RemoveFromWishlistByUser(*userID, *productID)
	response.SendResponse(w, "Product removed from wishlist successfully", err, http.StatusInternalServerError)
}

func (h *WishlistHandler) ClearUserWishlist(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	err = h.service.ClearUserWishlist(*userID)
	response.SendResponse(w, "Wishlist cleared successfully", err, http.StatusInternalServerError)
}

func (h *WishlistHandler) DeleteWishlistById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteWishlistById(*id)
	response.SendResponse(w, "Wishlist delete successfully", err, http.StatusInternalServerError)
}

func (h *WishlistHandler) UndoDeleteWishlistById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeleteWishlistById(*id)
	response.SendResponse(w, "Wishlist delete successfully", err, http.StatusInternalServerError)
}

func (h *WishlistHandler) GetAllWishlistsPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	userIDFilter, _ := utils.ParseUint(q.Get(c.WishlistUserID))
	productIDFilter, _ := utils.ParseUint(q.Get(c.WishlistProductID))
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	paginatedResponse, err := h.service.GetAllWishlistsPaginated(showDeleted, userIDFilter, productIDFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *WishlistHandler) GetWishlistCountByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	count, err := h.service.GetWishlistCount(*userID)
	response.SendResponse(w, map[string]int64{"count": count}, err, http.StatusInternalServerError)
}
