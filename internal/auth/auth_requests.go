package auth

import (
	"strings"

	"github.com/easy-comerce/backend/pkg/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResendVerifyLinkRequest struct {
	Email string `json:"email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SocialLoginRequest struct {
	IdToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	AuthType    string `json:"auth_type"`
}

func (r *LoginRequest) Sanitize() {
	if r.Email != "" {
		r.Email = utils.Trim(strings.ToLower(r.Email))
	}
}

func (r *RegisterRequest) Sanitize() {
	if r.Email != "" {
		r.Email = utils.Trim(strings.ToLower(r.Email))
	}
	if r.Name != "" {
		r.Name = utils.Trim(r.Name)
	}
}

func (r *ForgotPasswordRequest) Sanitize() {
	if r.Email != "" {
		r.Email = utils.Trim(strings.ToLower(r.Email))
	}
}
