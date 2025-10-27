package wishlist

import (
	"net/http"

	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type WishlistHandler interface {
	GetAllWishlistsPaginatedByUser(w http.ResponseWriter, r *http.Request)
	AddToWishlistsByUser(w http.ResponseWriter, r *http.Request)
	RemoveFromWishlistByUser(w http.ResponseWriter, r *http.Request)
	ClearUserWishlist(w http.ResponseWriter, r *http.Request)
	DeleteWishlistById(w http.ResponseWriter, r *http.Request)
	UndoDeleteWishlistById(w http.ResponseWriter, r *http.Request)
	GetWishlistCountByUser(w http.ResponseWriter, r *http.Request)
	GetAllWishlistsPaginated(w http.ResponseWriter, r *http.Request)
}

type wishlistHandler struct {
	service WishlistService
}

func NewWishlistHandler(service WishlistService) WishlistHandler {
	return &wishlistHandler{
		service: service,
	}
}

func (h *wishlistHandler) GetAllWishlistsPaginatedByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
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

func (h *wishlistHandler) AddToWishlistsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
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

func (h *wishlistHandler) RemoveFromWishlistByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
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

func (h *wishlistHandler) ClearUserWishlist(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	err = h.service.ClearUserWishlist(*userID)
	response.SendResponse(w, "Wishlist cleared successfully", err, http.StatusInternalServerError)
}

func (h *wishlistHandler) DeleteWishlistById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteWishlistById(*id)
	response.SendResponse(w, "Wishlist delete successfully", err, http.StatusInternalServerError)
}

func (h *wishlistHandler) UndoDeleteWishlistById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid wishlist ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeleteWishlistById(*id)
	response.SendResponse(w, "Wishlist delete successfully", err, http.StatusInternalServerError)
}

func (h *wishlistHandler) GetAllWishlistsPaginated(w http.ResponseWriter, r *http.Request) {
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

func (h *wishlistHandler) GetWishlistCountByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	count, err := h.service.GetWishlistCount(*userID)
	response.SendResponse(w, map[string]int64{"count": count}, err, http.StatusInternalServerError)
}
