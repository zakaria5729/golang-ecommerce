package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/attribute_type"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AttributeTypeHandler struct {
	attributeTypeUseCase *attribute_type.AttributeTypeUseCase
}

func NewAttributeTypeHandler() *AttributeTypeHandler {
	return &AttributeTypeHandler{
		attributeTypeUseCase: attribute_type.NewAttributeTypeUseCase(),
	}
}

func (h *AttributeTypeHandler) GetAllAttributeTypesPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeTypes := h.getAllAttributeTypesData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeTypes)
}

func (h *AttributeTypeHandler) GetAttributeTypeByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeType := h.getAttributeTypeDataById(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeType)
}

func (h *AttributeTypeHandler) GetAllAttributeTypes(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeTypes := h.getAllAttributeTypesData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeTypes)
}

func (h *AttributeTypeHandler) GetAttributeTypeByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeType := h.getAttributeTypeDataById(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, attributeType)
}

func (h *AttributeTypeHandler) CreateAttributeType(w http.ResponseWriter, r *http.Request) {
	var req attribute_type.CreateAttributeTypeRequest
	if !utils.DecodeJSON(w, r, &req, "CreateAttributeType") {
		return
	}

	if validationErrors := h.validateCreateAttributeTypeRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeType, err := h.attributeTypeUseCase.CreateAttributeType(&req)
	if err != nil {
		logger.Logger.Error("Create attribute type failed", "method", "CreateAttributeType", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, attributeType, http.StatusCreated)
}

func (h *AttributeTypeHandler) UpdateAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	var req attribute_type.UpdateAttributeTypeRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateAttributeType") {
		return
	}

	if validationErrors := h.validateUpdateAttributeTypeRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeType, err := h.attributeTypeUseCase.UpdateAttributeType(*id, &req)
	if err != nil {
		logger.Logger.Error("Update attribute type failed", "method", "UpdateAttributeType", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, attributeType)
}

func (h *AttributeTypeHandler) DeleteAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	if err := h.attributeTypeUseCase.DeleteAttributeType(*id); err != nil {
		logger.Logger.Error("Delete attribute type failed", "method", "DeleteAttributeType", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Attribute type deleted successfully")
}

func (h *AttributeTypeHandler) UndoDeletedAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	if err := h.attributeTypeUseCase.UndoDeletedAttributeType(*id); err != nil {
		logger.Logger.Error("Undo deleted attribute type failed", "method", "UndoDeletedAttributeType", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Attribute type restored successfully")
}

func (h *AttributeTypeHandler) getAllAttributeTypesData(r *http.Request, showDeleted *bool) (error, []attribute_type.AttributeType) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	attributeTypes, err := h.attributeTypeUseCase.GetAllAttributeTypes(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all attribute types"), nil
	}

	return nil, attributeTypes
}

func (h *AttributeTypeHandler) getAttributeTypeDataById(r *http.Request, showDeleted *bool) (error, *attribute_type.AttributeType) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid attribute type ID"), nil
	}

	include := r.URL.Query().Get(c.Include)
	attributeType, err := h.attributeTypeUseCase.GetAttributeTypeByID(*id, include, showDeleted)
	if err != nil {
		return errors.New("attribute type not found"), nil
	}

	return nil, attributeType
}

func (h *AttributeTypeHandler) validateCreateAttributeTypeRequest(req *attribute_type.CreateAttributeTypeRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}

func (h *AttributeTypeHandler) validateUpdateAttributeTypeRequest(req *attribute_type.UpdateAttributeTypeRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
