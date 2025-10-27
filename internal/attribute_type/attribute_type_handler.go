package attribute_type

import (
	"errors"
	"net/http"

	m "github.com/easy-comerce/backend/internal/attribute_type/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AttributeTypeHandler interface {
	GetAllAttributeTypesPublic(w http.ResponseWriter, r *http.Request)
	GetAttributeTypeByIDPublic(w http.ResponseWriter, r *http.Request)
	GetAllAttributeTypes(w http.ResponseWriter, r *http.Request)
	GetAttributeTypeByID(w http.ResponseWriter, r *http.Request)
	CreateAttributeType(w http.ResponseWriter, r *http.Request)
	UpdateAttributeType(w http.ResponseWriter, r *http.Request)
	DeleteAttributeType(w http.ResponseWriter, r *http.Request)
	UndoDeletedAttributeType(w http.ResponseWriter, r *http.Request)
}

type attributeTypeHandler struct {
	service AttributeTypeService
}

func NewAttributeTypeHandler(service AttributeTypeService) AttributeTypeHandler {
	return &attributeTypeHandler{
		service: service,
	}
}

func (h *attributeTypeHandler) GetAllAttributeTypesPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeTypes := getAllAttributeTypesData(r, nil, h.service)
	response.SendResponse(w, attributeTypes, err, http.StatusInternalServerError)
}

func (h *attributeTypeHandler) GetAttributeTypeByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeType := getAttributeTypeDataById(r, nil, h.service)
	response.SendResponse(w, attributeType, err, http.StatusInternalServerError)
}

func (h *attributeTypeHandler) GetAllAttributeTypes(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeTypes := getAllAttributeTypesData(r, showDeleted, h.service)
	response.SendResponse(w, attributeTypes, err, http.StatusInternalServerError)
}

func (h *attributeTypeHandler) GetAttributeTypeByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeType := getAttributeTypeDataById(r, showDeleted, h.service)
	response.SendResponse(w, attributeType, err, http.StatusInternalServerError)
}

func (h *attributeTypeHandler) CreateAttributeType(w http.ResponseWriter, r *http.Request) {
	var req m.CreateAttributeTypeRequest
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

func (h *attributeTypeHandler) UpdateAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	var req m.UpdateAttributeTypeRequest
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

func (h *attributeTypeHandler) DeleteAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteAttributeType(*id)
	response.SendResponse(w, "Attribute type deleted successfully", err, http.StatusInternalServerError)
}

func (h *attributeTypeHandler) UndoDeletedAttributeType(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute type ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedAttributeType(*id)
	response.SendResponse(w, "Attribute type restored successfully", err, http.StatusInternalServerError)
}

func getAllAttributeTypesData(r *http.Request, showDeleted *bool, service AttributeTypeService) (error, []AttributeTypeEntity) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	attributeTypes, err := service.GetAllAttributeTypes(showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all attribute types"), nil
	}

	return nil, attributeTypes
}

func getAttributeTypeDataById(r *http.Request, showDeleted *bool, service AttributeTypeService) (error, *AttributeTypeEntity) {
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

func validateCreateAttributeTypeRequest(req *m.CreateAttributeTypeRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}

func validateUpdateAttributeTypeRequest(req *m.UpdateAttributeTypeRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
