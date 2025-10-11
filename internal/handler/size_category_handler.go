package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/size_category"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type SizeCategoryHandler struct {
	service *size_category.SizeCategoryService
}

func NewSizeCategoryHandler(service *size_category.SizeCategoryService) *SizeCategoryHandler {
	return &SizeCategoryHandler{
		service: service,
	}
}

func (h *SizeCategoryHandler) GetAllSizeCategoriesPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeCategories := getAllSizeCategoriesData(r, nil, h.service)
	response.SendResponse(w, sizeCategories, err, http.StatusInternalServerError)
}

func (h *SizeCategoryHandler) GetSizeCategoryByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeCategory := getSizeCategoryDataById(r, nil, h.service)
	response.SendResponse(w, sizeCategory, err, http.StatusInternalServerError)
}

func (h *SizeCategoryHandler) GetAllSizeCategories(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeCategories := getAllSizeCategoriesData(r, showDeleted, h.service)
	response.SendResponse(w, sizeCategories, err, http.StatusInternalServerError)
}

func (h *SizeCategoryHandler) GetSizeCategoryByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeCategory := getSizeCategoryDataById(r, showDeleted, h.service)
	response.SendResponse(w, sizeCategory, err, http.StatusInternalServerError)
}

func (h *SizeCategoryHandler) CreateSizeCategory(w http.ResponseWriter, r *http.Request) {
	var req size_category.CreateSizeCategoryRequest
	if !utils.DecodeJSON(w, r, &req, "CreateSizeCategory") {
		return
	}

	if validationErrors := validateCreateSizeCategoryRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeCategory, err := h.service.CreateSizeCategory(&req)
	response.SendResponse(w, sizeCategory, err, http.StatusInternalServerError)
}

func (h *SizeCategoryHandler) UpdateSizeCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size category ID", http.StatusBadRequest)
		return
	}

	var req size_category.UpdateSizeCategoryRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateSizeCategory") {
		return
	}

	if validationErrors := validateUpdateSizeCategoryRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeCategory, err := h.service.UpdateSizeCategory(*id, &req)
	response.SendResponse(w, sizeCategory, err, http.StatusInternalServerError)
}

func (h *SizeCategoryHandler) DeleteSizeCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size category ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSizeCategory(*id)
	response.SendResponse(w, "Size category deleted successfully", err, http.StatusInternalServerError)
}

func (h *SizeCategoryHandler) UndoDeletedSizeCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size category ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedSizeCategory(*id)
	response.SendResponse(w, "Size category restored successfully", err, http.StatusInternalServerError)
}

func getAllSizeCategoriesData(r *http.Request, showDeleted *bool, service *size_category.SizeCategoryService) (error, []size_category.SizeCategory) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	sizeCategories, err := service.GetAllSizeCategories(showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all size categories"), nil
	}

	return nil, sizeCategories
}

func getSizeCategoryDataById(r *http.Request, showDeleted *bool, service *size_category.SizeCategoryService) (error, *size_category.SizeCategory) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid size category ID"), nil
	}

	sizeCategory, err := service.GetSizeCategoryByID(*id, showDeleted)
	if err != nil {
		return errors.New("failed to get size category by ID"), nil
	}

	return nil, sizeCategory
}

func validateCreateSizeCategoryRequest(req *size_category.CreateSizeCategoryRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}

func validateUpdateSizeCategoryRequest(req *size_category.UpdateSizeCategoryRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
