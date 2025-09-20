package permission

import (
	"errors"

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

func (uc *PermissionUseCase) GetAllPermissions(sortBy string, sortOrder string) ([]Permission, error) {
	permissions, err := uc.permissionRepo.GetAllPermissions(sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to get all permissions", "method", "GetAllPermissions", "error", err)
		return nil, err
	}
	return permissions, nil
}

func (uc *PermissionUseCase) GetPermissionByID(id uint) (*Permission, error) {
	permission, err := uc.permissionRepo.GetPermissionByID(id)
	if err != nil {
		logger.Logger.Error("Failed to get permission by ID", "method", "GetPermissionByID", "error", err, "id", id)
		return nil, err
	}
	return permission, nil
}

func (uc *PermissionUseCase) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, hasPermission bool, err error) {
	if userID <= 0 {
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

	if len(permissions) == 0 {
		return false, false, false, errors.New("permissions cannot be empty")
	}

	return uc.permissionRepo.GetUserStatusAndAnyPermission(userID, permissions)
}
