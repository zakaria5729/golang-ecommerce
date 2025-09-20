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

func (h *ProductStatsHandler) validateProductStatsRequest(req *product_stats.IncreaseProductStatsRequest) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if req.ProductID <= 0 {
		errors.AddError("product_id", "product_id must be a positive integer")
	}

	return errors
}
