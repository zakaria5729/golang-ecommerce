package auth

import (
	"github.com/easy-comerce/backend/internal/feature/user"
)

type LoginResponse struct {
	ExpiresAt    int64              `json:"expires_at"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	User         *user.UserResponse `json:"user,omitempty"`
}
