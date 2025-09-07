package auth

import (
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

func (r *RoleRepository) GetAllRoles(include []string, roleType *string, sortBy, sortOrder string) ([]Role, error) {
	var roles []Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if roleType != nil {
		query = query.Where(RoleRoleType+" = ?", *roleType)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, description, created_at, updated_at")
	}).Find(&roles).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch roles", "method", "GetAllRoles", "error", err, "include", include, "roleType", roleType, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return roles, err
}

func (r *RoleRepository) GetRoleByID(id uint, include []string) (*Role, error) {
	var role Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, description, created_at, updated_at")
	}).Where(constants.FieldID+" = ?", id).First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by ID", "method", "GetRoleByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetRoleByName(name string, include []string) (*Role, error) {
	var role Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, description, created_at, updated_at")
	}).Where(RoleRoleName+" = ?", name).First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by name", "method", "GetRoleByName", "error", err, "name", name, "include", include)
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetRoleByType(roleType string, include []string) (*Role, error) {
	var role Role

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, description, created_at, updated_at")
	}).Where(RoleRoleType+" = ?", roleType).First(&role).Error; err != nil {
		logger.Logger.Error("Failed to fetch role by type", "method", "GetRoleByType", "error", err, "roleType", roleType, "include", include)
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) CreateRole(role *Role) error {
	err := r.db.Create(role).Error
	if err != nil {
		logger.Logger.Error("Failed to create role", "method", "CreateRole", "error", err, "role", role)
	}
	return err
}

func (r *RoleRepository) UpdateRole(role *Role) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if err := tx.Model(role).Association("Permissions").Replace(role.Permissions); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Logger.Error("Failed to update role", "method", "UpdateRole", "error", err, "role", role)
	}
	return err
}

func (r *RoleRepository) DeleteRole(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Role{}).Where(constants.FieldID+" = ?", id).Association("Permissions").Clear(); err != nil {
			return err
		}

		if err := tx.Model(&Role{}).Where(constants.FieldID+" = ?", id).Association("Users").Clear(); err != nil {
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

func (r *RoleRepository) RoleExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&Role{}).Where(constants.FieldID+" = ?", id).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if role exists", "method", "RoleExists", "error", err, "id", id)
	}
	return count > 0, err
}

func (r *RoleRepository) RoleExistsByName(name string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&Role{}).Where(RoleRoleName+" = ?", name)

	if excludeID != nil {
		query = query.Where(constants.FieldID+" != ?", *excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if role exists by name", "method", "RoleExistsByName", "error", err, "name", name, "excludeID", excludeID)
	}
	return count > 0, err
}

func (r *RoleRepository) RoleExistsByType(roleType string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&Role{}).Where(RoleRoleType+" = ?", roleType)

	if excludeID != nil {
		query = query.Where(constants.FieldID+" != ?", *excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if role exists by type", "method", "RoleExistsByType", "error", err, "roleType", roleType, "excludeID", excludeID)
	}
	return count > 0, err
}

func (r *RoleRepository) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var permissions []Permission
		if err := tx.Select("id, name, description, created_at, updated_at").Where(constants.FieldID+" IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		var role Role
		if err := tx.Select("id").First(&role, roleID).Error; err != nil {
			return err
		}

		return tx.Model(&role).Association("Permissions").Replace(permissions)
	})
	if err != nil {
		logger.Logger.Error("Failed to assign permissions to role", "method", "AssignPermissionsToRole", "error", err, "roleID", roleID, "permissionIDs", permissionIDs)
	}
	return err
}

func (r *RoleRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, RoleRoleName, RoleRoleType, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{"description"}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
