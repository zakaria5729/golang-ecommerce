package auth

import (
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/golang-jwt/jwt/v5"
)

type LoginResponse struct {
	ExpiresAt    int64              `json:"expires_at"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	User         *user.UserResponse `json:"user,omitempty"`
}

type JwtClaims struct {
	UserID   uint     `json:"user_id"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	jwt.RegisteredClaims
}
