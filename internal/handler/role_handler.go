package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/pkg/constants"
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
	include := q.Get(constants.Include)
	roleType := q.Get(constants.RoleRoleType)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)
	showDeleted := q.Get(constants.ShowDeleted)

	roles, err := h.service.GetAllRoles(include, showDeleted, roleType, sortBy, sortOrder)
	response.SendResponse(w, roles, err, http.StatusInternalServerError)
}

func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	include := q.Get(constants.Include)
	showDeleted := q.Get(constants.ShowDeleted)

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
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
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
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteRole(r.Context(), *id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Role deleted successfully")
}

func (h *RoleHandler) UndoDeletedRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.service.UndoDeletedRole(r.Context(), *id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Undo role deleted successfully")
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

	if err := h.service.AssignRoleToUser(req.UserID, &req); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendCommonResponseJSON(w, "Roles assigned successfully")
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

	if err := h.service.AddPermissionsToRole(req.RoleID, &req); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendCommonResponseJSON(w, "Permissions added to role successfully")
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
