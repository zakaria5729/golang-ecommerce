package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/db"
	p "github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	t "github.com/easy-comerce/backend/pkg/types"
)

type PermissionMiddleware struct {
	permissionService *p.PermissionService
	userService       *user.UserService
	jwtSecret         string
}

func NewPermissionMiddleware() *PermissionMiddleware {
	db := db.GetDB()
	cfg := config.GetConfig()
	roleRepo := role.NewRoleRepository(db)
	userRepo := user.NewUserRepository(db)
	permissionRepo := p.NewPermissionRepository(db)

	return &PermissionMiddleware{
		permissionService: p.NewPermissionService(permissionRepo),
		userService:       user.NewUserService(userRepo, roleRepo),
		jwtSecret:         cfg.JWTSecret,
	}
}

func (pm *PermissionMiddleware) RequireAuthUserStatus() t.MiddlewareHandler {
	return loadAuthUser(pm, false, false, false)
}

func (pm *PermissionMiddleware) RequireAuthUser() t.MiddlewareHandler {
	return loadAuthUser(pm, true, false, false)
}

func (pm *PermissionMiddleware) RequireAuthWithRolePermission() t.MiddlewareHandler {
	return loadAuthUser(pm, true, true, true)
}

func (pm *PermissionMiddleware) RequirePermission(permission string) t.MiddlewareHandler {
	return loadPermissionsStatus(pm, []string{permission}, "RequirePermission")
}

func (pm *PermissionMiddleware) RequireAnyPermission(permissions []string) t.MiddlewareHandler {
	return loadPermissionsStatus(pm, permissions, "RequireAnyPermission")
}

func (pm *PermissionMiddleware) GetJWTSecret() string {
	return pm.jwtSecret
}

func loadPermissionsStatus(pm *PermissionMiddleware, permissions []string, methodName string) t.MiddlewareHandler {
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
				banned, verified, refreshToken, hasPermission, err = pm.permissionService.GetUserStatusAndPermission(claims.UserID, permissions[0])
			} else {
				banned, verified, refreshToken, hasPermission, err = pm.permissionService.GetUserStatusAndAnyPermission(claims.UserID, permissions)
			}

			if err != nil {
				logger.Logger.Error("Failed to get user status and permissions", "method", methodName, "error", err, "userID", claims.UserID, "permissions", permissions)
				response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if refreshToken == nil {
				response.SendErrorJSON(w, "Invalid access/refresh token", http.StatusUnauthorized)
				return
			}

			if banned {
				response.SendErrorJSON(w, "Account is banned", http.StatusUnauthorized)
				return
			}

			if !verified {
				response.SendErrorJSON(w, "Account not verified", http.StatusUnauthorized)
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

func loadAuthUser(pm *PermissionMiddleware, loadFullUser bool, includeRoles bool, includePermissions bool) t.MiddlewareHandler {
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
				banned, verified, refreshToken, err = pm.userService.GetAuthUserStatusByID(claims.UserID)
			} else {
				user, err = pm.userService.GetAuthUserByID(claims.UserID, includeRoles, includePermissions)
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
				response.SendErrorJSON(w, "Invalid access/refresh token", http.StatusUnauthorized)
				return
			}

			if !verified {
				response.SendErrorJSON(w, "Account not verified yet", http.StatusUnauthorized)
				return
			}

			if banned {
				response.SendErrorJSON(w, "Account is banned", http.StatusUnauthorized)
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

func GetUserFromContext(ctx context.Context) (*user.User, error) {
	user, ok := ctx.Value(c.UserContextKey).(*user.User)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

func GetUserIDFromContext(ctx context.Context) (*uint, error) {
	userID, ok := ctx.Value(c.UserIDContextKey).(uint)
	if !ok || userID == 0 {
		return nil, errors.New("user ID not found or invalid type in context")
	}
	return &userID, nil
}

func GetUserIdOnlyFromContext(ctx context.Context) *uint {
	userID, ok := ctx.Value(c.UserIDContextKey).(uint)
	if !ok || userID == 0 {
		return nil
	}
	return &userID
}
