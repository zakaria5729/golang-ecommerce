package auth

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/easy-comerce/backend/internal/auth/model"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
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
	var req model.SocialLoginRequest
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
	var req model.RegisterRequest
	if !utils.DecodeJSON(w, r, &req, "Register") {
		return
	}

	if validationErrors := validateRegisterRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	userResponse, err := h.service.Register(&req)
	var msg string

	if err == nil {
		msg = "Account created. We sent a verification link to your email to verify your account."
	}
	if userResponse != nil && userResponse.VerificationLink != nil {
		msg += " Verification link: " + *userResponse.VerificationLink
	}
	response.SendResponse(w, msg, err, http.StatusBadRequest)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req model.ForgotPasswordRequest
	if !utils.DecodeJSON(w, r, &req, "ForgotPassword") {
		return
	}

	if validationErrors := validateForgotPasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	token, err := h.service.ForgotPassword(&req)
	msg := "Password reset email sent"
	if token != "" && config.GetActiveProfile() != c.EnvProd {
		msg += " with token: " + token
	}
	response.SendResponse(w, msg, err, http.StatusInternalServerError)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req model.ResetPasswordRequest
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
	var req model.RefreshTokenRequest
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

func (h *AuthHandler) ResendVerifyLink(w http.ResponseWriter, r *http.Request) {
	var req model.ResendVerifyLinkRequest
	if !utils.DecodeJSON(w, r, &req, "ResendVerifyLink") {
		return
	}

	if validationErrors := validateResendVerifyLinkRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	verificationLink, verified, err := h.service.ResendVerifyLink(&req)
	if verified {
		response.SendErrorJSON(w, "Account already verified!", http.StatusBadRequest)
		return
	}

	var msg string
	if err == nil {
		msg = " We sent a verification link to your email to verify your account."
	}

	if verificationLink != "" {
		msg += " Verification link: " + verificationLink
	}
	response.SendResponse(w, msg, err, http.StatusInternalServerError)
}

func (h *AuthHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	verificationToken := r.URL.Query().Get(c.UserVerificationToken)
	if verificationToken == "" {
		renderVerificationTemplate(w, false, "Verification token is required")
		return
	}

	err := h.service.VerifyEmail(verificationToken)
	if err != nil {
		renderVerificationTemplate(w, false, "Invalid or expired verification token")
		return
	}

	renderVerificationTemplate(w, true, "Your account is verified! You can now login.")
}

func renderVerificationTemplate(w http.ResponseWriter, success bool, message string) {
	templatePath := filepath.Join(utils.GetProjectRootPath(), "templates", "account_verification.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		l.Logger.Error("❌ Failed to parse template", "method", "renderVerificationTemplate", "error", err, "templatePath", templatePath)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Success":     success,
		"Message":     message,
		"ProjectName": c.ProjectName,
	}

	w.Header().Set(c.ContentType, "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		l.Logger.Error("❌ Failed to execute template", "method", "renderVerificationTemplate", "error", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func validateRefreshTokenRequest(req *model.RefreshTokenRequest) validator.ValidationErrors {
	return validator.ValidateRequired(req.RefreshToken, "refresh_token")
}

func validateLoginRequest(req *model.LoginRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateRequired(req.Password, "password"),
	)
}

func validateSocialLoginRequest(req *model.SocialLoginRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.AuthType, "auth_type"),
		validator.ValidateRequired(req.AccessToken, "access_token"),
	)
}

func validateRegisterRequest(req *model.RegisterRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidatePassword(req.Password, "password"),
		validator.ValidateMinLength(req.Name, "name", 2),
	)
}

func validateResetPasswordRequest(req *model.ResetPasswordRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Token, "token"),
		validator.ValidatePassword(req.NewPassword, "new_password"),
	)
}

func validateForgotPasswordRequest(req *model.ForgotPasswordRequest) validator.ValidationErrors {
	return validator.ValidateRequired(req.Email, "email")
}

func validateResendVerifyLinkRequest(req *model.ResendVerifyLinkRequest) validator.ValidationErrors {
	return validator.ValidateRequired(req.Email, "email")
}
