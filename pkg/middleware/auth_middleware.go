package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type AuthMiddleware struct {
	authUseCase *auth.AuthUseCase
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: auth.NewAuthUseCase(jwtSecret),
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			logger.Logger.Error("No token provided", "method", "RequireAuth")
			response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		claims, err := m.authUseCase.VerifyToken(token)
		if err != nil {
			logger.Logger.Error("Invalid token", "method", "RequireAuth", "error", err)
			response.SendErrorJSON(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, err := m.authUseCase.GetUserProfile(claims.UserID, "roles,permissions")
		if err != nil {
			logger.Logger.Error("User not found", "method", "RequireAuth", "error", err, "userID", claims.UserID)
			response.SendErrorJSON(w, "User not found", http.StatusUnauthorized)
			return
		}

		if user.Banned {
			logger.Logger.Error("User is banned", "method", "RequireAuth", "userID", user.ID)
			response.SendErrorJSON(w, "Account is banned", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) RequireRole(roleType string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := m.getUserFromContext(r)
			if user == nil {
				response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasRole, err := m.authUseCase.HasRole(user.ID, roleType)
			if err != nil {
				logger.Logger.Error("Failed to check user role", "method", "RequireRole", "error", err, "userID", user.ID, "roleType", roleType)
				response.SendErrorJSON(w, "Failed to verify role", http.StatusInternalServerError)
				return
			}

			if !hasRole {
				logger.Logger.Error("Insufficient role", "method", "RequireRole", "userID", user.ID, "roleType", roleType)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func (m *AuthMiddleware) RequirePermission(permission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := m.getUserFromContext(r)
			if user == nil {
				response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasPermission, err := m.authUseCase.HasPermission(user.ID, permission)
			if err != nil {
				logger.Logger.Error("Failed to check user permission", "method", "RequirePermission", "error", err, "userID", user.ID, "permission", permission)
				response.SendErrorJSON(w, "Failed to verify permission", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				logger.Logger.Error("Insufficient permission", "method", "RequirePermission", "userID", user.ID, "permission", permission)
				response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func (m *AuthMiddleware) RequireAnyRole(roleTypes ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := m.getUserFromContext(r)
			if user == nil {
				response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			for _, roleType := range roleTypes {
				hasRole, err := m.authUseCase.HasRole(user.ID, roleType)
				if err != nil {
					logger.Logger.Error("Failed to check user role", "method", "RequireAnyRole", "error", err, "userID", user.ID, "roleType", roleType)
					response.SendErrorJSON(w, "Failed to verify role", http.StatusInternalServerError)
					return
				}
				if hasRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			logger.Logger.Error("Insufficient role", "method", "RequireAnyRole", "userID", user.ID, "roleTypes", roleTypes)
			response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
		}
	}
}

func (m *AuthMiddleware) RequireAnyPermission(permissions ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := m.getUserFromContext(r)
			if user == nil {
				response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			for _, permission := range permissions {
				hasPermission, err := m.authUseCase.HasPermission(user.ID, permission)
				if err != nil {
					logger.Logger.Error("Failed to check user permission", "method", "RequireAnyPermission", "error", err, "userID", user.ID, "permission", permission)
					response.SendErrorJSON(w, "Failed to verify permission", http.StatusInternalServerError)
					return
				}
				if hasPermission {
					next.ServeHTTP(w, r)
					return
				}
			}

			logger.Logger.Error("Insufficient permission", "method", "RequireAnyPermission", "userID", user.ID, "permissions", permissions)
			response.SendErrorJSON(w, "Insufficient permissions", http.StatusForbidden)
		}
	}
}

func (m *AuthMiddleware) extractToken(r *http.Request) string {
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

func (m *AuthMiddleware) getUserFromContext(r *http.Request) *auth.User {
	user, ok := r.Context().Value(UserContextKey).(*auth.User)
	if !ok {
		return nil
	}
	return user
}

func GetUserFromContext(r *http.Request) *auth.User {
	user, ok := r.Context().Value(UserContextKey).(*auth.User)
	if !ok {
		return nil
	}
	return user
}
