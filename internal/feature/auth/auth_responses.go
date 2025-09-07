package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	User         *User  `json:"user"`
}

type Claims struct {
	UserID   uint     `json:"user_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Username string   `json:"username"`
	jwt.RegisteredClaims
}
