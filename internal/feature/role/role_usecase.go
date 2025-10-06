package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/internal/feature/permission"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type RoleUseCase struct {
	roleRepo       *RoleRepository
	permissionRepo *permission.PermissionRepository
}

func NewRoleUseCase() *RoleUseCase {
	return &RoleUseCase{
		roleRepo:       NewRoleRepository(),
		permissionRepo: permission.NewPermissionRepository(),
	}
}

func (uc *RoleUseCase) GetAllRoles(includeStr string, showDeletedStr string, roleTypeFilter string, sortBy, sortOrder string) ([]Role, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)

	if roleTypeFilter != "" {
		if !isValidRoleType(roleTypeFilter) {
			return nil, errors.New("invalid role type")
		}
	}

	roles, err := uc.roleRepo.GetAllRoles(include, showDeleted, &roleTypeFilter, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch roles", "method", "GetAllRoles", "error", err, "include", include, "roleType", roleTypeFilter, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	return roles, nil
}

func (uc *RoleUseCase) GetRoleByID(id uint, includeStr string, showDeletedStr string) (*Role, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)
	role, err := uc.roleRepo.GetRoleByID(id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "GetRoleByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("role not found: %w", err)
	}

	return role, nil
}

func (uc *RoleUseCase) CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error) {
	req.Sanitize()

	if !isValidRoleType(req.RoleType) {
		logger.Logger.Error("Invalid role type", "method", "CreateRole", "roleType", req.RoleType)
		return nil, errors.New("invalid role type")
	}

	if req.RoleType == c.RoleTypeSuperAdmin {
		logger.Logger.Error("Role type cannot create roles", "method", "CreateRole", "roleType", req.RoleType)
		return nil, errors.New("Super admin role already exists, you can't create role with this role type")
	}

	exists, err := uc.roleRepo.RoleExistsByName(req.RoleName)
	if exists {
		logger.Logger.Error("Role already exists", "method", "CreateRole", "roleName", req.RoleName)
		return nil, errors.New("role with this name already exists")
	}

	role := &Role{
		RoleName:    req.RoleName,
		RoleType:    req.RoleType,
		Description: req.Description,
	}

	userID, ok := ctx.Value(c.UserIDContextKey).(uint)
	if ok {
		role.CreatedBy = &userID
	}

	permissions, err := uc.permissionRepo.GetPermissionsByIDs(req.PermissionIDs)
	if err != nil {
		logger.Logger.Error("Failed to fetch permissions", "method", "CreateRole", "error", err, "PermissionIDs", req.PermissionIDs)
		return nil, fmt.Errorf("failed to fetch permissions: %w", err)
	}
	if len(permissions) == 0 {
		logger.Logger.Error("Failed to fetch permissions", "method", "CreateRole", "error", err, "PermissionIDs", req.PermissionIDs)
		return nil, fmt.Errorf("no permission found with this permission_ids: %w", err)
	}
	role.Permissions = permissions

	if role, err := uc.roleRepo.CreateRole(role); err != nil {
		logger.Logger.Error("Failed to create role", "method", "CreateRole", "error", err, "role", role)
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

func (uc *RoleUseCase) UpdateRole(ctx context.Context, id uint, req *UpdateRoleRequest) (*Role, error) {
	req.Sanitize()

	existingRole, err := uc.roleRepo.GetRoleByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "UpdateRole", "error", err, "id", id)
		return nil, fmt.Errorf("role not found: %w", err)
	}

	if existingRole.RoleType == c.RoleTypeSuperAdmin {
		logger.Logger.Error("Cannot update super admin role", "method", "UpdateRole", "id", id)
		return nil, errors.New("cannot update super admin role")
	}

	if req.RoleName != "" && req.RoleName != existingRole.RoleName {
		exists, err := uc.roleRepo.RoleExistsByName(req.RoleName, id)
		if err != nil {
			logger.Logger.Error("Failed to check if role name exists", "method", "UpdateRole", "error", err, "id", id, "roleName", req.RoleName)
			return nil, fmt.Errorf("failed to check role name: %w", err)
		}
		if exists {
			logger.Logger.Error("Role name already exists", "method", "UpdateRole", "id", id, "roleName", req.RoleName)
			return nil, errors.New("role with this name already exists")
		}
		existingRole.RoleName = req.RoleName
	}

	if req.Description != nil {
		existingRole.Description = req.Description
	}

	userID, ok := ctx.Value(c.UserIDContextKey).(uint)
	if ok {
		existingRole.UpdatedBy = &userID
	}

	if len(req.PermissionIDs) > 0 {
		permissions, err := uc.permissionRepo.GetPermissionsByIDs(req.PermissionIDs)
		if err != nil {
			logger.Logger.Error("Failed to fetch permissions", "method", "UpdateRole", "error", err, "id", id, "permissions", req.PermissionIDs)
			return nil, fmt.Errorf("failed to fetch permissions: %w", err)
		}
		if len(permissions) == 0 {
			logger.Logger.Error("Failed to fetch permissions", "method", "UpdateRole", "error", err, "id", id, "permissions", req.PermissionIDs)
			return nil, fmt.Errorf("no permission found with this permission_ids: %w", err)
		}
		existingRole.Permissions = permissions

		if err := uc.roleRepo.UpdateRole(existingRole); err != nil {
			logger.Logger.Error("Failed to update role", "method", "UpdateRole", "error", err, "role", existingRole)
			return nil, fmt.Errorf("failed to update role: %w", err)
		}
	} else {
		if err := uc.roleRepo.UpdateRoleWithoutPermissions(existingRole); err != nil {
			logger.Logger.Error("Failed to update role without permissions", "method", "UpdateRole", "error", err, "role", existingRole)
			return nil, fmt.Errorf("failed to update role without permissions: %w", err)
		}
	}

	return existingRole, nil
}

func (uc *RoleUseCase) DeleteRole(ctx context.Context, id uint) error {
	role, err := uc.roleRepo.GetRoleByID(id, nil, nil)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "DeleteRole", "error", err, "id", id)
		return fmt.Errorf("role not found: %w", err)
	}

	if role.RoleType == c.RoleTypeSuperAdmin {
		logger.Logger.Error("Cannot delete super admin role", "method", "DeleteRole", "id", id)
		return errors.New("cannot delete super admin role")
	}

	if err := uc.roleRepo.DeleteRole(ctx, id); err != nil {
		logger.Logger.Error("Failed to delete role", "method", "DeleteRole", "error", err, "id", id)
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

func (uc *RoleUseCase) UndoDeletedRole(ctx context.Context, id uint) error {
	showDeleted := true
	exists, err := uc.roleRepo.RoleExists(id, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("Role not found", "method", "UndoDeletedRole", "error", err, "id", id)
		return fmt.Errorf("role not found: %w", err)
	}

	if err := uc.roleRepo.UndoDeletedRole(ctx, id); err != nil {
		logger.Logger.Error("Failed to undo deleted role", "method", "UndoDeletedRole", "error", err, "id", id)
		return fmt.Errorf("failed to undo deleted role: %w", err)
	}

	return nil
}

func (uc *RoleUseCase) AssignRoleToUser(userID uint, req *AssignRoleRequest) error {
	if err := uc.roleRepo.AssignRoleToUser(userID, req.RoleId); err != nil {
		logger.Logger.Error("Failed to assign role to user", "method", "AssignRoleToUser", "error", err, "userID", userID, "roleId", req.RoleId)
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (uc *RoleUseCase) AddPermissionsToRole(roleID uint, req *AddPermissionsToRoleRequest) error {
	if err := uc.roleRepo.AddPermissionsToRoleByIds(roleID, req.PermissionIds, nil); err != nil {
		logger.Logger.Error("Failed to add permissions to role", "method", "AddPermissionsToRole", "error", err, "roleID", roleID, "permissionNames", req.PermissionIds)
		return fmt.Errorf("failed to add permissions to role: %w", err)
	}

	return nil
}

func (uc *RoleUseCase) GetRoleByType(roleType string) (*Role, error) {
	role, err := uc.roleRepo.GetRoleByType(roleType, nil)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "GetRoleByType", "error", err, "roleType", roleType)
		return nil, errors.New("role not found")
	}

	return role, nil
}

func isValidRoleType(roleType string) bool {
	switch roleType {
	case c.RoleTypeSuperAdmin, c.RoleTypeAdmin, c.RoleTypeMaintainer, c.RoleTypeSeller, c.RoleTypeUser:
		return true
	default:
		return false
	}
}
