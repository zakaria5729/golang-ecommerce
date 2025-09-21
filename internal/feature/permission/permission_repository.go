package permission

import (
	"errors"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{
		db: db.GetDB(),
	}
}

func (r *PermissionRepository) CreatePermissionsIfNotExists(permissionNames []string, showDeleted *bool) error {
	var existingPermissions []Permission
	query := r.db.Model(&Permission{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(constants.PermissionName+" IN ?", permissionNames).Find(&existingPermissions).Error; err != nil {
		logger.Logger.Error("Failed to fetch existing permissions", "method", "CreatePermissionsIfNotExists", "error", err)
		return err
	}

	existingNamesMap := make(map[string]bool, len(existingPermissions))
	for _, perm := range existingPermissions {
		existingNamesMap[perm.Name] = true
	}

	var toCreatePermissions []Permission
	for _, name := range permissionNames {
		if !existingNamesMap[name] {
			desc := constants.PermissionMap[name]

			toCreatePermissions = append(
				toCreatePermissions,
				Permission{
					Name:        name,
					Description: &desc,
				},
			)
		}
	}

	if len(toCreatePermissions) > 0 {
		if err := r.db.Create(&toCreatePermissions).Error; err != nil {
			logger.Logger.Error("Failed to create permissions", "method", "CreatePermissionsIfNotExists", "error", err)
			return err
		}
	}

	return nil
}

func (r *PermissionRepository) ExistsByName(permissionName string) (exists bool, err error) {
	var count int64
	if err := r.db.Model(&Permission{}).Where(constants.PermissionName+" = ?", permissionName).Count(&count).Error; err != nil {
		logger.Logger.Error("Failed to check if permission exists", "method", "ExistsByName", "error", err, "permissionName", permissionName)
		return false, err
	}
	return count > 0, nil
}

func (r *PermissionRepository) GetAllPermissions(sortBy string, sortOrder string, showDeleted *bool) ([]Permission, error) {
	var permissions []Permission
	query := r.db.Model(&Permission{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&permissions).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch permissions", "method", "GetAllPermissions", "error", err, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return permissions, err
}

func (r *PermissionRepository) GetPermissionByID(id uint) (*Permission, error) {
	var permission Permission

	if err := r.db.Model(&Permission{}).Where(constants.FieldID+" = ?", id).First(&permission).Error; err != nil {
		logger.Logger.Error("Failed to fetch permission by ID", "method", "GetPermissionByID", "error", err, "id", id)
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) GetPermissionByName(name string, include []string) (*Permission, error) {
	var permission Permission

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if err := query.Where(constants.PermissionName+" = ?", name).First(&permission).Error; err != nil {
		logger.Logger.Error("Failed to fetch permission by name", "method", "GetPermissionByName", "error", err, "name", name, "include", include)
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) GetPermissionsByIDs(ids []uint) ([]Permission, error) {
	var permissions []Permission

	if err := r.db.Model(&Permission{}).Where(constants.FieldID+" IN ?", ids).Find(&permissions).Error; err != nil {
		logger.Logger.Error("Failed to fetch permissions by IDs", "method", "GetPermissionsByIDs", "error", err, "ids", ids)
		return nil, err
	}

	return permissions, nil
}

func (r *PermissionRepository) HasPermission(userID uint, permission string) (bool, error) {
	var count int64

	err := r.buildPermissionJoinQuery().
		Where(constants.FieldID+" = ? AND "+constants.PermissionName+" = ?", userID, permission).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to check user permission", "method", "HasPermission", "error", err, "userID", userID, "permission", permission)
	}

	return count > 0, err
}

func (r *PermissionRepository) HasAnyPermission(userID uint, permissions []string) (bool, error) {
	if len(permissions) == 0 {
		return true, nil
	}

	var count int64

	err := r.buildPermissionJoinQuery().
		Where(constants.FieldID+" = ? AND "+constants.PermissionName+" IN ?", userID, permissions).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to check user permissions", "method", "HasAnyPermission", "error", err, "userID", userID, "permissions", permissions)
	}

	return count > 0, err
}

func (r *PermissionRepository) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, hasPermission bool, err error) {
	var userStatus PermissionUserStatus

	err = r.db.Select(constants.UserBanned, constants.UserVerified).
		Where("id = ?", userID).
		First(&userStatus).Error

	if err != nil {
		logger.Logger.Error("Failed to get user status", "method", "GetUserStatusAndPermission", "error", err, "userID", userID)
		return false, false, false, err
	}

	if !userStatus.Verified {
		return false, false, false, errors.New("account is not verified yet")
	}

	if userStatus.Banned {
		return false, false, false, errors.New("account is banned")
	}

	hasPermission, err = r.HasPermission(userID, permission)
	if err != nil {
		logger.Logger.Error("Failed to check user permission", "method", "GetUserStatusAndPermission", "error", err, "userID", userID, "permission", permission)
		return false, false, false, err
	}

	return userStatus.Banned, userStatus.Verified, hasPermission, nil
}

func (r *PermissionRepository) GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, hasPermission bool, err error) {
	var userStatus PermissionUserStatus

	err = r.db.Select(constants.UserBanned, constants.UserVerified).
		Where("id = ?", userID).
		First(&userStatus).Error

	if err != nil {
		logger.Logger.Error("Failed to get user status", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID)
		return false, false, false, err
	}

	hasPermission, err = r.HasAnyPermission(userID, permissions)
	if err != nil {
		logger.Logger.Error("Failed to check user permissions", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID, "permissions", permissions)
		return false, false, false, err
	}

	return userStatus.Banned, userStatus.Verified, hasPermission, nil
}

func (r *PermissionRepository) buildPermissionJoinQuery() *gorm.DB {
	return r.db.Table(constants.TableUser + " u").
		Joins("JOIN " + constants.TableUserRole + " ur ON u.id = ur.user_id").
		Joins("JOIN " + constants.TableRole + " r ON ur.role_id = r.id").
		Joins("JOIN " + constants.TableRolePermission + " rp ON r.id = rp.role_id").
		Joins("JOIN " + constants.TablePermission + " p ON rp.permission_id = p.id").
		Where("u.deleted_at IS NULL AND r.deleted_at IS NULL AND p.deleted_at IS NULL")
}

func (r *PermissionRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.PermissionName, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{constants.PermissionDescription}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
