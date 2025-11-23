package tokenutil

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	e "github.com/easy-comerce/backend/pkg/app_error"
	cfg "github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserID     uint     `json:"user_id"`
	RoleTypess []string `json:"role_types"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
	jwt.RegisteredClaims
}

func GenerateNewToken(isFallback ...bool) (string, error) {
	ts := time.Now().UnixNano()
	randBytes := make([]byte, 6)
	if _, err := rand.Read(randBytes); err != nil {
		if len(isFallback) > 0 && isFallback[0] {
			return timeutil.NowUTC().Format("20060102150405"), nil
		}
		return "", e.WrapServerError("failed to generate token", err)
	}
	return fmt.Sprintf("%x%x", ts, randBytes), nil
}

func GenerateNewTokenWithExpiryTime(expiryHours int) (string, time.Time, error) {
	token, err := GenerateNewToken(true)
	if err != nil {
		return "", time.Time{}, e.WrapServerError("failed to generate token", err)
	}

	expiresAt := timeutil.AddHoursUTC(expiryHours)
	return token, expiresAt, nil
}

func VerifyJwtToken(tokenString string, jwtSecret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			l.Error("❌ unexpected signing method: %v", token.Header["alg"], "method", "VerifyJwtToken")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, e.WrapServerError("Invalid/expired jwt token", err)
	}

	claims, ok := token.Claims.(*JwtClaims)
	if ok && token.Valid {
		return claims, nil
	}

	if claims.ExpiresAt.Before(timeutil.NowUTC()) {
		l.Error("❌ Jwt Token expired", "method", "VerifyJwtToken", "error", err)
		return nil, errors.New("jwt token expired")
	}

	return nil, errors.New("jwt token invalid")
}

func GenerateNewJwtToken(userID uint, userEmail string, userName string, roleTypes []string, jwtSecret string) (string, int64, error) {
	now := timeutil.NowUTC()
	expirationTime := timeutil.AddHoursUTC(c.AccessTokenExpiryHours)
	expiresAt := expirationTime.Unix()

	claims := &JwtClaims{
		UserID:     userID,
		Email:      userEmail,
		RoleTypess: roleTypes,
		Username:   userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    c.ProjectName,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, e.WrapServerError("failed to generate token", err)
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
		return nil, e.WrapServerError("Invalid/expired jwt token", err)
	}
	return claims, nil
}

func GenerateNewUUID() (uuid.UUID, error) {
	if strings.EqualFold(cfg.GetUuidType(), c.UuidSequence) {
		return uuid.NewV7()
	}
	return uuid.NewRandom()
}
