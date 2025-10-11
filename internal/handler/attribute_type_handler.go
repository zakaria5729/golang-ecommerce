package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/attribute_type"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AttributeTypeHandler struct {
	service *attribute_type.AttributeTypeService
}

func NewAttributeTypeHandler(service *attribute_type.AttributeTypeService) *AttributeTypeHandler {
	return &AttributeTypeHandler{
		service: service,
	}
}

func (h *AttributeTypeHandler) GetAllAttributeTypesPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeTypes := getAllAttributeTypesData(r, nil, h.service)
	response.SendResponse(w, attributeTypes, err, http.StatusInternalServerError)
}

func (h *AttributeTypeHandler) GetAttributeTypeByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeType := getAttributeTypeDataById(r, nil, h.service)
	response.SendResponse(w, attributeType, err, http.StatusInternalServerError)
}

func (h *AttributeTypeHandler) GetAllAttributeTypes(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeTypes := getAllAttributeTypesData(r, showDeleted, h.service)
	response.SendResponse(w, attributeTypes, err, http.StatusInternalServerError)
}

func (h *AttributeTypeHandler) GetAttributeTypeByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeType := getAttributeTypeDataById(r, showDeleted, h.service)
	response.SendResponse(w, attributeType, err, http.StatusInternalServerError)
}

func (h *AttributeTypeHandler) CreateAttributeType(w http.ResponseWriter, r *http.Request) {
	var req attribute_type.CreateAttributeTypeRequest
	if !utils.DecodeJSON(w, r, &req, "CreateAttributeType") {
		return
	}

	if validationErrors := validateCreateAttributeTypeRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeType, err := h.service.CreateAttributeType(&req)
	response.SendResponse(w, attributeType, err, http.StatusInternalServerError)
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

	if validationErrors := validateUpdateAttributeTypeRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeType, err := h.service.UpdateAttributeType(*id, &req)
	response.SendResponse(w, attributeType, err, http.StatusInternalServerError)
}

func (h *AttributeTypeHandler) DeleteAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteAttributeType(*id)
	response.SendResponse(w, "Attribute type deleted successfully", err, http.StatusInternalServerError)
}

func (h *AttributeTypeHandler) UndoDeletedAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedAttributeType(*id)
	response.SendResponse(w, "Attribute type restored successfully", err, http.StatusInternalServerError)
}

func getAllAttributeTypesData(r *http.Request, showDeleted *bool, service *attribute_type.AttributeTypeService) (error, []attribute_type.AttributeType) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	attributeTypes, err := service.GetAllAttributeTypes(showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all attribute types"), nil
	}

	return nil, attributeTypes
}

func getAttributeTypeDataById(r *http.Request, showDeleted *bool, service *attribute_type.AttributeTypeService) (error, *attribute_type.AttributeType) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid attribute type ID"), nil
	}

	attributeType, err := service.GetAttributeTypeByID(*id, showDeleted)
	if err != nil {
		return errors.New("attribute type not found"), nil
	}

	return nil, attributeType
}

func validateCreateAttributeTypeRequest(req *attribute_type.CreateAttributeTypeRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}

func validateUpdateAttributeTypeRequest(req *attribute_type.UpdateAttributeTypeRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
