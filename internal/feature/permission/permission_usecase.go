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
