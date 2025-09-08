package user_permission

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/logger"
	"gorm.io/gorm"
)

type UserPermissionRepository struct {
	db *gorm.DB
}

func NewUserPermissionRepository() *UserPermissionRepository {
	return &UserPermissionRepository{
		db: db.GetDB(),
	}
}

func (r *UserPermissionRepository) GetUserPermissions(userID uint) ([]string, error) {
	var permissions []string

	// Query the materialized view using GORM
	err := r.db.Table("user_permissions").
		Select("permission_name").
		Where("user_id = ?", userID).
		Pluck("permission_name", &permissions).Error

	if err != nil {
		logger.Logger.Error("Failed to fetch user permissions", "method", "GetUserPermissions", "error", err, "userID", userID)
	}

	return permissions, err
}

func (r *UserPermissionRepository) HasPermission(userID uint, permission string) (bool, error) {
	var count int64

	// Query the materialized view using GORM
	err := r.db.Table("user_permissions").
		Where("user_id = ? AND permission_name = ? AND banned = false", userID, permission).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to check user permission", "method", "HasPermission", "error", err, "userID", userID, "permission", permission)
	}

	return count > 0, err
}

func (r *UserPermissionRepository) HasAnyPermission(userID uint, permissions []string) (bool, error) {
	if len(permissions) == 0 {
		return true, nil
	}

	var count int64

	err := r.db.Table("user_permissions").
		Where("user_id = ? AND permission_name IN ? AND banned = false", userID, permissions).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to check user permissions", "method", "HasAnyPermission", "error", err, "userID", userID, "permissions", permissions)
	}

	return count > 0, err
}

func (r *UserPermissionRepository) GetUserPermissionsByRole(userID uint, roleType string) ([]string, error) {
	var permissions []string

	// Query the materialized view using GORM
	err := r.db.Table("user_permissions").
		Select("permission_name").
		Where("user_id = ? AND role_type = ? AND banned = false", userID, roleType).
		Pluck("permission_name", &permissions).Error

	if err != nil {
		logger.Logger.Error("Failed to fetch user permissions by role", "method", "GetUserPermissionsByRole", "error", err, "userID", userID, "roleType", roleType)
	}

	return permissions, err
}

func (r *UserPermissionRepository) SyncUserPermissions(userID uint, permissions []UserPermission) error {
	// For materialized view, we just refresh the view
	// The view automatically contains the latest data from source tables
	return r.RefreshMaterializedView()
}

func (r *UserPermissionRepository) GetUsersWithPermission(permission string) ([]uint, error) {
	var userIDs []uint

	// Query the materialized view using GORM
	err := r.db.Table("user_permissions").
		Select("user_id").
		Where("permission_name = ?", permission).
		Pluck("user_id", &userIDs).Error

	if err != nil {
		logger.Logger.Error("Failed to fetch users with permission", "method", "GetUsersWithPermission", "error", err, "permission", permission)
	}

	return userIDs, err
}

func (r *UserPermissionRepository) GetPermissionStats() (map[string]int, error) {
	var results []struct {
		PermissionName string `json:"permission_name"`
		Count          int    `json:"count"`
	}

	// Query the materialized view using GORM
	err := r.db.Table("user_permissions").
		Select("permission_name, COUNT(*) as count").
		Group("permission_name").
		Scan(&results).Error

	if err != nil {
		logger.Logger.Error("Failed to fetch permission stats", "method", "GetPermissionStats", "error", err)
		return nil, err
	}

	stats := make(map[string]int)
	for _, result := range results {
		stats[result.PermissionName] = result.Count
	}

	return stats, nil
}

func (r *UserPermissionRepository) CleanupOrphanedPermissions() error {
	// For materialized view, we just refresh the view
	// The view automatically contains the latest data from source tables
	return r.RefreshMaterializedView()
}

// RefreshMaterializedView refreshes the user_permissions materialized view
func (r *UserPermissionRepository) RefreshMaterializedView() error {
	err := r.db.Exec("SELECT refresh_user_permissions()").Error
	if err != nil {
		logger.Logger.Error("Failed to refresh materialized view", "method", "RefreshMaterializedView", "error", err)
	}
	return err
}

// GetUserStatusAndPermission gets user status from users table and checks permission from materialized view
func (r *UserPermissionRepository) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, hasPermission bool, err error) {
	// First get user status from users table
	var userStatus struct {
		Banned   bool `gorm:"column:banned"`
		Verified bool `gorm:"column:verified"`
	}

	err = r.db.Table("users").
		Select("banned, verified").
		Where("id = ?", userID).
		Scan(&userStatus).Error

	if err != nil {
		logger.Logger.Error("Failed to get user status from users table", "method", "GetUserStatusAndPermission", "error", err, "userID", userID)
		return false, false, false, err
	}

	// Then check permission from materialized view
	var count int64
	err = r.db.Table("user_permissions").
		Where("user_id = ? AND permission_name = ?", userID, permission).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to check permission from materialized view", "method", "GetUserStatusAndPermission", "error", err, "userID", userID, "permission", permission)
		return false, false, false, err
	}

	hasPermission = count > 0
	return userStatus.Banned, userStatus.Verified, hasPermission, nil
}

// GetUserStatusAndAnyPermission gets user status from users table and checks permissions from materialized view
func (r *UserPermissionRepository) GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, hasPermission bool, err error) {
	if len(permissions) == 0 {
		return false, false, true, nil
	}

	// First get user status from users table
	var userStatus struct {
		Banned   bool `gorm:"column:banned"`
		Verified bool `gorm:"column:verified"`
	}

	err = r.db.Table("users").
		Select("banned, verified").
		Where("id = ?", userID).
		Scan(&userStatus).Error

	if err != nil {
		logger.Logger.Error("Failed to get user status from users table", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID)
		return false, false, false, err
	}

	// Then check permissions from materialized view
	var count int64
	err = r.db.Table("user_permissions").
		Where("user_id = ? AND permission_name IN ?", userID, permissions).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to check permissions from materialized view", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID, "permissions", permissions)
		return false, false, false, err
	}

	hasPermission = count > 0
	return userStatus.Banned, userStatus.Verified, hasPermission, nil
}
