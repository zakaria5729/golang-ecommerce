package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type PermissionHandler struct {
	permissionUseCase *permission.PermissionUseCase
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permissionUseCase: permission.NewPermissionUseCase(),
	}
}

func (h *PermissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	permissions, err := h.permissionUseCase.GetAllPermissions(sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all permissions failed", "method", "GetAllPermissions", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, permissions)
}

func (h *PermissionHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.permissionUseCase.GetPermissionByID(*id)
	if err != nil {
		logger.Logger.Error("Get permission by ID failed", "method", "GetPermissionByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, permission)
}
