package auth

import (
	"strings"

	"github.com/easy-comerce/backend/pkg/constants"
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

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type CreateRoleRequest struct {
	RoleName    string   `json:"role_name"`
	RoleType    string   `json:"role_type"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type AssignRoleRequest struct {
	UserID uint     `json:"user_id"`
	Roles  []string `json:"roles"`
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

func (r *CreateRoleRequest) Sanitize() {
	if r.RoleName != "" {
		r.RoleName = utils.Trim(r.RoleName)
	}
	if r.Description != nil && *r.Description != "" {
		sanitized := utils.Trim(*r.Description)
		r.Description = &sanitized
	}
}

func (r *CreateRoleRequest) IsValidRoleType() bool {
	switch r.RoleType {
	case constants.RoleTypeAdmin, constants.RoleTypeManager, constants.RoleTypeSeller, constants.RoleTypeUser:
		return true
	default:
		return false
	}
}

func (r *CreateRoleRequest) CanCreateRoles() bool {
	return r.RoleType == constants.RoleTypeSuperAdmin || r.RoleType == constants.RoleTypeAdmin
}
