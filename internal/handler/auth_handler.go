package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/auth"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AuthHandler struct {
	service *auth.AuthService
}

func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) AppHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req auth.AppHealthRequest
	if !utils.DecodeJSON(w, r, &req, "AppHealthCheck") {
		return
	}

	if req.HealthToken != nil && *req.HealthToken == c.AppHealthCheckToken {
		healthResponse := h.service.HealthCheck(r.Context())
		response.SendResponse(w, healthResponse, nil, http.StatusOK)
		return
	}

	response.SendErrorJSON(w, "Invalid health token", http.StatusForbidden)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if !utils.DecodeJSON(w, r, &req, "Login") {
		return
	}

	if validationErrors := validateLoginRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	loginResponse, err := h.service.Login(&req)
	response.SendResponse(w, loginResponse, err, http.StatusOK)
}

func (h *AuthHandler) SocialLogin(w http.ResponseWriter, r *http.Request) {
	var req auth.SocialLoginRequest
	if !utils.DecodeJSON(w, r, &req, "SocialLogin") {
		return
	}

	if validationErrors := validateSocialLoginRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	loginResponse, err := h.service.SocialLogin(&req)
	response.SendResponse(w, loginResponse, err, http.StatusOK)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if !utils.DecodeJSON(w, r, &req, "Register") {
		return
	}

	if validationErrors := validateRegisterRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	user, err := h.service.Register(&req)
	response.SendResponse(w, user, err, http.StatusCreated)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ForgotPasswordRequest
	if !utils.DecodeJSON(w, r, &req, "ForgotPassword") {
		return
	}

	if validationErrors := validateForgotPasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.ForgotPassword(&req)
	response.SendResponse(w, "Password reset email sent", err, http.StatusInternalServerError)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ResetPasswordRequest
	if !utils.DecodeJSON(w, r, &req, "ResetPassword") {
		return
	}

	if validationErrors := validateResetPasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.ResetPassword(&req)
	response.SendResponse(w, "Password reset successfully", err, http.StatusInternalServerError)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshTokenRequest
	if !utils.DecodeJSON(w, r, &req, "RefreshToken") {
		return
	}

	if validationErrors := validateRefreshTokenRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	loginResponse, err := h.service.RefreshToken(&req)
	response.SendResponse(w, loginResponse, err, http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUint(r.PathValue(c.FieldUserID))
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.service.Logout(*userID)
	response.SendResponse(w, "Logged out successfully", err, http.StatusInternalServerError)
}

func validateRefreshTokenRequest(req *auth.RefreshTokenRequest) validator.ValidationErrors {
	return validator.ValidateRequired(req.RefreshToken, "refresh_token")
}

func validateLoginRequest(req *auth.LoginRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateRequired(req.Password, "password"),
	)
}

func validateSocialLoginRequest(req *auth.SocialLoginRequest) validator.ValidationErrors {
	var validationErrors validator.ValidationErrors
	switch req.AuthType {
	case c.AuthTypeGoogle:
		validationErrors = validator.ValidateRequired(req.IdToken, "id_token")
	case c.AuthTypeFacebook:
		validationErrors = validator.ValidateRequired(req.AccessToken, "access_token")
	}
	return validationErrors
}

func validateRegisterRequest(req *auth.RegisterRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidatePassword(req.Password, "password"),
		validator.ValidateMinLength(req.Name, "name", 2),
	)
}

func validateResetPasswordRequest(req *auth.ResetPasswordRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Token, "token"),
		validator.ValidatePassword(req.NewPassword, "new_password"),
	)
}

func validateForgotPasswordRequest(req *auth.ForgotPasswordRequest) validator.ValidationErrors {
	return validator.ValidateRequired(req.Email, "email")
}
