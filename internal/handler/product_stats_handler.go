package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/product_stats"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ProductStatsHandler struct {
	useCase *product_stats.ProductStatsUseCase
}

func NewProductStatsHandler() *ProductStatsHandler {
	return &ProductStatsHandler{
		useCase: product_stats.NewProductStatsUseCase(),
	}
}

// **REQUIRED
func (h *ProductStatsHandler) GetAllProductStatsPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	productIDFilter := q.Get(c.SortOrder)
	dateFromFilter := q.Get("date_from")
	dateToFilter := q.Get(c.SortOrder)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	paginatedResponse, err := h.useCase.GetAllProductStatsPaginated(pageStr, pageSizeStr, productIDFilter, dateFromFilter, dateToFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch product stats", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

// **REQUIRED
func (h *ProductStatsHandler) GetProductStatsByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid product stats ID", http.StatusBadRequest)
		return
	}

	history, err := h.useCase.GetProductStatsByID(*id)
	if err != nil {
		response.SendErrorJSON(w, "Product stats not found", http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, history)
}

// **REQUIRED
func (h *ProductStatsHandler) IncreaseProductStats(w http.ResponseWriter, r *http.Request) {
	var req product_stats.IncreaseProductStatsRequest
	if !utils.DecodeJSON(w, r, &req, "IncreaseProductStats") {
		return
	}

	if validationErrors := h.validateProductStatsRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.useCase.IncreaseProductStats(&req)
	if err != nil {
		response.SendErrorJSON(w, "Failed to increase product stats", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, "Product stats increased successfully")
}

// **REQUIRED
func (h *ProductStatsHandler) validateProductStatsRequest(req *product_stats.IncreaseProductStatsRequest) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if req.ProductID <= 0 {
		errors.AddError("product_id", "product_id must be a positive integer")
	}

	return errors
}

// func (h *BrowsingHistoryHandler) DeleteBrowsingHistory(w http.ResponseWriter, r *http.Request) {
// 	userID := h.getUserID()

// 	id, err := utils.ParseUint(r.PathValue(c.FieldID))
// 	if err != nil || id == nil || *id == 0 {
// 		response.SendErrorJSON(w, "Invalid browsing history ID", http.StatusBadRequest)
// 		return
// 	}

// 	if userID == 0 {
// 		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.useCase.DeleteBrowsingHistory(*id, userID); err != nil {
// 		response.SendErrorJSON(w, "Failed to delete browsing history")
// 		return
// 	}

// 	response.SendDeleteJSON(w, "Browsing history deleted successfully")
// }

// func (h *BrowsingHistoryHandler) ClearBrowsingHistory(w http.ResponseWriter, r *http.Request) {
// 	userID := h.getUserID()

// 	if userID == 0 {
// 		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.useCase.ClearBrowsingHistory(userID); err != nil {
// 		response.SendErrorJSON(w, "Failed to clear browsing history")
// 		return
// 	}

// 	response.SendDeleteJSON(w, "Browsing history cleared successfully")
// }

// func (h *BrowsingHistoryHandler) GetRecentBrowsingHistory(w http.ResponseWriter, r *http.Request) {
// 	userID := h.getUserID()

// 	if userID == 0 {
// 		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	limit := h.getLimitFromQuery(r.URL.Query().Get("limit"))
// 	history, err := h.useCase.GetRecentBrowsingHistory(userID, limit)
// 	if err != nil {
// 		response.SendErrorJSON(w, "Failed to get recent browsing history", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccessJSON(w, history)
// }

// func (h *BrowsingHistoryHandler) GetMostViewedProducts(w http.ResponseWriter, r *http.Request) {
// 	userID := h.getUserID()

// 	if userID == 0 {
// 		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	limit := h.getLimitFromQuery(r.URL.Query().Get("limit"))
// 	productIDs, err := h.useCase.GetMostViewedProducts(userID, limit)
// 	if err != nil {
// 		response.SendErrorJSON(w, "Failed to get most viewed products", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccessJSON(w, productIDs)
// }

// func (h *BrowsingHistoryHandler) getLimitFromQuery(limitStr string) int {
// 	limit := h.getDefaultLimit()
// 	if limitStr != "" {
// 		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
// 			limit = l
// 		}
// 	}
// 	return limit
// }

// func (h *BrowsingHistoryHandler) getDefaultLimit() int {
// 	return 10
// }

// func (h *BrowsingHistoryHandler) GetAllBrowsingHistory(w http.ResponseWriter, r *http.Request) {
// 	userID := h.getUserID()

// 	q := r.URL.Query()
// 	includeStr := q.Get(constants.Include)
// 	productIDFilter := q.Get(browsing_history.BrowsingHistoryProductID)
// 	dateFromFilter := q.Get("date_from")
// 	dateToFilter := q.Get("date_to")
// 	sortBy := q.Get(constants.SortBy)
// 	sortOrder := q.Get(constants.SortOrder)

// 	history, err := h.useCase.GetAllBrowsingHistory(userID, includeStr, productIDFilter, dateFromFilter, dateToFilter, sortBy, sortOrder)
// 	if err != nil {
// 		response.SendErrorJSON(w, "Failed to fetch browsing history", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccessJSON(w, history)
// }
