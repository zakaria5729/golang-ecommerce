package role

import (
	"os/user"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		db: db.GetDB(),
	}
}

func (r *RoleRepository) GetAllRoles(include []string, showDeleted *bool, roleType *string, sortBy, sortOrder string) ([]Role, error) {
	var roles []Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if roleType != nil && *roleType != "" {
		query = query.Where(constants.RoleRoleType+" = ?", *roleType)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if utils.ContainsString(include, constants.RolePermissions) {
		query = query.Preload(constants.RolePermissionsCapitalized)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&roles).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch roles", "method", "GetAllRoles", "error", err, "include", include, "roleType", roleType, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return roles, err
}

func (r *RoleRepository) GetRoleByID(id uint, include []string, showDeleted *bool) (*Role, error) {
	var role Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if utils.ContainsString(include, constants.RolePermissions) {
		query = query.Preload(constants.RolePermissionsCapitalized)
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by ID", "method", "GetRoleByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetRoleWithPermissionsByType(roleType string) (*Role, error) {
	var role Role

	if err := r.db.Model(&Role{}).
		Preload(constants.RolePermissionsCapitalized).
		Where(constants.RoleRoleType+" = ?", roleType).
		First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by type", "method", "GetRoleWithPermissionsByType", "error", err, "roleType", roleType)
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetRoleByType(roleType string, showDeleted *bool) (*Role, error) {
	var role Role

	query := r.db.Model(&Role{}).Where(constants.RoleRoleType+" = ?", roleType)
	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Preload(constants.RolePermissionsCapitalized).
		First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by type", "method", "GetRoleByType", "error", err, "roleType", roleType)
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) CreateRole(role *Role) (*Role, error) {
	err := r.db.Create(role).Error
	if err != nil {
		logger.Logger.Error("Failed to create role", "method", "CreateRole", "error", err, "role", role)
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) UpdateRole(role *Role) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if err := tx.Model(role).Association(constants.RolePermissionsCapitalized).Replace(role.Permissions); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error("Failed to update role", "method", "UpdateRole", "error", err, "role", role)
	}
	return err
}

func (r *RoleRepository) UpdateRoleWithoutPermissions(role *Role) error {
	err := r.db.Model(&Role{}).Where(constants.FieldID+" = ?", role.ID).Updates(map[string]any{
		constants.RoleRoleName:    role.RoleName,
		constants.RoleRoleType:    role.RoleType,
		constants.RoleDescription: role.Description,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to update role", "method", "UpdateRoleWithoutPermissions", "error", err, "role", role)
	}
	return err
}

func (r *RoleRepository) DeleteRole(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&Role{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete role", "method", "DeleteRole", "error", err, "id", id)
	}

	return err
}

func (r *RoleRepository) UndoDeletedRole(id uint) error {
	err := r.db.Unscoped().Model(&Role{}).Where(constants.FieldID+" = ?", id).Update(constants.FieldDeletedAt, nil).Error

	if err != nil {
		logger.Logger.Error("Failed to undo deleted role", "method", "UndoDeletedRole", "error", err, "id", id)
	}

	return err
}

func (r *RoleRepository) AssignRoleToUser(userID uint, roleID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Where(constants.FieldID+" = ?", roleID).First(&role).Error; err != nil {
			return err
		}

		var user user.User
		if err := tx.Select(constants.FieldID).First(&user, userID).Error; err != nil {
			return err
		}

		return tx.Model(&user).Association(constants.UserRolesCapitalized).Replace(&role)
	})

	if err != nil {
		logger.Logger.Error("Failed to assign roles to user", "method", "AssignRolesToUser", "error", err, "userID", userID, "roleID", roleID)
	}
	return err
}

func (r *RoleRepository) AddPermissionsToRole(roleID uint, permissionNames []string, showDeleted *bool) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var role Role
		query := tx.Where(constants.FieldID+" = ?", roleID)

		if showDeleted != nil && *showDeleted {
			query = query.Unscoped()
		}
		if err := query.Preload(constants.RolePermissionsCapitalized).First(&role).Error; err != nil {
			return err
		}

		var permissions []permission.Permission
		query = tx.Where(constants.PermissionName+" IN ?", permissionNames)

		if showDeleted != nil && *showDeleted {
			query = query.Unscoped()
		}
		if err := query.Find(&permissions).Error; err != nil {
			return err
		}

		existingPermissionNames := make(map[string]bool)
		for _, perm := range role.Permissions {
			existingPermissionNames[perm.Name] = true
		}

		var newPermissions []permission.Permission
		for _, perm := range permissions {
			if !existingPermissionNames[perm.Name] {
				newPermissions = append(newPermissions, perm)
			}
		}

		if len(newPermissions) > 0 {
			query = tx.Model(&role)
			if showDeleted != nil && *showDeleted {
				query = query.Unscoped()
			}
			return query.Association(constants.RolePermissionsCapitalized).Append(newPermissions)
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error("Failed to add permissions to role", "method", "AddPermissionsToRole", "error", err, "roleID", roleID, "permissionNames", permissionNames)
	}
	return err
}

func (r *RoleRepository) RoleExists(id uint, showDeleted *bool) (bool, error) {
	var count int64
	query := r.db.Model(&Role{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if role exists", "method", "RoleExists", "error", err, "id", id)
	}
	return count > 0, err
}

func (r *RoleRepository) RoleExistsByName(name string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&Role{}).Where(constants.RoleRoleName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(constants.FieldID+" != ?", excludeID[0])
	}

	err := r.db.Model(&Role{}).Where(constants.RoleRoleName+" = ?", name).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if role exists by name", "method", "RoleExistsByName", "error", err, "name", name)
	}
	return count > 0, err
}

func (r *RoleRepository) RoleExistsByType(roleType string) (bool, error) {
	var count int64

	err := r.db.Model(&Role{}).Where(constants.RoleRoleType+" = ?", roleType).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if role exists by type", "method", "RoleExistsByType", "error", err, "roleType", roleType)
	}
	return count > 0, err
}

func (r *RoleRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.RoleRoleName, constants.RoleRoleType, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{constants.RoleDescription}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
