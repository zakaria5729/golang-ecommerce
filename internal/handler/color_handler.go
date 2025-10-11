package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/color"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ColorHandler struct {
	service *color.ColorService
}

func NewColorHandler(service *color.ColorService) *ColorHandler {
	return &ColorHandler{
		service: service,
	}
}

func (h *ColorHandler) GetAllColorsPublic(w http.ResponseWriter, r *http.Request) {
	err, colors := getAllColorsData(r, nil, h.service)
	response.SendResponse(w, colors, err, http.StatusInternalServerError)
}

func (h *ColorHandler) GetColorByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, color := getColorDataById(r, nil, h.service)
	response.SendResponse(w, color, err, http.StatusInternalServerError)
}

func (h *ColorHandler) GetAllColors(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, colors := getAllColorsData(r, showDeleted, h.service)
	response.SendResponse(w, colors, err, http.StatusInternalServerError)
}

func (h *ColorHandler) GetColorByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, color := getColorDataById(r, showDeleted, h.service)
	response.SendResponse(w, color, err, http.StatusInternalServerError)
}

func (h *ColorHandler) CreateColor(w http.ResponseWriter, r *http.Request) {
	var req color.CreateColorRequest
	if !utils.DecodeJSON(w, r, &req, "CreateColor") {
		return
	}

	if validationErrors := validateCreateColorRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	color, err := h.service.CreateColor(&req)
	response.SendResponse(w, color, err, http.StatusInternalServerError)
}

func (h *ColorHandler) UpdateColor(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid color ID", http.StatusBadRequest)
		return
	}

	var req color.UpdateColorRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateColor") {
		return
	}

	if validationErrors := validateUpdateColorRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	color, err := h.service.UpdateColor(*id, &req)
	response.SendResponse(w, color, err, http.StatusInternalServerError)
}

func (h *ColorHandler) DeleteColor(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid color ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteColor(*id)
	response.SendResponse(w, "Color deleted successfully", err, http.StatusInternalServerError)
}

func (h *ColorHandler) UndoDeletedColor(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid color ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedColor(*id)
	response.SendResponse(w, "Color restored successfully", err, http.StatusInternalServerError)
}

func getAllColorsData(r *http.Request, showDeleted *bool, service *color.ColorService) (error, []color.Color) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	colors, err := service.GetAllColors(showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all colors"), nil
	}

	return nil, colors
}

func getColorDataById(r *http.Request, showDeleted *bool, service *color.ColorService) (error, *color.Color) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid color ID"), nil
	}

	color, err := service.GetColorByID(*id, showDeleted)
	if err != nil {
		return errors.New("color not found"), nil
	}

	return nil, color
}

func validateCreateColorRequest(req *color.CreateColorRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}

func validateUpdateColorRequest(req *color.UpdateColorRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
