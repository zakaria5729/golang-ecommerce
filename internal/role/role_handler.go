package role

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/role/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type RoleHandler interface {
	GetAllRoles(w http.ResponseWriter, r *http.Request)
	GetRoleByID(w http.ResponseWriter, r *http.Request)
	CreateRole(w http.ResponseWriter, r *http.Request)
	UpdateRole(w http.ResponseWriter, r *http.Request)
	DeleteRole(w http.ResponseWriter, r *http.Request)
	UndoDeletedRole(w http.ResponseWriter, r *http.Request)
	AssignRoleToUser(w http.ResponseWriter, r *http.Request)
	AppendPermissionsToRole(w http.ResponseWriter, r *http.Request)
}

type roleHandler struct {
	service RoleService
}

func NewRoleHandler(service RoleService) RoleHandler {
	return &roleHandler{
		service: service,
	}
}

func (h *roleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	include := q.Get(c.Include)
	roleType := q.Get(c.RoleRoleType)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := q.Get(c.ShowDeleted)

	roles, err := h.service.GetAllRoles(r.Context(), include, showDeleted, roleType, sortBy, sortOrder)
	response.SendResponse(w, roles, err, http.StatusInternalServerError)
}

func (h *roleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	include := q.Get(c.Include)
	showDeleted := q.Get(c.ShowDeleted)

	role, err := h.service.GetRoleByID(r.Context(), *id, include, showDeleted)
	response.SendResponse(w, role, err, http.StatusNotFound)
}

func (h *roleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRoleRequest
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

func (h *roleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var req model.UpdateRoleRequest
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

func (h *roleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteRole(r.Context(), *id)
	response.SendResponse(w, "Role deleted successfully", err, http.StatusInternalServerError)
}

func (h *roleHandler) UndoDeletedRole(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedRole(r.Context(), *id)
	response.SendResponse(w, "Undo role deleted successfully", err, http.StatusInternalServerError)
}

func (h *roleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req model.AssignRoleRequest
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

func (h *roleHandler) AppendPermissionsToRole(w http.ResponseWriter, r *http.Request) {
	var req model.AppendPermissionsToRoleRequest
	if !utils.DecodeJSON(w, r, &req, "AppendPermissionsToRole") {
		return
	}

	if validationErrors := validateAppendPermissionsToRoleRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.AppendPermissionsToRole(req.RoleID, &req)
	response.SendResponse(w, "Permissions appended to role successfully", err, http.StatusInternalServerError)
}

func validateCreateRoleRequest(req *model.CreateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
		validator.ValidateRequired(req.RoleType, "role_type"),
		validator.ValidateRequiredBool(len(req.PermissionIDs) > 0, "permission_ids"),
	)
}

func validateUpdateRoleRequest(req *model.UpdateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
		validator.ValidateRequired(req.RoleType, "role_type"),
	)
}

func validateAssignRoleRequest(req *model.AssignRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidatePositiveInteger(req.UserID, "user_id"),
		validator.ValidatePositiveInteger(req.RoleId, "role_id"),
	)
}

func validateAppendPermissionsToRoleRequest(req *model.AppendPermissionsToRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidatePositiveInteger(req.RoleID, "role_id"),
		validator.ValidateRequiredBool(len(req.PermissionIds) > 0, "permission_ids"),
	)
}
