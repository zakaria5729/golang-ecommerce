package role

import (
	"os/user"
	"strings"

	"github.com/easy-comerce/backend/db"
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

// **REQUIRED
func (r *RoleRepository) GetAllRoles(include []string, roleType *string, sortBy, sortOrder string) ([]Role, error) {
	var roles []Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if roleType != nil {
		query = query.Where(constants.RoleRoleType+" = ?", *roleType)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if utils.ContainsString(include, constants.RolePermissions) {
		query = query.Preload(constants.RolePermissionsCapitalized)
	}

	err := query.Find(&roles).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch roles", "method", "GetAllRoles", "error", err, "include", include, "roleType", roleType, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return roles, err
}

// **REQUIRED
func (r *RoleRepository) GetRoleByID(id uint, include []string) (*Role, error) {
	var role Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if utils.ContainsString(include, constants.RolePermissions) {
		query = query.Preload(constants.RolePermissionsCapitalized)
	}

	if err := query.Where(constants.FieldID+" = ?", id).First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by ID", "method", "GetRoleByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &role, nil
}

// **REQUIRED
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

// **REQUIRED
func (r *RoleRepository) GetRoleByType(roleType string) (*Role, error) {
	var role Role

	if err := r.db.Model(&Role{}).Preload(constants.RolePermissionsCapitalized).
		Where(constants.RoleRoleType+" = ?", roleType).
		First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by type", "method", "GetRoleByType", "error", err, "roleType", roleType)
		return nil, err
	}

	return &role, nil
}

// **REQUIRED
func (r *RoleRepository) CreateRole(role *Role) error {
	err := r.db.Create(role).Error
	if err != nil {
		logger.Logger.Error("Failed to create role", "method", "CreateRole", "error", err, "role", role)
	}
	return err
}

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
func (r *RoleRepository) DeleteRole(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Role{}).Where(constants.FieldID+" = ?", id).Association(constants.RolePermissionsCapitalized).Clear(); err != nil {
			return err
		}

		if err := tx.Delete(&Role{}, id).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error("Failed to delete role", "method", "DeleteRole", "error", err, "id", id)
	}

	return err
}

// **REQUIRED
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

// **REQUIRED
func (r *RoleRepository) RoleExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&Role{}).Where(constants.FieldID+" = ?", id).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if role exists", "method", "RoleExists", "error", err, "id", id)
	}
	return count > 0, err
}

// **REQUIRED
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

// **REQUIRED
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

// func (r *RoleRepository) GetRoleByName(name string, include []string) (*Role, error) {
// 	var role Role
// 	selectFields := r.getSelectableFields(include)
// 	query := r.db.Select(strings.Join(selectFields, ", "))

// 	if err := query.Preload(RolePermissionsCapitalized).
// 		Where(RoleRoleName+" = ?", name).
// 		First(&role).Error; err != nil {
// 		logger.Logger.Error("Failed to fetch role by name", "method", "GetRoleByName", "error", err, "name", name, "include", include)
// 		return nil, err
// 	}

// 	return &role, nil
// }

// func (r *RoleRepository) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
// 	err := r.db.Transaction(func(tx *gorm.DB) error {
// 		var permissions []permission.Permission
// 		if err := tx.Select("id, name, description, created_at, updated_at").Where(constants.FieldID+" IN ?", permissionIDs).Find(&permissions).Error; err != nil {
// 			return err
// 		}

// 		var role Role
// 		if err := tx.Select("id").First(&role, roleID).Error; err != nil {
// 			return err
// 		}

// 		return tx.Model(&role).Association("Permissions").Replace(permissions)
// 	})
// 	if err != nil {
// 		logger.Logger.Error("Failed to assign permissions to role", "method", "AssignPermissionsToRole", "error", err, "roleID", roleID, "permissionIDs", permissionIDs)
// 	}
// 	return err
// }
