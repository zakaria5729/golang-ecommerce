package middleware

import (
	"context"
	"net/http"

	p "github.com/easy-comerce/backend/internal/permission"
	"github.com/easy-comerce/backend/internal/user"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/tokenutil"
)

type MiddlewareHandler func(http.Handler) http.Handler

type HandlerFunc func(http.ResponseWriter, *http.Request)

type PermissionMiddleware struct {
	permissionService *p.PermissionService
	userRepo          *user.UserRepository
	jwtSecret         string
}

func NewPermissionMiddleware(
	jwtSecret string,
	userRepo *user.UserRepository,
	permissionRepo *p.PermissionRepository,
) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionService: p.NewPermissionService(permissionRepo),
		userRepo:          userRepo,
		jwtSecret:         jwtSecret,
	}
}

func (pm *PermissionMiddleware) RequireAuthUserStatus() MiddlewareHandler {
	return loadAuthUser(pm, false, false, false)
}

func (pm *PermissionMiddleware) RequireAuthUserWithRolePermission() MiddlewareHandler {
	return loadAuthUser(pm, true, true, true)
}

func (pm *PermissionMiddleware) RequirePermission(permission string) MiddlewareHandler {
	return loadPermissionsStatus(pm, []string{permission}, "RequirePermission")
}

func (pm *PermissionMiddleware) RequireAnyPermission(permissions []string) MiddlewareHandler {
	return loadPermissionsStatus(pm, permissions, "RequireAnyPermission")
}

func loadPermissionsStatus(pm *PermissionMiddleware, permissions []string, methodName string) MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
			if err != nil {
				l.Logger.Error("❌ Token validation failed", "method", methodName, "error", err)
				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			err = nil
			banned := false
			verified := true
			hasPermission := false
			var refreshToken *string

			if len(permissions) == 1 {
				banned, verified, refreshToken, hasPermission, err = pm.permissionService.GetUserStatusAndPermission(claims.UserID, permissions[0])
			} else {
				banned, verified, refreshToken, hasPermission, err = pm.permissionService.GetUserStatusAndAnyPermission(claims.UserID, permissions)
			}

			if err != nil {
				l.Logger.Error("❌ Failed to get user status and permissions", "method", methodName, "error", err, "userID", claims.UserID, "permissions", permissions)
				response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if !isUserActionValid(w, refreshToken, banned, verified) {
				return
			}

			if !hasPermission {
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), c.UserIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func loadAuthUser(pm *PermissionMiddleware, loadFullUser bool, includeRoles bool, includePermissions bool) MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
			if err != nil {
				l.Logger.Error("❌ Invalid/expired jwt token", "method", "RequireAuth", "error", err)
				response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
				return
			}

			banned := false
			verified := false
			var refreshToken *string
			var user *user.UserEntity

			if !loadFullUser {
				banned, verified, refreshToken, err = pm.userRepo.GetAuthUserStatusByID(claims.UserID)
			} else {
				user, err = pm.userRepo.GetAuthUserByID(claims.UserID, includeRoles, includePermissions)
				if user != nil {
					banned = user.Banned
					verified = user.Verified
					refreshToken = user.RefreshToken
				}
			}

			if err != nil {
				l.Logger.Error("❌ User not found", "method", "RequireAuth", "error", err, "userID", claims.UserID)
				response.SendErrorJSON(w, "User not found", http.StatusUnauthorized)
				return
			}

			if !isUserActionValid(w, refreshToken, banned, verified) {
				return
			}

			ctx := context.WithValue(r.Context(), c.UserIDContextKey, claims.UserID)
			if loadFullUser {
				ctx = context.WithValue(ctx, c.UserContextKey, user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isUserActionValid(w http.ResponseWriter, refreshToken *string, banned bool, verified bool) bool {
	if refreshToken == nil {
		response.SendErrorJSON(w, "Invalid access/refresh token", http.StatusUnauthorized)
		return false
	}

	if banned {
		response.SendErrorJSON(w, "Account is banned", http.StatusUnauthorized)
		return false
	}

	if !verified {
		response.SendErrorJSON(w, "Account not verified", http.StatusUnauthorized)
		return false
	}

	return true
}
