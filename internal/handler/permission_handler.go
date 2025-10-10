package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type PermissionHandler struct {
	service *permission.PermissionService
}

func NewPermissionHandler(service *permission.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		service: service,
	}
}

func (h *PermissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	permissions, err := h.service.GetAllPermissions(sortBy, sortOrder)
	response.SendResponse(w, permissions, err, http.StatusInternalServerError)
}

func (h *PermissionHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.service.GetPermissionByID(*id)
	response.SendResponse(w, permission, err, http.StatusNotFound)
}
