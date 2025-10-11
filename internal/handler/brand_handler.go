package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/brand"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type BrandHandler struct {
	service *brand.BrandService
}

func NewBrandHandler(service *brand.BrandService) *BrandHandler {
	return &BrandHandler{
		service: service,
	}
}

func (h *BrandHandler) GetAllBrandsPublic(w http.ResponseWriter, r *http.Request) {
	err, brands := getAllBrandsData(r, nil, h.service)
	response.SendResponse(w, brands, err, http.StatusInternalServerError)
}

func (h *BrandHandler) GetBrandByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, brand := getBrandDataById(r, nil, h.service)
	response.SendResponse(w, brand, err, http.StatusInternalServerError)
}

func (h *BrandHandler) GetAllBrands(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, brands := getAllBrandsData(r, showDeleted, h.service)
	response.SendResponse(w, brands, err, http.StatusInternalServerError)
}

func (h *BrandHandler) GetBrandByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, brand := getBrandDataById(r, showDeleted, h.service)
	response.SendResponse(w, brand, err, http.StatusInternalServerError)
}

func (h *BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var req brand.CreateBrandRequest
	if !utils.DecodeJSON(w, r, &req, "CreateBrand") {
		return
	}

	if validationErrors := validateCreateBrandRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	brand, err := h.service.CreateBrand(&req)
	response.SendResponse(w, brand, err, http.StatusInternalServerError)
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

	if validationErrors := validateUpdateBrandRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	brand, err := h.service.UpdateBrand(*id, &req)
	response.SendResponse(w, brand, err, http.StatusInternalServerError)
}

func (h *BrandHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteBrand(*id)
	response.SendResponse(w, "Brand deleted successfully", err, http.StatusInternalServerError)
}

func (h *BrandHandler) UndoDeletedBrand(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedBrand(*id)
	response.SendResponse(w, "Brand restored successfully", err, http.StatusInternalServerError)
}

func getAllBrandsData(r *http.Request, showDeleted *bool, service *brand.BrandService) (error, []brand.Brand) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	brands, err := service.GetAllBrands(showDeleted, sortBy, sortOrder)

	if err != nil {
		return errors.New("failed to get all brands"), nil
	}

	return nil, brands
}

func getBrandDataById(r *http.Request, showDeleted *bool, service *brand.BrandService) (error, *brand.Brand) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid brand ID"), nil
	}

	brand, err := service.GetBrandByID(*id, showDeleted)
	if err != nil {
		return errors.New("brand not found"), nil
	}

	return nil, brand
}

func validateCreateBrandRequest(req *brand.CreateBrandRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 200),
	)
}

func validateUpdateBrandRequest(req *brand.UpdateBrandRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 200),
	)
}
