package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AuthHandler struct {
	authUseCase *auth.AuthUseCase
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authUseCase: auth.NewAuthUseCase(jwtSecret),
	}
}

// **REQUIRED
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if !utils.DecodeJSON(w, r, &req, "Login") {
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

// **REQUIRED
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if !utils.DecodeJSON(w, r, &req, "Register") {
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

// **REQUIRED
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req auth.ChangePasswordRequest
	if !utils.DecodeJSON(w, r, &req, "ChangePassword") {
		return
	}

	if validationErrors := h.validateChangePasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err := h.authUseCase.ChangePassword(*user, &req); err != nil {
		logger.Logger.Error("Change password failed", "method", "ChangePassword", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendCommonResponseJSON(w, "Password changed successfully")
}

// **REQUIRED
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ForgotPasswordRequest
	if !utils.DecodeJSON(w, r, &req, "ForgotPassword") {
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

	response.SendCommonResponseJSON(w, "Password reset email sent")
}

// **REQUIRED
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ResetPasswordRequest
	if !utils.DecodeJSON(w, r, &req, "ResetPassword") {
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

	response.SendCommonResponseJSON(w, "Password reset successfully")
}

// **REQUIRED
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshTokenRequest
	if !utils.DecodeJSON(w, r, &req, "RefreshToken") {
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

// **REQUIRED
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.authUseCase.Logout(*id); err != nil {
		logger.Logger.Error("Logout failed", "method", "Logout", "error", err, "userID", id)
		response.SendErrorJSON(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	response.SendCommonResponseJSON(w, "Logged out successfully")
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
