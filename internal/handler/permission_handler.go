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

// **REQUIRED
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

// **REQUIRED
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

// func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
// 	var req permission.Permission
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		logger.Logger.Error("Failed to decode create permission request", "method", "CreatePermission", "error", err)
// 		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if validationErrors := h.validateCreatePermissionRequest(&req); len(validationErrors) > 0 {
// 		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
// 		return
// 	}

// 	permission, err := h.permissionUseCase.CreatePermission(&req)
// 	if err != nil {
// 		logger.Logger.Error("Create permission failed", "method", "CreatePermission", "error", err)
// 		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	response.SendSuccessJSON(w, permission, http.StatusCreated)
// }

// func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	if idStr == "" {
// 		response.SendErrorJSON(w, "Permission ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		response.SendErrorJSON(w, "Invalid permission ID", http.StatusBadRequest)
// 		return
// 	}

// 	var req permission.Permission
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		logger.Logger.Error("Failed to decode update permission request", "method", "UpdatePermission", "error", err)
// 		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if validationErrors := h.validateCreatePermissionRequest(&req); len(validationErrors) > 0 {
// 		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
// 		return
// 	}

// 	permission, err := h.permissionUseCase.UpdatePermission(uint(id), &req)
// 	if err != nil {
// 		logger.Logger.Error("Update permission failed", "method", "UpdatePermission", "error", err, "id", id)
// 		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	response.SendSuccessJSON(w, permission)
// }

// func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	if idStr == "" {
// 		response.SendErrorJSON(w, "Permission ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		response.SendErrorJSON(w, "Invalid permission ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.permissionUseCase.DeletePermission(uint(id)); err != nil {
// 		logger.Logger.Error("Delete permission failed", "method", "DeletePermission", "error", err, "id", id)
// 		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	response.SendDeleteJSON(w, "Permission deleted successfully")
// }

// func (h *PermissionHandler) validateCreatePermissionRequest(req *permission.Permission) validator.ValidationErrors {
// 	return validator.MergeValidationErrors(
// 		validator.ValidateRequired(req.Name, "name"),
// 		validator.ValidateMinLength(req.Name, "name", 2),
// 	)
// }
