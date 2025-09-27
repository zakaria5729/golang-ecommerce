package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/color"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ColorHandler struct {
	colorUseCase *color.ColorUseCase
}

func NewColorHandler() *ColorHandler {
	return &ColorHandler{
		colorUseCase: color.NewColorUseCase(),
	}
}

func (h *ColorHandler) GetAllColorsPublic(w http.ResponseWriter, r *http.Request) {
	err, colors := h.getAllColorsData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, colors)
}

func (h *ColorHandler) GetColorByIDPublic(w http.ResponseWriter, r *http.Request) {
	err, color := h.getColorDataById(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, color)
}

func (h *ColorHandler) GetAllColors(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, colors := h.getAllColorsData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, colors)
}

func (h *ColorHandler) GetColorByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, color := h.getColorDataById(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, color)
}

func (h *ColorHandler) CreateColor(w http.ResponseWriter, r *http.Request) {
	var req color.CreateColorRequest
	if !utils.DecodeJSON(w, r, &req, "CreateColor") {
		return
	}

	if validationErrors := h.validateCreateColorRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	color, err := h.colorUseCase.CreateColor(&req)
	if err != nil {
		logger.Logger.Error("Create color failed", "method", "CreateColor", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, color, http.StatusCreated)
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

	if validationErrors := h.validateUpdateColorRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	color, err := h.colorUseCase.UpdateColor(*id, &req)
	if err != nil {
		logger.Logger.Error("Update color failed", "method", "UpdateColor", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, color)
}

func (h *ColorHandler) DeleteColor(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid color ID", http.StatusBadRequest)
		return
	}

	if err := h.colorUseCase.DeleteColor(*id); err != nil {
		logger.Logger.Error("Delete color failed", "method", "DeleteColor", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Color deleted successfully")
}

func (h *ColorHandler) UndoDeletedColor(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid color ID", http.StatusBadRequest)
		return
	}

	if err := h.colorUseCase.UndoDeletedColor(*id); err != nil {
		logger.Logger.Error("Undo deleted color failed", "method", "UndoDeletedColor", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Color restored successfully")
}

func (h *ColorHandler) getAllColorsData(r *http.Request, showDeleted *bool) (error, []color.Color) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	colors, err := h.colorUseCase.GetAllColors(include, showDeleted, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to get all colors"), nil
	}

	return nil, colors
}

func (h *ColorHandler) getColorDataById(r *http.Request, showDeleted *bool) (error, *color.Color) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid color ID"), nil
	}

	include := r.URL.Query().Get(c.Include)
	color, err := h.colorUseCase.GetColorByID(*id, include, showDeleted)
	if err != nil {
		return errors.New("color not found"), nil
	}

	return nil, color
}

func (h *ColorHandler) validateCreateColorRequest(req *color.CreateColorRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}

func (h *ColorHandler) validateUpdateColorRequest(req *color.UpdateColorRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateMaxLength(req.Name, "name", 100),
	)
}
