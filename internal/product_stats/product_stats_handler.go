package product_stats

import (
	"net/http"

	m "github.com/easy-comerce/backend/internal/product_stats/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ProductStatsHandler interface {
	GetAllProductStatsPaginated(w http.ResponseWriter, r *http.Request)
	GetProductStatsByID(w http.ResponseWriter, r *http.Request)
	IncreaseProductStats(w http.ResponseWriter, r *http.Request)
}

type productStatsHandler struct {
	service ProductStatsService
}

func NewProductStatsHandler(service ProductStatsService) ProductStatsHandler {
	return &productStatsHandler{
		service: service,
	}
}

func (h *productStatsHandler) GetAllProductStatsPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	productIDFilter := q.Get(c.SortOrder)
	dateFromFilter := q.Get(c.From)
	dateToFilter := q.Get(c.To)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	paginatedResponse, err := h.service.GetAllProductStatsPaginated(pageStr, pageSizeStr, productIDFilter, dateFromFilter, dateToFilter, sortBy, sortOrder)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *productStatsHandler) GetProductStatsByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid product stats ID", http.StatusBadRequest)
		return
	}

	history, err := h.service.GetProductStatsByID(*id)
	response.SendResponse(w, history, err, http.StatusInternalServerError)
}

func (h *productStatsHandler) IncreaseProductStats(w http.ResponseWriter, r *http.Request) {
	var req m.IncreaseProductStatsRequest
	if !utils.DecodeJSON(w, r, &req, "IncreaseProductStats") {
		return
	}

	if validationErrors := validateProductStatsRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.IncreaseProductStats(r.Context(), &req)
	response.SendResponse(w, "Product stats increased successfully", err, http.StatusInternalServerError)
}

func validateProductStatsRequest(req *m.IncreaseProductStatsRequest) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if req.ProductID <= 0 {
		errors.AddError("product_id", "product_id must be a positive integer")
	}

	return errors
}
