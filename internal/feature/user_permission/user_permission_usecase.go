package user_permission

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/internal/feature/shared"
	"github.com/easy-comerce/backend/pkg/logger"
)

type UserPermissionUseCase struct {
	userPermissionRepo *UserPermissionRepository
	userRepo           shared.UserRepositoryInterface
}

func NewUserPermissionUseCase(userRepo shared.UserRepositoryInterface) *UserPermissionUseCase {
	return &UserPermissionUseCase{
		userPermissionRepo: NewUserPermissionRepository(),
		userRepo:           userRepo,
	}
}

func (uc *UserPermissionUseCase) GetUserPermissions(userID uint) ([]string, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	permissions, err := uc.userPermissionRepo.GetUserPermissions(userID)
	if err != nil {
		logger.Logger.Error("Failed to get user permissions", "method", "GetUserPermissions", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	return permissions, nil
}

func (uc *UserPermissionUseCase) GetUserPermissionsByRole(userID uint, roleType string) ([]string, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	if roleType == "" {
		return nil, errors.New("role type cannot be empty")
	}

	permissions, err := uc.userPermissionRepo.GetUserPermissionsByRole(userID, roleType)
	if err != nil {
		logger.Logger.Error("Failed to get user permissions by role", "method", "GetUserPermissionsByRole", "error", err, "userID", userID, "roleType", roleType)
		return nil, fmt.Errorf("failed to get user permissions by role: %w", err)
	}

	return permissions, nil
}

func (uc *UserPermissionUseCase) HasPermission(userID uint, permission string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	if permission == "" {
		return false, errors.New("permission cannot be empty")
	}

	hasPermission, err := uc.userPermissionRepo.HasPermission(userID, permission)
	if err != nil {
		logger.Logger.Error("Failed to check user permission", "method", "HasPermission", "error", err, "userID", userID, "permission", permission)
		return false, fmt.Errorf("failed to check user permission: %w", err)
	}

	return hasPermission, nil
}

func (uc *UserPermissionUseCase) HasAnyPermission(userID uint, permissions []string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	if len(permissions) == 0 {
		return true, nil
	}

	hasPermission, err := uc.userPermissionRepo.HasAnyPermission(userID, permissions)
	if err != nil {
		logger.Logger.Error("Failed to check user permissions", "method", "HasAnyPermission", "error", err, "userID", userID, "permissions", permissions)
		return false, fmt.Errorf("failed to check user permissions: %w", err)
	}

	return hasPermission, nil
}

func (uc *UserPermissionUseCase) SyncUserPermissions(userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	// Refresh the materialized view (PostgreSQL handles the sync automatically)
	err := uc.userPermissionRepo.RefreshMaterializedView()
	if err != nil {
		logger.Logger.Error("Failed to refresh materialized view", "method", "SyncUserPermissions", "error", err, "userID", userID)
		return fmt.Errorf("failed to refresh materialized view: %w", err)
	}

	logger.Logger.Info("User permissions synced successfully via materialized view", "method", "SyncUserPermissions", "userID", userID)
	return nil
}

func (uc *UserPermissionUseCase) SyncAllUsersPermissions() error {
	// Refresh the materialized view (PostgreSQL handles the sync automatically)
	err := uc.userPermissionRepo.RefreshMaterializedView()
	if err != nil {
		logger.Logger.Error("Failed to refresh materialized view", "method", "SyncAllUsersPermissions", "error", err)
		return fmt.Errorf("failed to refresh materialized view: %w", err)
	}

	logger.Logger.Info("All user permissions synced successfully via materialized view", "method", "SyncAllUsersPermissions")
	return nil
}

// GetUserStatusAndPermission gets user status and checks for specific permission in single query
func (uc *UserPermissionUseCase) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, hasPermission bool, err error) {
	if userID == 0 {
		return false, false, false, errors.New("invalid user ID")
	}

	if permission == "" {
		return false, false, false, errors.New("permission cannot be empty")
	}

	return uc.userPermissionRepo.GetUserStatusAndPermission(userID, permission)
}

// GetUserStatusAndAnyPermission gets user status and checks for any of the specified permissions in single query
func (uc *UserPermissionUseCase) GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, hasPermission bool, err error) {
	if userID == 0 {
		return false, false, false, errors.New("invalid user ID")
	}

	return uc.userPermissionRepo.GetUserStatusAndAnyPermission(userID, permissions)
}

func (uc *UserPermissionUseCase) GetUsersWithPermission(permission string) ([]uint, error) {
	if permission == "" {
		return nil, errors.New("permission cannot be empty")
	}

	userIDs, err := uc.userPermissionRepo.GetUsersWithPermission(permission)
	if err != nil {
		logger.Logger.Error("Failed to get users with permission", "method", "GetUsersWithPermission", "error", err, "permission", permission)
		return nil, fmt.Errorf("failed to get users with permission: %w", err)
	}

	return userIDs, nil
}

func (uc *UserPermissionUseCase) GetPermissionStats() (map[string]int, error) {
	stats, err := uc.userPermissionRepo.GetPermissionStats()
	if err != nil {
		logger.Logger.Error("Failed to get permission stats", "method", "GetPermissionStats", "error", err)
		return nil, fmt.Errorf("failed to get permission stats: %w", err)
	}

	return stats, nil
}

func (uc *UserPermissionUseCase) CleanupOrphanedPermissions() error {
	err := uc.userPermissionRepo.CleanupOrphanedPermissions()
	if err != nil {
		logger.Logger.Error("Failed to cleanup orphaned permissions", "method", "CleanupOrphanedPermissions", "error", err)
		return fmt.Errorf("failed to cleanup orphaned permissions: %w", err)
	}

	logger.Logger.Info("Orphaned permissions cleaned up successfully", "method", "CleanupOrphanedPermissions")
	return nil
}
