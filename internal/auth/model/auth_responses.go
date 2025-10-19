package model

import (
	"github.com/easy-comerce/backend/internal/user/model"
)

type LoginResponse struct {
	ExpiresAt    int64               `json:"expires_at"`
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         *model.UserResponse `json:"user,omitempty"`
}
