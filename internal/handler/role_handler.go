package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type RoleHandler struct {
	roleUseCase *role.RoleUseCase
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleUseCase: role.NewRoleUseCase(),
	}
}

// **REQUIRED
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	include := q.Get(constants.Include)
	roleType := q.Get(constants.RoleRoleType)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	roles, err := h.roleUseCase.GetAllRoles(include, roleType, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all roles failed", "method", "GetAllRoles", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, roles)
}

// **REQUIRED
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	role, err := h.roleUseCase.GetRoleByID(*id, include)
	if err != nil {
		logger.Logger.Error("Get role by ID failed", "method", "GetRoleByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, role)
}

// **REQUIRED
func (h *RoleHandler) AddToWishlistsByUser(w http.ResponseWriter, r *http.Request) {
	var req role.CreateRoleRequest
	if !utils.DecodeJSON(w, r, &req, "CreateRole") {
		return
	}

	if validationErrors := h.validateCreateRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	role, err := h.roleUseCase.CreateRole(&req)
	if err != nil {
		logger.Logger.Error("Create role failed", "method", "CreateRole", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, role, http.StatusCreated)
}

// **REQUIRED
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

	if validationErrors := h.validateUpdateRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	role, err := h.roleUseCase.UpdateRole(*id, &req)
	if err != nil {
		logger.Logger.Error("Update role failed", "method", "UpdateRole", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, role)
}

// **REQUIRED
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.roleUseCase.DeleteRole(*id); err != nil {
		logger.Logger.Error("Delete role failed", "method", "DeleteRole", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Role deleted successfully")
}

// **REQUIRED
func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req role.AssignRoleRequest
	if !utils.DecodeJSON(w, r, &req, "AssignRoleToUser") {
		return
	}

	if validationErrors := h.validateAssignRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err := h.roleUseCase.AssignRoleToUser(req.UserID, &req); err != nil {
		logger.Logger.Error("Assign roles to user failed", "method", "AssignRoleToUser", "error", err, "userID", req.UserID)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendCommonResponseJSON(w, "Roles assigned successfully")
}

func (h *RoleHandler) validateCreateRoleRequest(req *role.CreateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
		validator.ValidateRequired(req.RoleType, "role_type"),
		validator.ValidateRequiredBool(len(req.PermissionIDs) > 0, "permission_ids"),
	)
}

func (h *RoleHandler) validateUpdateRoleRequest(req *role.UpdateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
		validator.ValidateRequired(req.RoleType, "role_type"),
	)
}

func (h *RoleHandler) validateAssignRoleRequest(req *role.AssignRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidatePositiveInteger(req.UserID, "user_id"),
		validator.ValidatePositiveInteger(req.RoleId, "role_id"),
	)
}
