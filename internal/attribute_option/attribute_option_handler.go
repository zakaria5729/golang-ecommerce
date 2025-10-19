package attribute_option

import (
	"errors"
	"net/http"

	m "github.com/easy-comerce/backend/internal/attribute_option/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AttributeOptionHandler struct {
	service *AttributeOptionService
}

func NewAttributeOptionHandler(service *AttributeOptionService) *AttributeOptionHandler {
	return &AttributeOptionHandler{
		service: service,
	}
}

func (h *AttributeOptionHandler) GetAllAttributeOptionsPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeOptions := getAllAttributeOptionsData(r, nil, h.service)
	response.SendResponse(w, attributeOptions, err, http.StatusInternalServerError)
}

func (h *AttributeOptionHandler) GetAttributeOptionByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, attributeOption := getAllAttributeOptionsData(r, nil, h.service)
	response.SendResponse(w, attributeOption, err, http.StatusInternalServerError)
}

func (h *AttributeOptionHandler) GetAllAttributeOptions(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeOptions := getAllAttributeOptionsData(r, showDeleted, h.service)
	response.SendResponse(w, attributeOptions, err, http.StatusInternalServerError)
}

func (h *AttributeOptionHandler) GetAttributeOptionByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, attributeOption := getAllAttributeOptionsData(r, showDeleted, h.service)
	response.SendResponse(w, attributeOption, err, http.StatusInternalServerError)
}

func (h *AttributeOptionHandler) CreateAttributeOption(w http.ResponseWriter, r *http.Request) {
	var req m.CreateAttributeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "CreateAttributeOption") {
		return
	}

	if validationErrors := validateCreateAttributeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeOption, err := h.service.CreateAttributeOption(&req)
	response.SendResponse(w, attributeOption, err, http.StatusInternalServerError)
}

func (h *AttributeOptionHandler) UpdateAttributeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute option ID", http.StatusBadRequest)
		return
	}

	var req m.UpdateAttributeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateAttributeOption") {
		return
	}

	if validationErrors := validateUpdateAttributeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	attributeOption, err := h.service.UpdateAttributeOption(*id, &req)
	response.SendResponse(w, attributeOption, err, http.StatusInternalServerError)
}

func (h *AttributeOptionHandler) DeleteAttributeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute option ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteAttributeOption(*id)
	response.SendResponse(w, "Attribute option deleted successfully", err, http.StatusInternalServerError)
}

func (h *AttributeOptionHandler) UndoDeletedAttributeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid attribute option ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedAttributeOption(*id)
	response.SendResponse(w, "Attribute option restored successfully", err, http.StatusInternalServerError)
}

func getAllAttributeOptionsData(r *http.Request, showDeleted *bool, service *AttributeOptionService) (error, []AttributeOptionEntity) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	attributeTypeID := q.Get("attribute_type_id")
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	attributeOptions, err := service.GetAllAttributeOptions(include, showDeleted, attributeTypeID, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all attribute options"), nil
	}

	return nil, attributeOptions
}

func validateCreateAttributeOptionRequest(req *m.CreateAttributeOptionRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidatePositiveInteger(req.AttributeTypeID, "attribute_type_id"),
		validator.ValidateRequired(req.AttributeOptionName, "attribute_option_name"),
		validator.ValidateMinLength(req.AttributeOptionName, "attribute_option_name", 2),
		validator.ValidateMaxLength(req.AttributeOptionName, "attribute_option_name", 100),
	)
}

func validateUpdateAttributeOptionRequest(req *m.UpdateAttributeOptionRequest) validator.ValidationErrors {
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
