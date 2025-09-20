package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/internal/feature/user"
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

func NewPermissionMiddleware(jwtSecret string) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionUseCase: permission.NewPermissionUseCase(),
		authUseCase:       auth.NewAuthUseCase(jwtSecret),
		userUseCase:       user.NewUserUseCase(),
		jwtSecret:         jwtSecret,
	}
}

// **REQUIRED
func (pm *PermissionMiddleware) RequireAuthUserStatus() t.MiddlewareHandler {
	return pm.loadAuthUser(false, false, false)
}

// **REQUIRED
func (pm *PermissionMiddleware) RequireAuthUser() t.MiddlewareHandler {
	return pm.loadAuthUser(true, false, false)
}

// **REQUIRED
func (pm *PermissionMiddleware) RequireAuthWithRolePermission() t.MiddlewareHandler {
	return pm.loadAuthUser(true, true, true)
}

// **REQUIRED
func (pm *PermissionMiddleware) RequirePermission(permission string) t.MiddlewareHandler {
	return pm.loadPermissionsStatus([]string{permission})
}

// **REQUIRED
func (pm *PermissionMiddleware) RequireAnyPermission(permissions []string) t.MiddlewareHandler {
	return pm.loadPermissionsStatus(permissions)
}

func (pm *PermissionMiddleware) GetJWTSecret() string {
	return pm.jwtSecret
}

// **REQUIRED
func (pm *PermissionMiddleware) loadPermissionsStatus(permissions []string) t.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
			if err != nil {
				logger.Logger.Error("Token validation failed", "method", "RequireAnyPermission", "error", err)
				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			err = nil
			banned := false
			verified := true
			hasPermission := false

			if len(permissions) == 1 {
				banned, verified, hasPermission, err = pm.permissionUseCase.GetUserStatusAndPermission(claims.UserID, permissions[0])
			} else {
				banned, verified, hasPermission, err = pm.permissionUseCase.GetUserStatusAndAnyPermission(claims.UserID, permissions)
			}

			if err != nil {
				logger.Logger.Error("Failed to get user status and permissions", "method", "RequireAnyPermission", "error", err, "userID", claims.UserID, "permissions", permissions)
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if banned {
				logger.Logger.Warn("Banned user attempted to access protected resource", "method", "RequireAnyPermission", "userID", claims.UserID)
				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
				return
			}

			if !verified {
				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", "RequireAnyPermission", "userID", claims.UserID)
				response.SendErrorJSON(w, "Account not verified", http.StatusForbidden)
				return
			}

			if !hasPermission {
				logger.Logger.Warn("User lacks any of the required permissions", "method", "RequireAnyPermission", "userID", claims.UserID, "permissions", permissions)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), constants.UserIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// **REQUIRED
func (pm *PermissionMiddleware) loadAuthUser(loadFullUser bool, includeRoles bool, includePermissions bool) t.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
			if err != nil {
				logger.Logger.Error("Invalid/expired token", "method", "RequireAuth", "error", err)
				response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
				return
			}

			banned := false
			verified := false
			user := &user.User{}

			if !loadFullUser {
				banned, verified, err = pm.userUseCase.GetAuthUserStatusByID(claims.UserID)
			} else {
				user, err = pm.userUseCase.GetAuthUserByID(claims.UserID, includeRoles, includePermissions)
			}

			if err != nil {
				logger.Logger.Error("User not found", "method", "RequireAuth", "error", err, "userID", claims.UserID)
				response.SendErrorJSON(w, "User not found", http.StatusUnauthorized)
				return
			}

			if !verified {
				logger.Logger.Error("User is not verified", "method", "RequireAuth", "userID", claims.UserID)
				response.SendErrorJSON(w, "Account not verified yet", http.StatusForbidden)
				return
			}

			if banned {
				logger.Logger.Error("User is banned", "method", "RequireAuth", "userID", claims.UserID)
				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
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

// func (pm *PermissionMiddleware) RequireRole(roleType string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Validate JWT token and get user ID (no database calls)
// 			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
// 			if err != nil {
// 				logger.Logger.Error("Token validation failed", "method", "RequireRole", "error", err)
// 				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			// Check user status from users table (banned, verified)
// 			banned, verified, err := pm.userUseCase.GetUserStatus(claims.UserID)
// 			if err != nil {
// 				logger.Logger.Error("Failed to get user status", "method", "RequireRole", "error", err, "userID", claims.UserID)
// 				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
// 				return
// 			}

// 			if banned {
// 				logger.Logger.Warn("Banned user attempted to access protected resource", "method", "RequireRole", "userID", claims.UserID)
// 				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
// 				return
// 			}

// 			if !verified {
// 				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", "RequireRole", "userID", claims.UserID)
// 				response.SendErrorJSON(w, "Account not verified", http.StatusForbidden)
// 				return
// 			}

// 			// Get user permissions for the specific role using 5-table joins
// 			permissions, err := pm.permissionUseCase.GetUserPermissionsByRole(claims.UserID, roleType)
// 			if err != nil {
// 				logger.Logger.Error("Failed to get user permissions by role", "method", "RequireRole", "error", err, "userID", claims.UserID, "roleType", roleType)
// 				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
// 				return
// 			}

// 			if len(permissions) == 0 {
// 				logger.Logger.Warn("User lacks required role", "method", "RequireRole", "userID", claims.UserID, "roleType", roleType)
// 				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
// 				return
// 			}

// 			// Add user ID and role to context for downstream handlers
// 			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
// 			ctx = context.WithValue(ctx, "role_type", roleType)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// // RequireSuperAdmin creates middleware that requires super admin role
// func (pm *PermissionMiddleware) RequireSuperAdmin() func(http.Handler) http.Handler {
// 	return pm.RequireRole("SUPER_ADMIN")
// }

// // RequireAdmin creates middleware that requires admin role or higher
// func (pm *PermissionMiddleware) RequireAdmin() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Validate JWT token and get user ID (no database calls)
// 			claims, err := tokenutil.ValidateTokenAndGetJwtClaims(r, pm.jwtSecret)
// 			if err != nil {
// 				logger.Logger.Error("Token validation failed", "method", "RequireAdmin", "error", err)
// 				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			// Check user status and admin permissions in single query
// 			adminPermissions := []string{
// 				"system.admin",
// 				"system.super_admin",
// 			}

// 			banned, verified, hasPermission, err := pm.permissionUseCase.GetUserStatusAndAnyPermission(claims.UserID, adminPermissions)
// 			if err != nil {
// 				logger.Logger.Error("Failed to get user status and admin permissions", "method", "RequireAdmin", "error", err, "userID", claims.UserID)
// 				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
// 				return
// 			}

// 			if banned {
// 				logger.Logger.Warn("Banned user attempted to access protected resource", "method", "RequireAdmin", "userID", claims.UserID)
// 				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
// 				return
// 			}

// 			if !verified {
// 				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", "RequireAdmin", "userID", claims.UserID)
// 				response.SendErrorJSON(w, "Account not verified", http.StatusForbidden)
// 				return
// 			}

// 			if !hasPermission {
// 				logger.Logger.Warn("User lacks admin permissions", "method", "RequireAdmin", "userID", claims.UserID)
// 				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
// 				return
// 			}

// 			// Add user ID to context for downstream handlers
// 			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// extractToken extracts JWT token from Authorization header
// func (pm *AuthPermissionMiddleware) extractJwtToken(r *http.Request) string {
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		return ""
// 	}

// 	parts := strings.Split(authHeader, " ")
// 	if len(parts) != 2 || parts[0] != "Bearer" {
// 		return ""
// 	}

// 	return parts[1]
// }

// Helper function to extract user ID from JWT token
// func ExtractUserIDFromToken(r *http.Request) (uint, error) {
// 	// Get user from context (set by auth middleware)
// 	user, ok := r.Context().Value(constants.UserContextKey).(*user.User)
// 	if !ok || user == nil {
// 		return 0, errors.New("user not found in context")
// 	}
// 	return user.ID, nil
// }

// // Helper function to check permission in handlers
// func CheckPermission(userID uint, permission string) (bool, error) {
// 	pm := NewPermissionMiddleware("")
// 	return pm.permissionUseCase.HasPermission(userID, permission)
// }

// // Helper function to check multiple permissions in handlers
// func CheckAnyPermission(userID uint, permissions []string) (bool, error) {
// 	pm := NewPermissionMiddleware("")
// 	return pm.permissionUseCase.HasAnyPermission(userID, permissions)
// }

// validateTokenAndGetUserID validates JWT token and returns user ID without any database calls
// func (pm *AuthPermissionMiddleware) validateTokenAndGetUserID(r *http.Request) (uint, error) {
// 	token := pm.extractJwtToken(r)
// 	if token == "" {
// 		return 0, errors.New("no token provided")
// 	}

// 	claims, err := pm.authUseCase.VerifyToken(token)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Return user ID directly - no database calls needed
// 	// Permission checking will handle user validation through 5-table joins
// 	return claims.UserID, nil
// }
