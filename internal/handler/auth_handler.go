package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AuthHandler struct {
	authUseCase *auth.AuthUseCase
	roleUseCase *auth.RoleUseCase
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authUseCase: auth.NewAuthUseCase(jwtSecret),
		roleUseCase: auth.NewRoleUseCase(),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode login request", "method", "Login", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateLoginRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	loginResponse, err := h.authUseCase.Login(&req)
	if err != nil {
		logger.Logger.Error("Login failed", "method", "Login", "error", err, "email", req.Email)
		response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.SendSuccessJSON(w, loginResponse)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode register request", "method", "Register", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateRegisterRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	user, err := h.authUseCase.Register(&req)
	if err != nil {
		logger.Logger.Error("Registration failed", "method", "Register", "error", err, "email", req.Email)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, user, http.StatusCreated)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req auth.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode change password request", "method", "ChangePassword", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateChangePasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err := h.authUseCase.ChangePassword(user.ID, &req); err != nil {
		logger.Logger.Error("Change password failed", "method", "ChangePassword", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Password changed successfully"})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode forgot password request", "method", "ForgotPassword", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateForgotPasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err := h.authUseCase.ForgotPassword(&req); err != nil {
		logger.Logger.Error("Forgot password failed", "method", "ForgotPassword", "error", err, "email", req.Email)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Password reset email sent"})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode reset password request", "method", "ResetPassword", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateResetPasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err := h.authUseCase.ResetPassword(&req); err != nil {
		logger.Logger.Error("Reset password failed", "method", "ResetPassword", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Password reset successfully"})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	include := r.URL.Query().Get("include")
	profile, err := h.authUseCase.GetUserProfile(user.ID, include)
	if err != nil {
		logger.Logger.Error("Get profile failed", "method", "GetProfile", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, profile)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	include := r.URL.Query().Get("include")
	profile, err := h.authUseCase.GetUserProfile(user.ID, include)
	if err != nil {
		logger.Logger.Error("Get me failed", "method", "GetMe", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, profile)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req auth.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode update profile request", "method", "UpdateProfile", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.authUseCase.UpdateUserProfile(user.ID, &req)
	if err != nil {
		logger.Logger.Error("Update profile failed", "method", "UpdateProfile", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, updatedUser)
}

func (h *AuthHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

func (h *AuthHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	include := r.URL.Query().Get("include")
	role, err := h.roleUseCase.GetRoleByID(uint(id), include)
	if err != nil {
		logger.Logger.Error("Get role by ID failed", "method", "GetRoleByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, role)
}

func (h *AuthHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.CreateRoleRequest
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

func (h *AuthHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
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

	var req auth.CreateRoleRequest
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

func (h *AuthHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.AssignRoleRequest
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

func (h *AuthHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	include := r.URL.Query().Get("include")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	permissions, err := h.roleUseCase.GetAllPermissions(include, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all permissions failed", "method", "GetAllPermissions", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, permissions)
}

func (h *AuthHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.SendErrorJSON(w, "Permission ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendErrorJSON(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get("include")
	permission, err := h.roleUseCase.GetPermissionByID(uint(id), include)
	if err != nil {
		logger.Logger.Error("Get permission by ID failed", "method", "GetPermissionByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, permission)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode refresh token request", "method", "RefreshToken", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateRefreshTokenRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	loginResponse, err := h.authUseCase.RefreshToken(&req)
	if err != nil {
		logger.Logger.Error("Refresh token failed", "method", "RefreshToken", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.SendSuccessJSON(w, loginResponse)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if err := h.authUseCase.Logout(user.ID); err != nil {
		logger.Logger.Error("Logout failed", "method", "Logout", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) validateRefreshTokenRequest(req *auth.RefreshTokenRequest) validator.ValidationErrors {
	return validator.ValidateRequired(req.RefreshToken, "refresh_token")
}

func (h *AuthHandler) validateLoginRequest(req *auth.LoginRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateRequired(req.Password, "password"),
	)
}

func (h *AuthHandler) validateRegisterRequest(req *auth.RegisterRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidatePassword(req.Password, "password"),
		validator.ValidateMinLength(req.Name, "name", 2),
	)
}

func (h *AuthHandler) validateChangePasswordRequest(req *auth.ChangePasswordRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.CurrentPassword, "current_password"),
		validator.ValidatePassword(req.NewPassword, "new_password"),
	)
}

func (h *AuthHandler) validateResetPasswordRequest(req *auth.ResetPasswordRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Token, "token"),
		validator.ValidatePassword(req.NewPassword, "new_password"),
	)
}

func (h *AuthHandler) validateForgotPasswordRequest(req *auth.ForgotPasswordRequest) validator.ValidationErrors {
	return validator.ValidateRequired(req.Email, "email")
}

func (h *AuthHandler) validateCreateRoleRequest(req *auth.CreateRoleRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.RoleName, "role_name"),
		validator.ValidateMinLength(req.RoleName, "role_name", 2),
	)
}

func (h *AuthHandler) validateAssignRoleRequest(req *auth.AssignRoleRequest) validator.ValidationErrors {
	var errors validator.ValidationErrors
	if req.UserID == 0 {
		errors.AddError("user_id", "user_id is required")
	}
	if len(req.Roles) == 0 {
		errors.AddError("roles", "at least one role is required")
	}
	return errors
}
