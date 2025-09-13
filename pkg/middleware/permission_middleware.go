package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
)

type AuthPermissionMiddleware struct {
	permissionUseCase *permission.PermissionUseCase
	authUseCase       *auth.AuthUseCase
	userUseCase       *user.UserUseCase
	jwtSecret         string
}

func NewAuthPermissionMiddleware(jwtSecret string) *AuthPermissionMiddleware {
	return &AuthPermissionMiddleware{
		permissionUseCase: permission.NewPermissionUseCase(),
		authUseCase:       auth.NewAuthUseCase(jwtSecret),
		userUseCase:       user.NewUserUseCase(),
		jwtSecret:         jwtSecret,
	}
}

// **REQUIRED
func (pm *AuthPermissionMiddleware) RequireAuth(includeRoles bool, includePermissions bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := pm.extractJwtToken(r)

			if token == "" {
				logger.Logger.Error("No token provided", "method", "RequireAuth")
				response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			claims, err := pm.authUseCase.VerifyToken(token)
			if err != nil {
				logger.Logger.Error("Invalid/expired token", "method", "RequireAuth", "error", err)
				response.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
				return
			}

			user, err := pm.userUseCase.GetAuthUserByID(claims.UserID, includeRoles, includePermissions)
			if err != nil {
				logger.Logger.Error("User not found", "method", "RequireAuth", "error", err, "userID", claims.UserID)
				response.SendErrorJSON(w, "User not found", http.StatusUnauthorized)
				return
			}

			if user.Verified {
				logger.Logger.Error("User is not verified", "method", "RequireAuth", "userID", user.ID)
				response.SendErrorJSON(w, "Account not verified yet", http.StatusForbidden)
				return
			}

			if user.Banned {
				logger.Logger.Error("User is banned", "method", "RequireAuth", "userID", user.ID)
				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), constants.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePermission creates middleware that requires authentication and a specific permission
func (pm *AuthPermissionMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate JWT token and get user ID (no database calls)
			userID, err := pm.validateTokenAndGetUserID(r)
			if err != nil {
				logger.Logger.Error("Token validation failed", "method", "RequirePermission", "error", err)
				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check user status and permission in single query
			banned, verified, hasPermission, err := pm.permissionUseCase.GetUserStatusAndPermission(userID, permission)
			if err != nil {
				logger.Logger.Error("Failed to get user status and permission", "method", "RequirePermission", "error", err, "userID", userID, "permission", permission)
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if banned {
				logger.Logger.Warn("Banned user attempted to access protected resource", "method", "RequirePermission", "userID", userID)
				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
				return
			}

			if !verified {
				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", "RequirePermission", "userID", userID)
				response.SendErrorJSON(w, "Account not verified", http.StatusForbidden)
				return
			}

			if !hasPermission {
				logger.Logger.Warn("User lacks required permission", "method", "RequirePermission", "userID", userID, "permission", permission)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user ID and permission to context for downstream handlers
			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "permission", permission)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAnyPermission creates middleware that requires any of the specified permissions
func (pm *AuthPermissionMiddleware) RequireAnyPermission(permissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate JWT token and get user ID (no database calls)
			userID, err := pm.validateTokenAndGetUserID(r)
			if err != nil {
				logger.Logger.Error("Token validation failed", "method", "RequireAnyPermission", "error", err)
				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check user status and permissions in single query
			banned, verified, hasPermission, err := pm.permissionUseCase.GetUserStatusAndAnyPermission(userID, permissions)
			if err != nil {
				logger.Logger.Error("Failed to get user status and permissions", "method", "RequireAnyPermission", "error", err, "userID", userID, "permissions", permissions)
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if banned {
				logger.Logger.Warn("Banned user attempted to access protected resource", "method", "RequireAnyPermission", "userID", userID)
				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
				return
			}

			if !verified {
				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", "RequireAnyPermission", "userID", userID)
				response.SendErrorJSON(w, "Account not verified", http.StatusForbidden)
				return
			}

			if !hasPermission {
				logger.Logger.Warn("User lacks any of the required permissions", "method", "RequireAnyPermission", "userID", userID, "permissions", permissions)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user ID and permissions to context for downstream handlers
			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "permissions", permissions)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole creates middleware that requires a specific role type
func (pm *AuthPermissionMiddleware) RequireRole(roleType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate JWT token and get user ID (no database calls)
			userID, err := pm.validateTokenAndGetUserID(r)
			if err != nil {
				logger.Logger.Error("Token validation failed", "method", "RequireRole", "error", err)
				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check user status from users table (banned, verified)
			banned, verified, err := pm.userUseCase.GetUserStatus(userID)
			if err != nil {
				logger.Logger.Error("Failed to get user status", "method", "RequireRole", "error", err, "userID", userID)
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if banned {
				logger.Logger.Warn("Banned user attempted to access protected resource", "method", "RequireRole", "userID", userID)
				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
				return
			}

			if !verified {
				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", "RequireRole", "userID", userID)
				response.SendErrorJSON(w, "Account not verified", http.StatusForbidden)
				return
			}

			// Get user permissions for the specific role using 5-table joins
			permissions, err := pm.permissionUseCase.GetUserPermissionsByRole(userID, roleType)
			if err != nil {
				logger.Logger.Error("Failed to get user permissions by role", "method", "RequireRole", "error", err, "userID", userID, "roleType", roleType)
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if len(permissions) == 0 {
				logger.Logger.Warn("User lacks required role", "method", "RequireRole", "userID", userID, "roleType", roleType)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user ID and role to context for downstream handlers
			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "role_type", roleType)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireSuperAdmin creates middleware that requires super admin role
func (pm *AuthPermissionMiddleware) RequireSuperAdmin() func(http.Handler) http.Handler {
	return pm.RequireRole("SUPER_ADMIN")
}

// RequireAdmin creates middleware that requires admin role or higher
func (pm *AuthPermissionMiddleware) RequireAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate JWT token and get user ID (no database calls)
			userID, err := pm.validateTokenAndGetUserID(r)
			if err != nil {
				logger.Logger.Error("Token validation failed", "method", "RequireAdmin", "error", err)
				response.SendErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check user status and admin permissions in single query
			adminPermissions := []string{
				"system.admin",
				"system.super_admin",
			}

			banned, verified, hasPermission, err := pm.permissionUseCase.GetUserStatusAndAnyPermission(userID, adminPermissions)
			if err != nil {
				logger.Logger.Error("Failed to get user status and admin permissions", "method", "RequireAdmin", "error", err, "userID", userID)
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if banned {
				logger.Logger.Warn("Banned user attempted to access protected resource", "method", "RequireAdmin", "userID", userID)
				response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
				return
			}

			if !verified {
				logger.Logger.Warn("Unverified user attempted to access protected resource", "method", "RequireAdmin", "userID", userID)
				response.SendErrorJSON(w, "Account not verified", http.StatusForbidden)
				return
			}

			if !hasPermission {
				logger.Logger.Warn("User lacks admin permissions", "method", "RequireAdmin", "userID", userID)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user ID to context for downstream handlers
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken extracts JWT token from Authorization header
func (pm *AuthPermissionMiddleware) extractJwtToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// Helper function to extract user ID from JWT token
func ExtractUserIDFromToken(r *http.Request) (uint, error) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value(constants.UserContextKey).(*user.User)
	if !ok || user == nil {
		return 0, errors.New("user not found in context")
	}
	return user.ID, nil
}

// Helper function to check permission in handlers
func CheckPermission(userID uint, permission string) (bool, error) {
	pm := NewAuthPermissionMiddleware("")
	return pm.permissionUseCase.HasPermission(userID, permission)
}

// Helper function to check multiple permissions in handlers
func CheckAnyPermission(userID uint, permissions []string) (bool, error) {
	pm := NewAuthPermissionMiddleware("")
	return pm.permissionUseCase.HasAnyPermission(userID, permissions)
}

// GetJWTSecret returns the JWT secret used by this middleware
func (pm *AuthPermissionMiddleware) GetJWTSecret() string {
	return pm.jwtSecret
}

// validateTokenAndGetUserID validates JWT token and returns user ID without any database calls
func (pm *AuthPermissionMiddleware) validateTokenAndGetUserID(r *http.Request) (uint, error) {
	token := pm.extractJwtToken(r)
	if token == "" {
		return 0, errors.New("no token provided")
	}

	claims, err := pm.authUseCase.VerifyToken(token)
	if err != nil {
		return 0, err
	}

	// Return user ID directly - no database calls needed
	// Permission checking will handle user validation through 5-table joins
	return claims.UserID, nil
}
