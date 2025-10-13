package permission

import (
	"net/http"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type PermissionHandler struct {
	service *PermissionService
}

func NewPermissionHandler(service *PermissionService) *PermissionHandler {
	return &PermissionHandler{
		service: service,
	}
}

func (h *PermissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	permissions, err := h.service.GetAllPermissions(sortBy, sortOrder)
	response.SendResponse(w, permissions, err, http.StatusInternalServerError)
}

func (h *PermissionHandler) GetAllPermissionsGroup(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	permissions, err := h.service.GetAllPermissionsGroup(sortBy, sortOrder)
	response.SendResponse(w, permissions, err, http.StatusInternalServerError)
}

func (h *PermissionHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.service.GetPermissionByID(*id)
	response.SendResponse(w, permission, err, http.StatusNotFound)
}
