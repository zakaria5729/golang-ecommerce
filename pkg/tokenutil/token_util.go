package tokenutil

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	userEntity "github.com/easy-comerce/backend/internal/user"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserID   uint     `json:"user_id"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	jwt.RegisteredClaims
}

func GenerateNewToken(isFallback ...bool) (string, error) {
	ts := time.Now().UnixNano()
	randBytes := make([]byte, 6)
	if _, err := rand.Read(randBytes); err != nil {
		if len(isFallback) > 0 && isFallback[0] {
			return timeutil.NowUTC().Format("20060102150405"), nil
		}
		return "", err
	}
	return fmt.Sprintf("%x%x", ts, randBytes), nil
}

func GenerateNewTokenWithExpiryTime(expiryHours int) (string, time.Time, error) {
	token, err := GenerateNewToken(true)
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := timeutil.AddHoursUTC(expiryHours)
	return token, expiresAt, nil
}

func VerifyJwtToken(tokenString string, jwtSecret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			l.Logger.Error("❌ unexpected signing method: %v", token.Header["alg"], "method", "VerifyJwtToken")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if ok && token.Valid {
		return claims, nil
	}

	if claims.ExpiresAt.Before(timeutil.NowUTC()) {
		l.Logger.Error("❌ Jwt Token expired", "method", "VerifyJwtToken", "error", err)
		return nil, errors.New("jwt token expired")
	}

	return nil, errors.New("jwt token invalid")
}

func GenerateNewJwtToken(user *userEntity.UserEntity, jwtSecret string) (string, int64, error) {
	now := timeutil.NowUTC()
	expirationTime := timeutil.AddHoursUTC(c.AccessTokenExpiryHours)
	expiresAt := expirationTime.Unix()

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.RoleType)
	}

	claims := &JwtClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Roles:    roleNames,
		Username: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    c.ProjectName,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

func ExtractJwtToken(r *http.Request) string {
	authHeader := r.Header.Get(c.Authorization)
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != c.Bearer {
		return ""
	}
	return parts[1]
}

func ValidateTokenAndGetJwtClaims(r *http.Request, jwtSecret string) (*JwtClaims, error) {
	token := ExtractJwtToken(r)
	if token == "" {
		return nil, errors.New("no jwt token provided")
	}

	claims, err := VerifyJwtToken(token, jwtSecret)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
