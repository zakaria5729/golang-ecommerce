package permission

import (
	"errors"
	"strings"

	"github.com/easy-comerce/backend/pkg/logger"
)

type PermissionUseCase struct {
	permissionRepo *PermissionRepository
}

func NewPermissionUseCase() *PermissionUseCase {
	return &PermissionUseCase{
		permissionRepo: NewPermissionRepository(),
	}
}

func (uc *PermissionUseCase) GetAllPermissions(include, sortBy, sortOrder string) ([]Permission, error) {
	var includeList []string
	if include != "" {
		includeList = strings.Split(include, ",")
	}

	permissions, err := uc.permissionRepo.GetAllPermissions(includeList, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to get all permissions", "method", "GetAllPermissions", "error", err)
		return nil, err
	}
	return permissions, nil
}

func (uc *PermissionUseCase) GetPermissionByID(id uint, include string) (*Permission, error) {
	var includeList []string
	if include != "" {
		includeList = strings.Split(include, ",")
	}

	permission, err := uc.permissionRepo.GetPermissionByID(id, includeList)
	if err != nil {
		logger.Logger.Error("Failed to get permission by ID", "method", "GetPermissionByID", "error", err, "id", id)
		return nil, err
	}
	return permission, nil
}

func (uc *PermissionUseCase) CreatePermission(permission *Permission) (*Permission, error) {
	// Sanitize input
	permission.Sanitize()

	// Validate permission
	if err := permission.IsValid(); err != nil {
		logger.Logger.Error("Invalid permission data", "method", "CreatePermission", "error", err)
		return nil, err
	}

	// Check if permission already exists
	exists, err := uc.permissionRepo.PermissionExistsByName(permission.Name, nil)
	if err != nil {
		logger.Logger.Error("Failed to check if permission exists", "method", "CreatePermission", "error", err)
		return nil, err
	}
	if exists {
		logger.Logger.Error("Permission already exists", "method", "CreatePermission", "name", permission.Name)
		return nil, errors.New("permission already exists")
	}

	// Create permission
	err = uc.permissionRepo.CreatePermission(permission)
	if err != nil {
		logger.Logger.Error("Failed to create permission", "method", "CreatePermission", "error", err)
		return nil, err
	}

	logger.Logger.Info("Permission created successfully", "method", "CreatePermission", "id", permission.ID, "name", permission.Name)
	return permission, nil
}

func (uc *PermissionUseCase) UpdatePermission(id uint, permission *Permission) (*Permission, error) {
	// Sanitize input
	permission.Sanitize()

	// Validate permission
	if err := permission.IsValid(); err != nil {
		logger.Logger.Error("Invalid permission data", "method", "UpdatePermission", "error", err)
		return nil, err
	}

	// Check if permission exists
	existingPermission, err := uc.permissionRepo.GetPermissionByID(id, []string{})
	if err != nil {
		logger.Logger.Error("Failed to get existing permission", "method", "UpdatePermission", "error", err, "id", id)
		return nil, err
	}

	// Check if new name conflicts with existing permission (excluding current one)
	if permission.Name != existingPermission.Name {
		exists, err := uc.permissionRepo.PermissionExistsByName(permission.Name, &id)
		if err != nil {
			logger.Logger.Error("Failed to check if permission exists", "method", "UpdatePermission", "error", err)
			return nil, err
		}
		if exists {
			logger.Logger.Error("Permission name already exists", "method", "UpdatePermission", "name", permission.Name)
			return nil, errors.New("permission name already exists")
		}
	}

	// Set the ID for update
	permission.ID = id

	// Update permission
	err = uc.permissionRepo.UpdatePermission(permission)
	if err != nil {
		logger.Logger.Error("Failed to update permission", "method", "UpdatePermission", "error", err, "id", id)
		return nil, err
	}

	logger.Logger.Info("Permission updated successfully", "method", "UpdatePermission", "id", id, "name", permission.Name)
	return permission, nil
}

func (uc *PermissionUseCase) DeletePermission(id uint) error {
	// Check if permission exists
	_, err := uc.permissionRepo.GetPermissionByID(id, []string{})
	if err != nil {
		logger.Logger.Error("Failed to get permission for deletion", "method", "DeletePermission", "error", err, "id", id)
		return err
	}

	// Delete permission
	if err := uc.permissionRepo.DeletePermission(id); err != nil {
		logger.Logger.Error("Failed to delete permission", "method", "DeletePermission", "error", err, "id", id)
		return err
	}

	logger.Logger.Info("Permission deleted successfully", "method", "DeletePermission", "id", id)
	return nil
}

func (uc *PermissionUseCase) GetUserPermissions(userID uint) ([]string, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	permissions, err := uc.permissionRepo.GetUserPermissions(userID)
	if err != nil {
		logger.Logger.Error("Failed to get user permissions", "method", "GetUserPermissions", "error", err, "userID", userID)
		return nil, err
	}

	return permissions, nil
}

func (uc *PermissionUseCase) GetUserPermissionsByRole(userID uint, roleType string) ([]string, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	if roleType == "" {
		return nil, errors.New("role type cannot be empty")
	}

	permissions, err := uc.permissionRepo.GetUserPermissionsByRole(userID, roleType)
	if err != nil {
		logger.Logger.Error("Failed to get user permissions by role", "method", "GetUserPermissionsByRole", "error", err, "userID", userID, "roleType", roleType)
		return nil, err
	}

	return permissions, nil
}

func (uc *PermissionUseCase) HasPermission(userID uint, permission string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	if permission == "" {
		return false, errors.New("permission cannot be empty")
	}

	hasPermission, err := uc.permissionRepo.HasPermission(userID, permission)
	if err != nil {
		logger.Logger.Error("Failed to check user permission", "method", "HasPermission", "error", err, "userID", userID, "permission", permission)
		return false, err
	}

	return hasPermission, nil
}

func (uc *PermissionUseCase) HasAnyPermission(userID uint, permissions []string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	if len(permissions) == 0 {
		return true, nil
	}

	hasPermission, err := uc.permissionRepo.HasAnyPermission(userID, permissions)
	if err != nil {
		logger.Logger.Error("Failed to check user permissions", "method", "HasAnyPermission", "error", err, "userID", userID, "permissions", permissions)
		return false, err
	}

	return hasPermission, nil
}

func (uc *PermissionUseCase) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, hasPermission bool, err error) {
	if userID == 0 {
		return false, false, false, errors.New("invalid user ID")
	}

	if permission == "" {
		return false, false, false, errors.New("permission cannot be empty")
	}

	return uc.permissionRepo.GetUserStatusAndPermission(userID, permission)
}

func (uc *PermissionUseCase) GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, hasPermission bool, err error) {
	if userID == 0 {
		return false, false, false, errors.New("invalid user ID")
	}

	return uc.permissionRepo.GetUserStatusAndAnyPermission(userID, permissions)
}

func (uc *PermissionUseCase) GetUsersWithPermission(permission string) ([]uint, error) {
	if permission == "" {
		return nil, errors.New("permission cannot be empty")
	}

	userIDs, err := uc.permissionRepo.GetUsersWithPermission(permission)
	if err != nil {
		logger.Logger.Error("Failed to get users with permission", "method", "GetUsersWithPermission", "error", err, "permission", permission)
		return nil, err
	}

	return userIDs, nil
}

func (uc *PermissionUseCase) GetPermissionStats() (map[string]int, error) {
	stats, err := uc.permissionRepo.GetPermissionStats()
	if err != nil {
		logger.Logger.Error("Failed to get permission stats", "method", "GetPermissionStats", "error", err)
		return nil, err
	}

	return stats, nil
}
