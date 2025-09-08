package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/validator"
)

type RoleHandler struct {
	roleUseCase *role.RoleUseCase
}

func NewRoleHandler() *RoleHandler {
	userRepo := user.NewUserRepository()
	return &RoleHandler{
		roleUseCase: role.NewRoleUseCase(userRepo.ToSharedInterface()),
	}
}

func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	include := r.URL.Query().Get("include")
	roleType := r.URL.Query().Get("role_type")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	roles, err := h.roleUseCase.GetAllRoles(include, roleType, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all roles failed", "method", "GetAllRoles", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, roles)
}

func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.SendErrorJSON(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get("include")
	role, err := h.roleUseCase.GetRoleByID(uint(id), include)
	if err != nil {
		logger.Logger.Error("Get role by ID failed", "method", "GetRoleByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, role)
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req role.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode create role request", "method", "CreateRole", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
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

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.SendErrorJSON(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var req role.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode update role request", "method", "UpdateRole", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateCreateRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	role, err := h.roleUseCase.UpdateRole(uint(id), &req)
	if err != nil {
		logger.Logger.Error("Update role failed", "method", "UpdateRole", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, role)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.SendErrorJSON(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.roleUseCase.DeleteRole(uint(id)); err != nil {
		logger.Logger.Error("Delete role failed", "method", "DeleteRole", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Role deleted successfully")
}

func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req role.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode assign role request", "method", "AssignRoleToUser", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateAssignRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err := h.roleUseCase.AssignRoleToUser(req.UserID, &req); err != nil {
		logger.Logger.Error("Assign role to user failed", "method", "AssignRoleToUser", "error", err, "userID", req.UserID)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Roles assigned successfully"})
}

func (h *RoleHandler) validateCreateRoleRequest(req *role.CreateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
	)
}

func (h *RoleHandler) validateAssignRoleRequest(req *role.AssignRoleRequest) validator.ValidationErrors {
	var errors validator.ValidationErrors
	if req.UserID == 0 {
		errors.AddError("user_id", "user_id is required")
	}
	if len(req.Roles) == 0 {
		errors.AddError("roles", "at least one role is required")
	}
	return errors
}
