package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/attribute_option"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AttributeOptionHandler struct {
	attributeOptionUseCase *attribute_option.AttributeOptionUseCase
}

func NewAttributeOptionHandler() *AttributeOptionHandler {
	return &AttributeOptionHandler{
		attributeOptionUseCase: attribute_option.NewAttributeOptionUseCase(),
	}
}

func (h *AttributeOptionHandler) GetAllAttributeOptionsPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeOptions := h.getAllAttributeOptionsData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeOptions)
}

func (h *AttributeOptionHandler) GetAttributeOptionByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeOption := h.getAttributeOptionDataById(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeOption)
}

func (h *AttributeOptionHandler) GetAllAttributeOptions(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeOptions := h.getAllAttributeOptionsData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeOptions)
}

func (h *AttributeOptionHandler) GetAttributeOptionByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeOption := h.getAttributeOptionDataById(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeOption)
}

func (h *AttributeOptionHandler) CreateAttributeOption(w http.ResponseWriter, r *http.Request) {
	var req attribute_option.CreateAttributeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "CreateAttributeOption") {
		return
	}

	if validationErrors := h.validateCreateAttributeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeOption, err := h.attributeOptionUseCase.CreateAttributeOption(&req)
	if err != nil {
		logger.Logger.Error("Create attribute option failed", "method", "CreateAttributeOption", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, attributeOption, http.StatusCreated)
}

func (h *AttributeOptionHandler) UpdateAttributeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute option ID", http.StatusBadRequest)
		return
	}

	var req attribute_option.UpdateAttributeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateAttributeOption") {
		return
	}

	if validationErrors := h.validateUpdateAttributeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeOption, err := h.attributeOptionUseCase.UpdateAttributeOption(*id, &req)
	if err != nil {
		logger.Logger.Error("Update attribute option failed", "method", "UpdateAttributeOption", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, attributeOption)
}

func (h *AttributeOptionHandler) DeleteAttributeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute option ID", http.StatusBadRequest)
		return
	}

	if err := h.attributeOptionUseCase.DeleteAttributeOption(*id); err != nil {
		logger.Logger.Error("Delete attribute option failed", "method", "DeleteAttributeOption", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Attribute option deleted successfully")
}

func (h *AttributeOptionHandler) UndoDeletedAttributeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute option ID", http.StatusBadRequest)
		return
	}

	if err := h.attributeOptionUseCase.UndoDeletedAttributeOption(*id); err != nil {
		logger.Logger.Error("Undo deleted attribute option failed", "method", "UndoDeletedAttributeOption", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Attribute option restored successfully")
}

func (h *AttributeOptionHandler) validateCreateAttributeOptionRequest(req *attribute_option.CreateAttributeOptionRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidatePositiveInteger(req.AttributeTypeID, "attribute_type_id"),
		validator.ValidateRequired(req.AttributeOptionName, "attribute_option_name"),
		validator.ValidateMinLength(req.AttributeOptionName, "attribute_option_name", 2),
		validator.ValidateMaxLength(req.AttributeOptionName, "attribute_option_name", 100),
	)
}

func (h *AttributeOptionHandler) validateUpdateAttributeOptionRequest(req *attribute_option.UpdateAttributeOptionRequest) validator.ValidationErrors {
	errors := validator.MergeValidationErrors(
		validator.ValidateRequired(req.AttributeOptionName, "attribute_option_name"),
		validator.ValidateMinLength(req.AttributeOptionName, "attribute_option_name", 2),
		validator.ValidateMaxLength(req.AttributeOptionName, "attribute_option_name", 100),
	)

	if req.AttributeTypeID != nil {
		errors = validator.MergeValidationErrors(errors, validator.ValidatePositiveInteger(*req.AttributeTypeID, "attribute_type_id"))
	}

	return errors
}

func (h *AttributeOptionHandler) getAllAttributeOptionsData(r *http.Request, showDeleted *bool) (error, []attribute_option.AttributeOption) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	attributeTypeID := q.Get("attribute_type_id")
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	attributeOptions, err := h.attributeOptionUseCase.GetAllAttributeOptions(include, showDeleted, attributeTypeID, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all attribute options"), nil
	}

	return nil, attributeOptions
}

func (h *AttributeOptionHandler) getAttributeOptionDataById(r *http.Request, showDeleted *bool) (error, *attribute_option.AttributeOption) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid attribute option ID"), nil
	}

	include := r.URL.Query().Get(c.Include)
	attributeOption, err := h.attributeOptionUseCase.GetAttributeOptionByID(*id, include, showDeleted)
	if err != nil {
		return errors.New("attribute option not found"), nil
	}

	return nil, attributeOption
}
