package role

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/pkg/constants"
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

// **REQUIRED
func (uc *RoleUseCase) GetAllRoles(includeStr string, roleTypeFilter string, sortBy, sortOrder string) ([]Role, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	var roleType *string
	if roleTypeFilter != "" {
		if !isValidRoleType(roleTypeFilter) {
			return nil, errors.New("invalid role type")
		}
		roleType = &roleTypeFilter
	}

	roles, err := uc.roleRepo.GetAllRoles(include, roleType, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch roles", "method", "GetAllRoles", "error", err, "include", include, "roleType", roleType, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	return roles, nil
}

// **REQUIRED
func (uc *RoleUseCase) GetRoleByID(id uint, includeStr string) (*Role, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	role, err := uc.roleRepo.GetRoleByID(id, include)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "GetRoleByID", "error", err, "id", id, "include", include)
		return nil, fmt.Errorf("role not found: %w", err)
	}

	return role, nil
}

// **REQUIRED
func (uc *RoleUseCase) CreateRole(req *CreateRoleRequest) (*Role, error) {
	req.Sanitize()

	if !isValidRoleType(req.RoleType) {
		logger.Logger.Error("Invalid role type", "method", "CreateRole", "roleType", req.RoleType)
		return nil, errors.New("invalid role type")
	}

	if !canCreateRoles(req.RoleType) {
		logger.Logger.Error("Role type cannot create roles", "method", "CreateRole", "roleType", req.RoleType)
		return nil, errors.New("this role type cannot create roles")
	}

	exists, err := uc.roleRepo.RoleExistsByName(req.RoleName)
	if err != nil {
		logger.Logger.Error("Failed to check if role exists", "method", "CreateRole", "error", err, "roleName", req.RoleName)
		return nil, fmt.Errorf("failed to check role existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Role already exists", "method", "CreateRole", "roleName", req.RoleName)
		return nil, errors.New("role with this name already exists")
	}

	exists, err = uc.roleRepo.RoleExistsByType(req.RoleType)
	if err != nil {
		logger.Logger.Error("Failed to check if role type exists", "method", "CreateRole", "error", err, "roleType", req.RoleType)
		return nil, fmt.Errorf("failed to check role type existence: %w", err)
	}
	if exists {
		logger.Logger.Error("Role type already exists", "method", "CreateRole", "roleType", req.RoleType)
		return nil, errors.New("role with this type already exists")
	}

	role := &Role{
		RoleName:    req.RoleName,
		RoleType:    req.RoleType,
		Description: req.Description,
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

	if err := uc.roleRepo.CreateRole(role); err != nil {
		logger.Logger.Error("Failed to create role", "method", "CreateRole", "error", err, "role", role)
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// **REQUIRED
func (uc *RoleUseCase) UpdateRole(id uint, req *UpdateRoleRequest) (*Role, error) {
	req.Sanitize()

	existingRole, err := uc.roleRepo.GetRoleByID(id, nil)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "UpdateRole", "error", err, "id", id)
		return nil, fmt.Errorf("role not found: %w", err)
	}

	if existingRole.RoleType == constants.RoleTypeSuperAdmin {
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

// **REQUIRED
func (uc *RoleUseCase) DeleteRole(id uint) error {
	role, err := uc.roleRepo.GetRoleByID(id, nil)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "DeleteRole", "error", err, "id", id)
		return fmt.Errorf("role not found: %w", err)
	}

	if role.RoleType == constants.RoleTypeSuperAdmin {
		logger.Logger.Error("Cannot delete super admin role", "method", "DeleteRole", "id", id)
		return errors.New("cannot delete super admin role")
	}

	if err := uc.roleRepo.DeleteRole(id); err != nil {
		logger.Logger.Error("Failed to delete role", "method", "DeleteRole", "error", err, "id", id)
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *RoleUseCase) AssignRoleToUser(userID uint, req *AssignRoleRequest) error {
	if err := uc.roleRepo.AssignRoleToUser(userID, req.RoleId); err != nil {
		logger.Logger.Error("Failed to assign role to user", "method", "AssignRoleToUser", "error", err, "userID", userID, "roleId", req.RoleId)
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *RoleUseCase) GetRoleByType(roleType string) (*Role, error) {
	role, err := uc.roleRepo.GetRoleByType(roleType)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "GetRoleByType", "error", err, "roleType", roleType)
		return nil, errors.New("role not found")
	}

	return role, nil
}

// **REQUIRED
func isValidRoleType(roleType string) bool {
	switch roleType {
	case constants.RoleTypeSuperAdmin, constants.RoleTypeAdmin, constants.RoleTypeManager, constants.RoleTypeSeller, constants.RoleTypeUser:
		return true
	default:
		return false
	}
}

// **REQUIRED
func canCreateRoles(roleType string) bool {
	return roleType == constants.RoleTypeSuperAdmin || roleType == constants.RoleTypeAdmin
}

// func (uc *RoleUseCase) GetAllPermissions(includeStr string, sortBy, sortOrder string) ([]permission.Permission, error) {
// 	include := utils.ParseCommaSeparatedString(includeStr)
// 	permissions, err := uc.permissionRepo.GetAllPermissions(include, sortBy, sortOrder)
// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch permissions", "method", "GetAllPermissions", "error", err, "include", include, "sortBy", sortBy, "sortOrder", sortOrder)
// 		return nil, fmt.Errorf("failed to fetch permissions: %w", err)
// 	}

// 	return permissions, nil
// }

// func (uc *RoleUseCase) GetPermissionByID(id uint, includeStr string) (*permission.Permission, error) {
// 	include := utils.ParseCommaSeparatedString(includeStr)
// 	permission, err := uc.permissionRepo.GetPermissionByID(id, include)
// 	if err != nil {
// 		logger.Logger.Error("Permission not found", "method", "GetPermissionByID", "error", err, "id", id, "include", include)
// 		return nil, fmt.Errorf("permission not found: %w", err)
// 	}

// 	return permission, nil
// }

// func (uc *RoleUseCase) CreatePermission(permission *permission.Permission) (*permission.Permission, error) {
// 	permission.Sanitize()

// 	exists, err := uc.permissionRepo.PermissionExistsByName(permission.Name, nil)
// 	if err != nil {
// 		logger.Logger.Error("Failed to check if permission exists", "method", "CreatePermission", "error", err, "name", permission.Name)
// 		return nil, fmt.Errorf("failed to check permission existence: %w", err)
// 	}
// 	if exists {
// 		logger.Logger.Error("Permission already exists", "method", "CreatePermission", "name", permission.Name)
// 		return nil, errors.New("permission with this name already exists")
// 	}

// 	if err := uc.permissionRepo.CreatePermission(permission); err != nil {
// 		logger.Logger.Error("Failed to create permission", "method", "CreatePermission", "error", err, "permission", permission)
// 		return nil, fmt.Errorf("failed to create permission: %w", err)
// 	}

// 	return permission, nil
// }

// func (uc *RoleUseCase) UpdatePermission(id uint, permission *permission.Permission) (*permission.Permission, error) {
// 	permission.Sanitize()

// 	existingPermission, err := uc.permissionRepo.GetPermissionByID(id, nil)
// 	if err != nil {
// 		logger.Logger.Error("Permission not found", "method", "UpdatePermission", "error", err, "id", id)
// 		return nil, fmt.Errorf("permission not found: %w", err)
// 	}

// 	if permission.Name != "" && permission.Name != existingPermission.Name {
// 		exists, err := uc.permissionRepo.PermissionExistsByName(permission.Name, &id)
// 		if err != nil {
// 			logger.Logger.Error("Failed to check if permission name exists", "method", "UpdatePermission", "error", err, "id", id, "name", permission.Name)
// 			return nil, fmt.Errorf("failed to check permission name: %w", err)
// 		}
// 		if exists {
// 			logger.Logger.Error("Permission name already exists", "method", "UpdatePermission", "id", id, "name", permission.Name)
// 			return nil, errors.New("permission with this name already exists")
// 		}
// 		existingPermission.Name = permission.Name
// 	}

// 	if permission.Description != nil {
// 		existingPermission.Description = permission.Description
// 	}

// 	if err := uc.permissionRepo.UpdatePermission(existingPermission); err != nil {
// 		logger.Logger.Error("Failed to update permission", "method", "UpdatePermission", "error", err, "permission", existingPermission)
// 		return nil, fmt.Errorf("failed to update permission: %w", err)
// 	}

// 	return existingPermission, nil
// }

// func (uc *RoleUseCase) DeletePermission(id uint) error {
// 	_, err := uc.permissionRepo.GetPermissionByID(id, nil)
// 	if err != nil {
// 		logger.Logger.Error("Permission not found", "method", "DeletePermission", "error", err, "id", id)
// 		return fmt.Errorf("permission not found: %w", err)
// 	}

// 	if err := uc.permissionRepo.DeletePermission(id); err != nil {
// 		logger.Logger.Error("Failed to delete permission", "method", "DeletePermission", "error", err, "id", id)
// 		return fmt.Errorf("failed to delete permission: %w", err)
// 	}

// 	return nil
// }
