package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/size_option"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type SizeOptionHandler struct {
	sizeOptionUseCase *size_option.SizeOptionUseCase
}

func NewSizeOptionHandler() *SizeOptionHandler {
	return &SizeOptionHandler{
		sizeOptionUseCase: size_option.NewSizeOptionUseCase(),
	}
}

func (h *SizeOptionHandler) GetAllSizeOptionsPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeOptions := h.getAllSizeOptionsData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeOptions)
}

func (h *SizeOptionHandler) GetSizeOptionByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeOption := h.getSizeOptionDataById(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeOption)
}

func (h *SizeOptionHandler) GetAllSizeOptions(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeOptions := h.getAllSizeOptionsData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeOptions)
}

func (h *SizeOptionHandler) GetSizeOptionByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeOption := h.getSizeOptionDataById(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeOption)
}

func (h *SizeOptionHandler) CreateSizeOption(w http.ResponseWriter, r *http.Request) {
	var req size_option.CreateSizeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "CreateSizeOption") {
		return
	}

	if validationErrors := h.validateCreateSizeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeOption, err := h.sizeOptionUseCase.CreateSizeOption(&req)
	if err != nil {
		logger.Logger.Error("Failed to create size option", "method", "CreateSizeOption", "error", err, "request", req)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeOption)
}

func (h *SizeOptionHandler) UpdateSizeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size option ID", http.StatusBadRequest)
		return
	}

	var req size_option.UpdateSizeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateSizeOption") {
		return
	}

	if validationErrors := h.validateUpdateSizeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeOption, err := h.sizeOptionUseCase.UpdateSizeOption(*id, &req)
	if err != nil {
		logger.Logger.Error("Failed to update size option", "method", "UpdateSizeOption", "error", err, "id", *id, "request", req)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, sizeOption)
}

func (h *SizeOptionHandler) DeleteSizeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size option ID", http.StatusBadRequest)
		return
	}

	err = h.sizeOptionUseCase.DeleteSizeOption(*id)
	if err != nil {
		logger.Logger.Error("Failed to delete size option", "method", "DeleteSizeOption", "error", err, "id", *id)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Size option deleted successfully")
}

func (h *SizeOptionHandler) UndoDeletedSizeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size option ID", http.StatusBadRequest)
		return
	}

	err = h.sizeOptionUseCase.UndoDeletedSizeOption(*id)
	if err != nil {
		logger.Logger.Error("Failed to undo deleted size option", "method", "UndoDeletedSizeOption", "error", err, "id", *id)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Size option restored successfully"})
}

func (h *SizeOptionHandler) getAllSizeOptionsData(r *http.Request, showDeleted *bool) (error, []size_option.SizeOption) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	sizeCategoryID := q.Get("size_category_id")
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	sizeOptions, err := h.sizeOptionUseCase.GetAllSizeOptions(include, showDeleted, sizeCategoryID, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all size options"), nil
	}

	return nil, sizeOptions
}

func (h *SizeOptionHandler) getSizeOptionDataById(r *http.Request, showDeleted *bool) (error, *size_option.SizeOption) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid size option ID"), nil
	}

	q := r.URL.Query()
	include := q.Get(c.Include)

	sizeOption, err := h.sizeOptionUseCase.GetSizeOptionByID(*id, include, showDeleted)
	if err != nil {
		return errors.New("failed to get size option by ID"), nil
	}

	return nil, sizeOption
}

func (h *SizeOptionHandler) validateCreateSizeOptionRequest(req *size_option.CreateSizeOptionRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
		validator.ValidatePositiveInteger(req.SizeCategoryID, "size_category_id"),
	)
}

func (h *SizeOptionHandler) validateUpdateSizeOptionRequest(req *size_option.UpdateSizeOptionRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
