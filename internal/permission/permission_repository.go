package permission

import (
	"errors"

	"github.com/easy-comerce/backend/internal/permission/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	CreatePermissionsIfNotExists(permissionNames []string, showDeleted *bool) error
	ExistsByName(permissionName string) (bool, error)
	GetAllPermissions(sortBy string, sortOrder string, showDeleted *bool) ([]PermissionEntity, error)
	GetAllPermissionsGroup(sortBy string, sortOrder string, showDeleted *bool) ([]model.PermissionGroup, error)
	GetPermissionByID(id uint) (*PermissionEntity, error)
	GetPermissionByName(name string) (*PermissionEntity, error)
	GetPermissionsByIDs(ids []uint) ([]PermissionEntity, error)
	HasPermission(userID uint, permission string) (bool, error)
	HasAnyPermission(userID uint, permissions []string) (bool, error)
	GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
	GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) CreatePermissionsIfNotExists(permissionNames []string, showDeleted *bool) error {
	var existingPermissions []PermissionEntity
	query := r.db.Model(&PermissionEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.PermissionName+" IN ?", permissionNames).Find(&existingPermissions).Error; err != nil {
		l.Logger.Error("Failed to fetch existing permissions", "method", "CreatePermissionsIfNotExists", "error", err)
		return err
	}

	existingNamesMap := make(map[string]bool, len(existingPermissions))
	for _, perm := range existingPermissions {
		existingNamesMap[perm.Name] = true
	}

	var toCreatePermissions []PermissionEntity
	for _, name := range permissionNames {
		if !existingNamesMap[name] {
			value := c.PermissionMap[name]

			toCreatePermissions = append(
				toCreatePermissions,
				PermissionEntity{
					Name:        name,
					Description: &value.Desc,
					GroupName:   &value.GroupName,
				},
			)
		}
	}

	if len(toCreatePermissions) > 0 {
		if err := r.db.Create(&toCreatePermissions).Error; err != nil {
			l.Logger.Error("Failed to create permissions", "method", "CreatePermissionsIfNotExists", "error", err)
			return err
		}
	}

	return nil
}

func (r *permissionRepository) ExistsByName(permissionName string) (exists bool, err error) {
	var permission PermissionEntity
	err = r.db.Model(&PermissionEntity{}).Where(c.PermissionName+" = ?", permissionName).Select(c.FieldID).Take(&permission).Error
	if err != nil {
		l.Logger.Error("Failed to check if permission exists", "method", "ExistsByName", "error", err, "permissionName", permissionName)
		return false, err
	}
	return permission.ID != 0, nil
}

func (r *permissionRepository) GetAllPermissions(sortBy string, sortOrder string, showDeleted *bool) ([]PermissionEntity, error) {
	var permissions []PermissionEntity
	query := r.db.Model(&PermissionEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	filters := []string{c.PermissionName, c.PermissionDescription, c.PermissionGroupName}
	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &filters); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&permissions).Error
	if err != nil {
		l.Logger.Error("Failed to fetch permissions", "method", "GetAllPermissions", "error", err, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return permissions, err
}

func (r *permissionRepository) GetAllPermissionsGroup(sortBy string, sortOrder string, showDeleted *bool) ([]model.PermissionGroup, error) {
	var permissions []model.PermissionResponse
	query := r.db.Model(&PermissionEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	filters := []string{c.PermissionName, c.PermissionDescription, c.PermissionGroupName}
	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &filters); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Select(c.FieldID, c.PermissionName, c.PermissionDescription, c.PermissionGroupName).Find(&permissions).Error
	if err != nil {
		l.Logger.Error("Failed to fetch permissions", "method", "GetAllPermissionsGroupMemory", "error", err)
		return nil, err
	}

	groupOrder := make([]string, 0, len(permissions))
	groupMap := make(map[string][]model.PermissionResponse, len(permissions))

	for _, p := range permissions {
		if p.GroupName == nil || *p.GroupName == "" {
			continue
		}

		key := *p.GroupName
		if _, exists := groupMap[key]; !exists {
			groupOrder = append(groupOrder, key)
		}

		groupMap[key] = append(groupMap[key], p)
	}

	result := make([]model.PermissionGroup, 0, len(groupMap))
	for _, groupName := range groupOrder {
		result = append(result, model.PermissionGroup{
			GroupName:   groupName,
			Permissions: groupMap[groupName],
		})
	}

	return result, nil
}

func (r *permissionRepository) GetPermissionByID(id uint) (*PermissionEntity, error) {
	var permission PermissionEntity

	if err := r.db.Model(&PermissionEntity{}).Where(c.FieldID+" = ?", id).First(&permission).Error; err != nil {
		l.Logger.Error("Failed to fetch permission by ID", "method", "GetPermissionByID", "error", err, "id", id)
		return nil, err
	}

	return &permission, nil
}

func (r *permissionRepository) GetPermissionByName(name string) (*PermissionEntity, error) {
	var permission PermissionEntity

	if err := r.db.Model(&PermissionEntity{}).Where(c.PermissionName+" = ?", name).First(&permission).Error; err != nil {
		l.Logger.Error("Failed to fetch permission by name", "method", "GetPermissionByName", "error", err, "name", name)
		return nil, err
	}

	return &permission, nil
}

func (r *permissionRepository) GetPermissionsByIDs(ids []uint) ([]PermissionEntity, error) {
	var permissions []PermissionEntity

	if err := r.db.Model(&PermissionEntity{}).Where(c.FieldID+" IN ?", ids).Find(&permissions).Error; err != nil {
		l.Logger.Error("Failed to fetch permissions by IDs", "method", "GetPermissionsByIDs", "error", err, "ids", ids)
		return nil, err
	}

	return permissions, nil
}

func (r *permissionRepository) HasPermission(userID uint, permission string) (bool, error) {
	var count int64

	err := buildPermissionJoinQuery(r.db).
		Where("u."+c.FieldID+" = ? AND "+"p."+c.PermissionName+" = ?", userID, permission).
		Count(&count).Error

	if err != nil {
		l.Logger.Error("Failed to check user permission", "method", "HasPermission", "error", err, "userID", userID, "permission", permission)
	}

	return count > 0, err
}

func (r *permissionRepository) HasAnyPermission(userID uint, permissions []string) (bool, error) {
	if len(permissions) == 0 {
		return true, nil
	}

	var count int64
	err := buildPermissionJoinQuery(r.db).
		Where("u."+c.FieldID+" = ? AND "+"p."+c.PermissionName+" IN ?", userID, permissions).
		Count(&count).Error

	if err != nil {
		l.Logger.Error("Failed to check user permissions", "method", "HasAnyPermission", "error", err, "userID", userID, "permissions", permissions)
	}

	return count > 0, err
}

func (r *permissionRepository) GetUserStatusAndPermission(userID uint, permission string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	var userStatus model.PermissionUserStatus

	err = r.db.Select(c.UserBanned, c.UserVerified, c.UserRefreshToken).
		Where("id = ?", userID).
		First(&userStatus).Error

	if err != nil {
		l.Logger.Error("Failed to get user status", "method", "GetUserStatusAndPermission", "error", err, "userID", userID)
		return false, false, nil, false, errors.New("failed to get user status")
	}

	if !userStatus.Verified {
		return false, false, nil, false, errors.New("account is not verified yet")
	}

	if userStatus.Banned {
		return false, false, nil, false, errors.New("account is banned")
	}

	if userStatus.RefreshToken == nil {
		return false, false, nil, false, errors.New("invalid access/refresh token (logged out)")
	}

	hasPermission, err = r.HasPermission(userID, permission)
	if err != nil {
		l.Logger.Error("Failed to check user permission", "method", "GetUserStatusAndPermission", "error", err, "userID", userID, "permission", permission)
		return false, false, nil, false, err
	}

	return userStatus.Banned, userStatus.Verified, userStatus.RefreshToken, hasPermission, nil
}

func (r *permissionRepository) GetUserStatusAndAnyPermission(userID uint, permissions []string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	var userStatus model.PermissionUserStatus

	err = r.db.Select(c.UserBanned, c.UserVerified, c.UserRefreshToken).
		Where("id = ?", userID).
		First(&userStatus).Error

	if err != nil {
		l.Logger.Error("Failed to get user status", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID)
		return false, false, nil, false, err
	}

	hasPermission, err = r.HasAnyPermission(userID, permissions)
	if err != nil {
		l.Logger.Error("Failed to check user permissions", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID, "permissions", permissions)
		return false, false, nil, false, err
	}

	return userStatus.Banned, userStatus.Verified, userStatus.RefreshToken, hasPermission, nil
}

func buildPermissionJoinQuery(db *gorm.DB) *gorm.DB {
	return db.Table(c.TableUser + " u").
		Joins("JOIN " + c.TableUserRole + " ur ON u.id = ur.user_id").
		Joins("JOIN " + c.TableRole + " r ON ur.role_id = r.id").
		Joins("JOIN " + c.TableRolePermission + " rp ON r.id = rp.role_id").
		Joins("JOIN " + c.TablePermission + " p ON rp.permission_id = p.id").
		Where("u.deleted_at IS NULL AND r.deleted_at IS NULL AND p.deleted_at IS NULL")
}
