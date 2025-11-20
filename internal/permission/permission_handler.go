package permission

import (
	"net/http"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type PermissionHandler interface {
	GetAllPermissions(w http.ResponseWriter, r *http.Request)
	GetAllPermissionsGroup(w http.ResponseWriter, r *http.Request)
	GetPermissionByID(w http.ResponseWriter, r *http.Request)
}

type permissionHandler struct {
	service PermissionService
}

func NewPermissionHandler(service PermissionService) PermissionHandler {
	return &permissionHandler{
		service: service,
	}
}

func (h *permissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	permissions, err := h.service.GetAllPermissions(r.Context(), sortBy, sortOrder, nil)
	response.SendApiResponse(w, permissions, err)
}

func (h *permissionHandler) GetAllPermissionsGroup(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	permissions, err := h.service.GetAllPermissionsGroup(r.Context(), sortBy, sortOrder)
	response.SendApiResponse(w, permissions, err)
}

func (h *permissionHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.service.GetPermissionByID(r.Context(), *id)
	response.SendApiResponse(w, permission, err)
}
