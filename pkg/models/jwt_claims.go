package models

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	UserID   uint     `json:"user_id"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	jwt.RegisteredClaims
}
