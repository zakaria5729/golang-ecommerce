package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/role"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type RoleHandler struct {
	service *role.RoleService
}

func NewRoleHandler(service *role.RoleService) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	roleType := q.Get(c.RoleRoleType)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := q.Get(c.ShowDeleted)

	roles, err := h.service.GetAllRoles(include, showDeleted, roleType, sortBy, sortOrder)
	response.SendResponse(w, roles, err, http.StatusInternalServerError)
}

func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	include := q.Get(c.Include)
	showDeleted := q.Get(c.ShowDeleted)

	role, err := h.service.GetRoleByID(*id, include, showDeleted)
	response.SendResponse(w, role, err, http.StatusNotFound)
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req role.CreateRoleRequest
	if !utils.DecodeJSON(w, r, &req, "CreateRole") {
		return
	}

	if validationErrors := validateCreateRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	role, err := h.service.CreateRole(r.Context(), &req)
	response.SendResponse(w, role, err, http.StatusBadRequest)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var req role.UpdateRoleRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateRole") {
		return
	}

	if validationErrors := validateUpdateRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	role, err := h.service.UpdateRole(r.Context(), *id, &req)
	response.SendResponse(w, role, err, http.StatusBadRequest)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteRole(r.Context(), *id)
	response.SendResponse(w, "Role deleted successfully", err, http.StatusInternalServerError)
}

func (h *RoleHandler) UndoDeletedRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedRole(r.Context(), *id)
	response.SendResponse(w, "Undo role deleted successfully", err, http.StatusInternalServerError)
}

func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req role.AssignRoleRequest
	if !utils.DecodeJSON(w, r, &req, "AssignRoleToUser") {
		return
	}

	if validationErrors := validateAssignRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.AssignRoleToUser(req.UserID, &req)
	response.SendResponse(w, "Roles assigned successfully", err, http.StatusInternalServerError)
}

func (h *RoleHandler) AddPermissionsToRole(w http.ResponseWriter, r *http.Request) {
	var req role.AddPermissionsToRoleRequest
	if !utils.DecodeJSON(w, r, &req, "AddPermissionsToRole") {
		return
	}

	if validationErrors := validateAddPermissionsToRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.AddPermissionsToRole(req.RoleID, &req)
	response.SendResponse(w, "Permissions added to role successfully", err, http.StatusInternalServerError)
}

func validateCreateRoleRequest(req *role.CreateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
		validator.ValidateRequired(req.RoleType, "role_type"),
		validator.ValidateRequiredBool(len(req.PermissionIDs) > 0, "permission_ids"),
	)
}

func validateUpdateRoleRequest(req *role.UpdateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
		validator.ValidateRequired(req.RoleType, "role_type"),
	)
}

func validateAssignRoleRequest(req *role.AssignRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidatePositiveInteger(req.UserID, "user_id"),
		validator.ValidatePositiveInteger(req.RoleId, "role_id"),
	)
}

func validateAddPermissionsToRoleRequest(req *role.AddPermissionsToRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidatePositiveInteger(req.RoleID, "role_id"),
		validator.ValidateRequiredBool(len(req.PermissionIds) > 0, "permission_ids"),
	)
}
