package size_option

import (
	"errors"
	"net/http"

	m "github.com/easy-comerce/backend/internal/size_option/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type SizeOptionHandler interface {
	GetAllSizeOptionsPublic(w http.ResponseWriter, r *http.Request)
	GetSizeOptionByIDPublic(w http.ResponseWriter, r *http.Request)
	GetAllSizeOptions(w http.ResponseWriter, r *http.Request)
	GetSizeOptionByID(w http.ResponseWriter, r *http.Request)
	CreateSizeOption(w http.ResponseWriter, r *http.Request)
	UpdateSizeOption(w http.ResponseWriter, r *http.Request)
	DeleteSizeOption(w http.ResponseWriter, r *http.Request)
	UndoDeletedSizeOption(w http.ResponseWriter, r *http.Request)
}

type sizeOptionHandler struct {
	service SizeOptionService
}

func NewSizeOptionHandler(service SizeOptionService) SizeOptionHandler {
	return &sizeOptionHandler{
		service: service,
	}
}

func (h *sizeOptionHandler) GetAllSizeOptionsPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeOptions := getAllSizeOptionsData(r, nil, h.service)
	response.SendResponse(w, sizeOptions, err, http.StatusInternalServerError)
}

func (h *sizeOptionHandler) GetSizeOptionByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, sizeOption := getSizeOptionDataById(r, nil, h.service)
	response.SendResponse(w, sizeOption, err, http.StatusInternalServerError)
}

func (h *sizeOptionHandler) GetAllSizeOptions(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeOptions := getAllSizeOptionsData(r, showDeleted, h.service)
	response.SendResponse(w, sizeOptions, err, http.StatusInternalServerError)
}

func (h *sizeOptionHandler) GetSizeOptionByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, sizeOption := getSizeOptionDataById(r, showDeleted, h.service)
	response.SendResponse(w, sizeOption, err, http.StatusInternalServerError)
}

func (h *sizeOptionHandler) CreateSizeOption(w http.ResponseWriter, r *http.Request) {
	var req m.CreateSizeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "CreateSizeOption") {
		return
	}

	if validationErrors := validateCreateSizeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeOption, err := h.service.CreateSizeOption(&req)
	response.SendResponse(w, sizeOption, err, http.StatusInternalServerError)
}

func (h *sizeOptionHandler) UpdateSizeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size option ID", http.StatusBadRequest)
		return
	}

	var req m.UpdateSizeOptionRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateSizeOption") {
		return
	}

	if validationErrors := validateUpdateSizeOptionRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	sizeOption, err := h.service.UpdateSizeOption(*id, &req)
	response.SendResponse(w, sizeOption, err, http.StatusInternalServerError)
}

func (h *sizeOptionHandler) DeleteSizeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size option ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSizeOption(*id)
	response.SendResponse(w, "Size option deleted successfully", err, http.StatusInternalServerError)
}

func (h *sizeOptionHandler) UndoDeletedSizeOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid size option ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedSizeOption(*id)
	response.SendResponse(w, "Size option restored successfully", err, http.StatusInternalServerError)
}

func getAllSizeOptionsData(r *http.Request, showDeleted *bool, service SizeOptionService) (error, []SizeOptionEntity) {
	q := r.URL.Query()
	sizeCategoryID := q.Get(c.SizeOptionSizeCategoryID)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	sizeOptions, err := service.GetAllSizeOptions(showDeleted, sizeCategoryID, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all size options"), nil
	}

	return nil, sizeOptions
}

func getSizeOptionDataById(r *http.Request, showDeleted *bool, service SizeOptionService) (error, *SizeOptionEntity) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid size option ID"), nil
	}

	sizeOption, err := service.GetSizeOptionByID(*id, showDeleted)
	if err != nil {
		return errors.New("failed to get size option by ID"), nil
	}

	return nil, sizeOption
}

func validateCreateSizeOptionRequest(req *m.CreateSizeOptionRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
		validator.ValidatePositiveInteger(req.SizeCategoryID, "size_category_id"),
	)
}

func validateUpdateSizeOptionRequest(req *m.UpdateSizeOptionRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
