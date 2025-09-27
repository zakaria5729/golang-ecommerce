package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/brand"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type BrandHandler struct {
	brandUseCase *brand.BrandUseCase
}

func NewBrandHandler() *BrandHandler {
	return &BrandHandler{
		brandUseCase: brand.NewBrandUseCase(),
	}
}

func (h *BrandHandler) GetAllBrandsPublic(w http.ResponseWriter, r *http.Request) {
	err, brands := h.getAllBrandsData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, brands)
}

func (h *BrandHandler) GetBrandByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, brand := h.getBrandDataById(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, brand)
}

func (h *BrandHandler) GetAllBrands(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, brands := h.getAllBrandsData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, brands)
}

func (h *BrandHandler) GetBrandByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, brand := h.getBrandDataById(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, brand)
}

func (h *BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var req brand.CreateBrandRequest
	if !utils.DecodeJSON(w, r, &req, "CreateBrand") {
		return
	}

	if validationErrors := h.validateCreateBrandRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	brand, err := h.brandUseCase.CreateBrand(&req)
	if err != nil {
		logger.Logger.Error("Create brand failed", "method", "CreateBrand", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, brand, http.StatusCreated)
}

func (h *BrandHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	var req brand.UpdateBrandRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateBrand") {
		return
	}

	if validationErrors := h.validateUpdateBrandRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	brand, err := h.brandUseCase.UpdateBrand(*id, &req)
	if err != nil {
		logger.Logger.Error("Update brand failed", "method", "UpdateBrand", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, brand)
}

func (h *BrandHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	if err := h.brandUseCase.DeleteBrand(*id); err != nil {
		logger.Logger.Error("Delete brand failed", "method", "DeleteBrand", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Brand deleted successfully")
}

func (h *BrandHandler) UndoDeletedBrand(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	if err := h.brandUseCase.UndoDeletedBrand(*id); err != nil {
		logger.Logger.Error("Undo deleted brand failed", "method", "UndoDeletedBrand", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Brand restored successfully")
}

func (h *BrandHandler) getAllBrandsData(r *http.Request, showDeleted *bool) (error, []brand.Brand) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	brands, err := h.brandUseCase.GetAllBrands(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all brands"), nil
	}

	return nil, brands
}

func (h *BrandHandler) getBrandDataById(r *http.Request, showDeleted *bool) (error, *brand.Brand) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid brand ID"), nil
	}

	include := r.URL.Query().Get(c.Include)
	brand, err := h.brandUseCase.GetBrandByID(*id, include, showDeleted)
	if err != nil {
		return errors.New("brand not found"), nil
	}

	return nil, brand
}

func (h *BrandHandler) validateCreateBrandRequest(req *brand.CreateBrandRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 200),
	)
}

func (h *BrandHandler) validateUpdateBrandRequest(req *brand.UpdateBrandRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 200),
	)
}
