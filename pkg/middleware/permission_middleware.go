package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	t "github.com/easy-comerce/backend/pkg/types"
)

type PermissionMiddleware struct {
	permissionUseCase *permission.PermissionUseCase
	authUseCase       *auth.AuthUseCase
	userUseCase       *user.UserUseCase
	jwtSecret         string
}

func NewPermissionMiddleware() *PermissionMiddleware {
	cfg := config.GetConfig()
	return &PermissionMiddleware{
		permissionUseCase: permission.NewPermissionUseCase(),
		authUseCase:       auth.NewAuthUseCase(cfg.JWTSecret),
		userUseCase:       user.NewUserUseCase(),
		jwtSecret:         cfg.JWTSecret,
	}
}

func (pm *PermissionMiddleware) RequireAuthUserStatus() t.MiddlewareHandler {
	return pm.loadAuthUser(false, false, false)
}

func (pm *PermissionMiddleware) RequireAuthUser() t.MiddlewareHandler {
	return pm.loadAuthUser(true, false, false)
}

func (pm *PermissionMiddleware) RequireAuthWithRolePermission() t.MiddlewareHandler {
	return pm.loadAuthUser(true, true, true)
}

func (pm *PermissionMiddleware) RequirePermission(permission string) t.MiddlewareHandler {
	return pm.loadPermissionsStatus([]string{permission}, "RequirePermission")
}

func (pm *PermissionMiddleware) RequireAnyPermission(permissions []string) t.MiddlewareHandler {
	return pm.loadPermissionsStatus(permissions, "RequireAnyPermission")
}

func (pm *PermissionMiddleware) GetJWTSecret() string {
	return pm.jwtSecret
}

func (pm *PermissionMiddleware) loadPermissionsStatus(permissions []string, methodName string) t.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
			if err != nil {
				logger.Logger.Error("Token validation failed", "method", methodName, "error", err)
				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			err = nil
			banned := false
			verified := true
			hasPermission := false
			var refreshToken *string

			if len(permissions) == 1 {
				banned, verified, refreshToken, hasPermission, err = pm.permissionUseCase.GetUserStatusAndPermission(claims.UserID, permissions[0])
			} else {
				banned, verified, refreshToken, hasPermission, err = pm.permissionUseCase.GetUserStatusAndAnyPermission(claims.UserID, permissions)
			}

			if err != nil {
				logger.Logger.Error("Failed to get user status and permissions", "method", methodName, "error", err, "userID", claims.UserID, "permissions", permissions)
				response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if refreshToken == nil {
				logger.Logger.Error("Invalid refresh token", "method", methodName, "userID", claims.UserID)
				response.SendErrorJSON(w, "Invalid access/refresh token", http.StatusUnauthorized)
				return
			}

			if banned {
				logger.Logger.Warn("Banned user attempted to access protected resource", "method", methodName, "userID", claims.UserID)
				response.SendErrorJSON(w, "Account is banned", http.StatusUnauthorized)
				return
			}

			if !verified {
				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", methodName, "userID", claims.UserID)
				response.SendErrorJSON(w, "Account not verified", http.StatusUnauthorized)
				return
			}

			if !hasPermission {
				logger.Logger.Warn("User lacks any of the required permissions", "method", methodName, "userID", claims.UserID, "permissions", permissions)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), constants.UserIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (pm *PermissionMiddleware) loadAuthUser(loadFullUser bool, includeRoles bool, includePermissions bool) t.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
			if err != nil {
				logger.Logger.Error("Invalid/expired jwt token", "method", "RequireAuth", "error", err)
				response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
				return
			}

			banned := false
			verified := false
			var refreshToken *string
			var user *user.User

			if !loadFullUser {
				banned, verified, refreshToken, err = pm.userUseCase.GetAuthUserStatusByID(claims.UserID)
			} else {
				user, err = pm.userUseCase.GetAuthUserByID(claims.UserID, includeRoles, includePermissions)
				if user != nil {
					banned = user.Banned
					verified = user.Verified
					refreshToken = user.RefreshToken
				}
			}

			if err != nil {
				logger.Logger.Error("User not found", "method", "RequireAuth", "error", err, "userID", claims.UserID)
				response.SendErrorJSON(w, "User not found", http.StatusUnauthorized)
				return
			}

			if refreshToken == nil {
				logger.Logger.Error("Invalid refresh token", "method", "RequireAuth", "userID", claims.UserID)
				response.SendErrorJSON(w, "Invalid access/refresh token", http.StatusUnauthorized)
				return
			}

			if !verified {
				logger.Logger.Error("Account is not verified", "method", "RequireAuth", "userID", claims.UserID)
				response.SendErrorJSON(w, "Account not verified yet", http.StatusUnauthorized)
				return
			}

			if banned {
				logger.Logger.Error("User is banned", "method", "RequireAuth", "userID", claims.UserID)
				response.SendErrorJSON(w, "Account is banned", http.StatusUnauthorized)
				return
			}

			if !loadFullUser {
				ctx := context.WithValue(r.Context(), constants.UserIDContextKey, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				ctx := context.WithValue(r.Context(), constants.UserContextKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func GetUserFromContext(r *http.Request) (*user.User, error) {
	user, ok := r.Context().Value(constants.UserContextKey).(*user.User)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

func GetUserIDFromContext(r *http.Request) (*uint, error) {
	userID, ok := r.Context().Value(constants.UserIDContextKey).(uint)
	if !ok {
		return nil, errors.New("user ID not found or invalid type in context")
	}
	return &userID, nil
}
