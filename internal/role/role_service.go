package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/internal/permission"
	"github.com/easy-comerce/backend/internal/role/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type RoleService interface {
	GetAllRoles(includeStr string, showDeletedStr string, roleTypeFilter string, sortBy, sortOrder string) ([]RoleEntity, error)
	GetRoleByID(id uint, includeStr string, showDeletedStr string) (*RoleEntity, error)
	CreateRole(ctx context.Context, req *model.CreateRoleRequest) (*RoleEntity, error)
	UpdateRole(ctx context.Context, id uint, req *model.UpdateRoleRequest) (*RoleEntity, error)
	DeleteRole(ctx context.Context, id uint) error
	UndoDeletedRole(ctx context.Context, id uint) error
	AssignRoleToUser(userID uint, req *model.AssignRoleRequest) error
	AddPermissionsToRole(roleID uint, req *model.AddPermissionsToRoleRequest) error
	GetRoleByType(roleType string) (*RoleEntity, error)
}

type roleService struct {
	roleRepo       RoleRepository
	permissionRepo permission.PermissionRepository
}

func NewRoleService(roleRepo RoleRepository, permissionRepo permission.PermissionRepository) RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *roleService) GetAllRoles(includeStr string, showDeletedStr string, roleTypeFilter string, sortBy, sortOrder string) ([]RoleEntity, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)

	if roleTypeFilter != "" {
		if !s.isValidRoleType(roleTypeFilter) {
			return nil, errors.New("invalid role type")
		}
	}

	roles, err := s.roleRepo.GetAllRoles(include, showDeleted, &roleTypeFilter, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	return roles, nil
}

func (s *roleService) GetRoleByID(id uint, includeStr string, showDeletedStr string) (*RoleEntity, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)

	role, err := s.roleRepo.GetRoleByID(id, include, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	return role, nil
}

func (s *roleService) CreateRole(ctx context.Context, req *model.CreateRoleRequest) (*RoleEntity, error) {
	req.Sanitize()

	if !s.isValidRoleType(req.RoleType) {
		return nil, errors.New("invalid role type")
	}

	if req.RoleType == c.RoleTypeSuperAdmin {
		return nil, errors.New("Super admin role already exists, you can't create role with this role type")
	}

	exists, err := s.roleRepo.RoleExistsByName(req.RoleName)
	if exists {
		return nil, errors.New("role with this name already exists")
	}

	role := &RoleEntity{
		RoleName:    req.RoleName,
		RoleType:    req.RoleType,
		Description: req.Description,
	}

	userID, _ := cu.GetUserIDFromContext(ctx)
	role.CreatedBy = userID

	permissions, err := s.permissionRepo.GetPermissionsByIDs(req.PermissionIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions: %w", err)
	}
	if len(permissions) == 0 {
		return nil, fmt.Errorf("no permission found with this permission_ids: %w", err)
	}
	role.Permissions = permissions

	if role, err := s.roleRepo.CreateRole(role); err != nil {
		logger.Logger.Error("Failed to create role", "method", "CreateRole", "error", err, "role", role)
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id uint, req *model.UpdateRoleRequest) (*RoleEntity, error) {
	req.Sanitize()

	existingRole, err := s.roleRepo.GetRoleByID(id, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	if existingRole.RoleType == c.RoleTypeSuperAdmin {
		return nil, errors.New("cannot update super admin role")
	}

	if req.RoleName != "" && req.RoleName != existingRole.RoleName {
		exists, err := s.roleRepo.RoleExistsByName(req.RoleName, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check role name: %w", err)
		}
		if exists {
			return nil, errors.New("role with this name already exists")
		}
		existingRole.RoleName = req.RoleName
	}

	if req.Description != nil {
		existingRole.Description = req.Description
	}

	userID, _ := cu.GetUserIDFromContext(ctx)
	existingRole.CreatedBy = userID

	if len(req.PermissionIDs) > 0 {
		permissions, err := s.permissionRepo.GetPermissionsByIDs(req.PermissionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch permissions: %w", err)
		}
		if len(permissions) == 0 {
			return nil, fmt.Errorf("no permission found with this permission_ids: %w", err)
		}
		existingRole.Permissions = permissions

		if err := s.roleRepo.UpdateRole(existingRole); err != nil {
			return nil, fmt.Errorf("failed to update role: %w", err)
		}
	} else {
		if err := s.roleRepo.UpdateRoleWithoutPermissions(existingRole); err != nil {
			return nil, fmt.Errorf("failed to update role without permissions: %w", err)
		}
	}

	return existingRole, nil
}

func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	role, err := s.roleRepo.GetRoleByID(id, nil, nil)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	if role.RoleType == c.RoleTypeSuperAdmin {
		return errors.New("cannot delete super admin role")
	}

	if err := s.roleRepo.DeleteRole(ctx, id); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

func (s *roleService) UndoDeletedRole(ctx context.Context, id uint) error {
	showDeleted := true
	exists, err := s.roleRepo.RoleExists(id, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("role not found: %w", err)
	}

	if err := s.roleRepo.UndoDeletedRole(ctx, id); err != nil {
		return fmt.Errorf("failed to undo deleted role: %w", err)
	}

	return nil
}

func (s *roleService) AssignRoleToUser(userID uint, req *model.AssignRoleRequest) error {
	if err := s.roleRepo.AssignRoleToUser(userID, req.RoleId); err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (s *roleService) AddPermissionsToRole(roleID uint, req *model.AddPermissionsToRoleRequest) error {
	if err := s.roleRepo.AddPermissionsToRoleByIds(roleID, req.PermissionIds, nil); err != nil {
		return fmt.Errorf("failed to add permissions to role: %w", err)
	}

	return nil
}

func (s *roleService) GetRoleByType(roleType string) (*RoleEntity, error) {
	role, err := s.roleRepo.GetRoleByType(roleType, nil)
	if err != nil {
		return nil, errors.New("role not found")
	}

	return role, nil
}

func (s *roleService) isValidRoleType(roleType string) bool {
	switch roleType {
	case c.RoleTypeSuperAdmin, c.RoleTypeAdmin, c.RoleTypeMaintainer, c.RoleTypeSeller, c.RoleTypeUser:
		return true
	default:
		return false
	}
}
