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

// **REQUIRED
func (r *PermissionRepository) GetAllPermissions(sortBy string, sortOrder string) ([]Permission, error) {
	var permissions []Permission
	query := r.db.Model(&Permission{})

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&permissions).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch permissions", "method", "GetAllPermissions", "error", err, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return permissions, err
}

// **REQUIRED
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

// **REQUIRED
func (r *PermissionRepository) GetPermissionsByIDs(ids []uint) ([]Permission, error) {
	var permissions []Permission

	if err := r.db.Model(&Permission{}).Where(constants.FieldID+" IN ?", ids).Find(&permissions).Error; err != nil {
		logger.Logger.Error("Failed to fetch permissions by IDs", "method", "GetPermissionsByIDs", "error", err, "ids", ids)
		return nil, err
	}

	return permissions, nil
}

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
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

// func (r *PermissionRepository) GetUsersWithPermission(permission string) ([]uint, error) {
// 	var userIDs []uint

// 	err := r.buildPermissionJoinQuery().
// 		Select("DISTINCT u.id").
// 		Where("p.name = ?", permission).
// 		Pluck("u.id", &userIDs).Error

// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch users with permission", "method", "GetUsersWithPermission", "error", err, "permission", permission)
// 	}

// 	return userIDs, err
// }

// func (r *PermissionRepository) GetPermissionStats() (map[string]int, error) {
// 	var results []struct {
// 		PermissionName string `json:"permission_name"`
// 		Count          int    `json:"count"`
// 	}

// 	err := r.buildPermissionJoinQuery().
// 		Select("p.name as permission_name, COUNT(DISTINCT u.id) as count").
// 		Group("p.name").
// 		Scan(&results).Error

// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch permission stats", "method", "GetPermissionStats", "error", err)
// 		return nil, err
// 	}

// 	stats := make(map[string]int)
// 	for _, result := range results {
// 		stats[result.PermissionName] = result.Count
// 	}

// 	return stats, nil
// }

// func (r *PermissionRepository) CreatePermission(permission *Permission) error {
// 	err := r.db.Create(permission).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to create permission", "method", "CreatePermission", "error", err, "permission", permission)
// 	}
// 	return err
// }

// func (r *PermissionRepository) UpdatePermission(permission *Permission) error {
// 	err := r.db.Save(permission).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to update permission", "method", "UpdatePermission", "error", err, "permission", permission)
// 	}
// 	return err
// }

// func (r *PermissionRepository) DeletePermission(id uint) error {
// 	err := r.db.Delete(&Permission{}, id).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to delete permission", "method", "DeletePermission", "error", err, "id", id)
// 	}
// 	return err
// }

// func (r *PermissionRepository) PermissionExists(id uint) (bool, error) {
// 	var count int64
// 	err := r.db.Model(&Permission{}).Where(constants.FieldID+" = ?", id).Count(&count).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to check if permission exists", "method", "PermissionExists", "error", err, "id", id)
// 	}
// 	return count > 0, err
// }

// func (r *PermissionRepository) PermissionExistsByName(name string, excludeID *uint) (bool, error) {
// 	var count int64
// 	query := r.db.Model(&Permission{}).Where(constants.PermissionName+" = ?", name)

// 	if excludeID != nil {
// 		query = query.Where(constants.FieldID+" != ?", *excludeID)
// 	}

// 	err := query.Count(&count).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to check if permission exists by name", "method", "PermissionExistsByName", "error", err, "name", name, "excludeID", excludeID)
// 	}
// 	return count > 0, err
// }

// func (r *PermissionRepository) GetUserPermissions(userID uint) ([]string, error) {
// 	var permissions []string

// 	err := r.buildPermissionJoinQuery().
// 		Select("DISTINCT p.name").
// 		Where("u.id = ?", userID).
// 		Pluck("p.name", &permissions).Error

// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch user permissions", "method", "GetUserPermissions", "error", err, "userID", userID)
// 	}

// 	return permissions, err
// }

// // **REQUIRED
// func (r *PermissionRepository) GetUserPermissionsByRole(userID uint, roleType string) ([]string, error) {
// 	var permissions []string

// 	err := r.buildPermissionJoinQuery().
// 		Select("DISTINCT p.name").
// 		Where("u.id = ? AND r.role_type = ?", userID, roleType).
// 		Pluck("p.name", &permissions).Error

// 	if err != nil {
// 		logger.Logger.Error("Failed to fetch user permissions by role", "method", "GetUserPermissionsByRole", "error", err, "userID", userID, "roleType", roleType)
// 	}

// 	return permissions, err
// }
