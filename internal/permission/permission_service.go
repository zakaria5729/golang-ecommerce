package permission

import (
	"errors"

	"github.com/easy-comerce/backend/internal/permission/model"
)

type PermissionService interface {
	GetAllPermissions(sortBy string, sortOrder string) ([]PermissionEntity, error)
	GetAllPermissionsGroup(sortBy string, sortOrder string) ([]model.PermissionGroup, error)
	GetPermissionByID(id uint) (*PermissionEntity, error)
	GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
	GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
}

type permissionService struct {
	repo PermissionRepository
}

func NewPermissionService(repo PermissionRepository) PermissionService {
	return &permissionService{
		repo: repo,
	}
}

func (s *permissionService) GetAllPermissions(sortBy string, sortOrder string) ([]PermissionEntity, error) {
	return s.repo.GetAllPermissions(sortBy, sortOrder, nil)
}

func (s *permissionService) GetAllPermissionsGroup(sortBy string, sortOrder string) ([]model.PermissionGroup, error) {
	return s.repo.GetAllPermissionsGroup(sortBy, sortOrder, nil)
}

func (s *permissionService) GetPermissionByID(id uint) (*PermissionEntity, error) {
	return s.repo.GetPermissionByID(id)
}

func (s *permissionService) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	if userID <= 0 {
		return false, false, nil, false, errors.New("invalid user ID")
	}

	if permission == "" {
		return false, false, nil, false, errors.New("permission cannot be empty")
	}

	return s.repo.GetUserStatusAndPermission(userID, permission)
}

func (s *permissionService) GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	if userID == 0 {
		return false, false, nil, false, errors.New("invalid user ID")
	}

	if len(permissions) == 0 {
		return false, false, nil, false, errors.New("permissions cannot be empty")
	}

	return s.repo.GetUserStatusAndAnyPermission(userID, permissions)
}
