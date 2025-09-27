package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/size_category"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type SizeCategoryHandler struct {
	sizeCategoryUseCase *size_category.SizeCategoryUseCase
}

func NewSizeCategoryHandler() *SizeCategoryHandler {
	return &SizeCategoryHandler{
		sizeCategoryUseCase: size_category.NewSizeCategoryUseCase(),
	}
}

func (h *SizeCategoryHandler) GetAllSizeCategoriesPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeCategories := h.getAllSizeCategoriesData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeCategories)
}

func (h *SizeCategoryHandler) GetSizeCategoryByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeCategory := h.getSizeCategoryDataById(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeCategory)
}

func (h *SizeCategoryHandler) GetAllSizeCategories(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeCategories := h.getAllSizeCategoriesData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeCategories)
}

func (h *SizeCategoryHandler) GetSizeCategoryByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeCategory := h.getSizeCategoryDataById(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeCategory)
}

func (h *SizeCategoryHandler) CreateSizeCategory(w http.ResponseWriter, r *http.Request) {
	var req size_category.CreateSizeCategoryRequest
	if !utils.DecodeJSON(w, r, &req, "CreateSizeCategory") {
		return
	}

	if validationErrors := h.validateCreateSizeCategoryRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeCategory, err := h.sizeCategoryUseCase.CreateSizeCategory(&req)
	if err != nil {
		logger.Logger.Error("Failed to create size category", "method", "CreateSizeCategory", "error", err, "request", req)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeCategory)
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

	if validationErrors := h.validateUpdateSizeCategoryRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeCategory, err := h.sizeCategoryUseCase.UpdateSizeCategory(*id, &req)
	if err != nil {
		logger.Logger.Error("Failed to update size category", "method", "UpdateSizeCategory", "error", err, "id", *id, "request", req)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeCategory)
}

func (h *SizeCategoryHandler) DeleteSizeCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size category ID", http.StatusBadRequest)
		return
	}

	err = h.sizeCategoryUseCase.DeleteSizeCategory(*id)
	if err != nil {
		logger.Logger.Error("Failed to delete size category", "method", "DeleteSizeCategory", "error", err, "id", *id)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Size category deleted successfully")
}

func (h *SizeCategoryHandler) UndoDeletedSizeCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size category ID", http.StatusBadRequest)
		return
	}

	err = h.sizeCategoryUseCase.UndoDeletedSizeCategory(*id)
	if err != nil {
		logger.Logger.Error("Failed to undo deleted size category", "method", "UndoDeletedSizeCategory", "error", err, "id", *id)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Size category restored successfully"})
}

func (h *SizeCategoryHandler) getAllSizeCategoriesData(r *http.Request, showDeleted *bool) (error, []size_category.SizeCategory) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	sizeCategories, err := h.sizeCategoryUseCase.GetAllSizeCategories(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all size categories"), nil
	}

	return nil, sizeCategories
}

func (h *SizeCategoryHandler) getSizeCategoryDataById(r *http.Request, showDeleted *bool) (error, *size_category.SizeCategory) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid size category ID"), nil
	}

	q := r.URL.Query()
	include := q.Get(c.Include)

	sizeCategory, err := h.sizeCategoryUseCase.GetSizeCategoryByID(*id, include, showDeleted)
	if err != nil {
		return errors.New("failed to get size category by ID"), nil
	}

	return nil, sizeCategory
}

func (h *SizeCategoryHandler) validateCreateSizeCategoryRequest(req *size_category.CreateSizeCategoryRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}

func (h *SizeCategoryHandler) validateUpdateSizeCategoryRequest(req *size_category.UpdateSizeCategoryRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
