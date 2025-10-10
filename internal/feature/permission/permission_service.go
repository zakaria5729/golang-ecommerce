package permission

import (
	"errors"
)

type PermissionService struct {
	repo *PermissionRepository
}

func NewPermissionService(repo *PermissionRepository) *PermissionService {
	return &PermissionService{
		repo: repo,
	}
}

func (s *PermissionService) GetAllPermissions(sortBy string, sortOrder string) ([]Permission, error) {
	return s.repo.GetAllPermissions(sortBy, sortOrder, nil)
}

func (s *PermissionService) GetPermissionByID(id uint) (*Permission, error) {
	return s.repo.GetPermissionByID(id)
}

func (s *PermissionService) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	if userID <= 0 {
		return false, false, nil, false, errors.New("invalid user ID")
	}

	if permission == "" {
		return false, false, nil, false, errors.New("permission cannot be empty")
	}

	return s.repo.GetUserStatusAndPermission(userID, permission)
}

func (s *PermissionService) GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	if userID == 0 {
		return false, false, nil, false, errors.New("invalid user ID")
	}

	if len(permissions) == 0 {
		return false, false, nil, false, errors.New("permissions cannot be empty")
	}

	return s.repo.GetUserStatusAndAnyPermission(userID, permissions)
}
