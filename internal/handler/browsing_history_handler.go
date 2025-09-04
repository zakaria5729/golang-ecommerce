package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/browsing_history"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type BrowsingHistoryHandler struct {
	useCase *browsing_history.BrowsingHistoryUseCase
}

func NewBrowsingHistoryHandler() *BrowsingHistoryHandler {
	return &BrowsingHistoryHandler{
		useCase: browsing_history.NewBrowsingHistoryUseCase(),
	}
}

func (h *BrowsingHistoryHandler) getUserID() uint {
	// TODO: Get from JWT token or session in the future
	return uint(1)
}

func (h *BrowsingHistoryHandler) GetAllBrowsingHistory(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	productIDFilter := q.Get(browsing_history.BrowsingHistoryProductID)
	dateFromFilter := q.Get("date_from")
	dateToFilter := q.Get("date_to")
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	history, err := h.useCase.GetAllBrowsingHistory(userID, includeStr, productIDFilter, dateFromFilter, dateToFilter, sortBy, sortOrder)
	if err != nil {
		response.JSONError(w, "Failed to fetch browsing history", http.StatusInternalServerError)
		return
	}

	response.JSON(w, history)
}

func (h *BrowsingHistoryHandler) GetAllBrowsingHistoryPaginated(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	productIDFilter := q.Get(browsing_history.BrowsingHistoryProductID)
	dateFromFilter := q.Get("date_from")
	dateToFilter := q.Get("date_to")
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetAllBrowsingHistoryPaginated(userID, includeStr, pageStr, pageSizeStr, productIDFilter, dateFromFilter, dateToFilter, sortBy, sortOrder)
	if err != nil {
		response.JSONError(w, "Failed to fetch browsing history", http.StatusInternalServerError)
		return
	}

	response.JSON(w, paginatedResponse)
}

func (h *BrowsingHistoryHandler) GetBrowsingHistoryByID(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid browsing history ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	history, err := h.useCase.GetBrowsingHistoryByID(*id, userID, include)
	if err != nil {
		response.JSONError(w, "Browsing history not found", http.StatusNotFound)
		return
	}

	response.JSON(w, history)
}

func (h *BrowsingHistoryHandler) CreateBrowsingHistory(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	var req browsing_history.BrowsingHistory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	history, err := h.useCase.CreateBrowsingHistory(userID, &req)
	if err != nil {
		response.JSONError(w, "Failed to create browsing history", http.StatusInternalServerError)
		return
	}

	response.JSONCreated(w, history)
}

func (h *BrowsingHistoryHandler) DeleteBrowsingHistory(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid browsing history ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteBrowsingHistory(*id, userID); err != nil {
		response.JSONError(w, "Failed to delete browsing history", http.StatusInternalServerError)
		return
	}

	response.JSONNoContent(w)
}

func (h *BrowsingHistoryHandler) ClearBrowsingHistory(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	if err := h.useCase.ClearBrowsingHistory(userID); err != nil {
		response.JSONError(w, "Failed to clear browsing history", http.StatusInternalServerError)
		return
	}

	response.JSONNoContent(w)
}

func (h *BrowsingHistoryHandler) GetRecentBrowsingHistory(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := h.useCase.GetRecentBrowsingHistory(userID, limit)
	if err != nil {
		response.JSONError(w, "Failed to get recent browsing history", http.StatusInternalServerError)
		return
	}

	response.JSON(w, history)
}

func (h *BrowsingHistoryHandler) GetMostViewedProducts(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	productIDs, err := h.useCase.GetMostViewedProducts(userID, limit)
	if err != nil {
		response.JSONError(w, "Failed to get most viewed products", http.StatusInternalServerError)
		return
	}

	response.JSON(w, productIDs)
}
