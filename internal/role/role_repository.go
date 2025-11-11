package role

import (
	"context"
	"errors"

	"github.com/easy-comerce/backend/internal/permission"
	e "github.com/easy-comerce/backend/pkg/app_error"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetAllRoles(include []string, showDeleted *bool, roleType *string, sortBy, sortOrder string) ([]RoleEntity, error)
	GetRoleByID(id uint, include []string, showDeleted *bool) (*RoleEntity, error)
	GetRoleWithPermissionsByType(roleType string) (*RoleEntity, error)
	GetRoleByType(roleType string, showDeleted *bool) (*RoleEntity, error)
	CreateRole(role *RoleEntity) (*RoleEntity, error)
	UpdateRole(role *RoleEntity) error
	UpdateRoleWithoutPermissions(role *RoleEntity) error
	DeleteRole(ctx context.Context, id uint) error
	UndoDeletedRole(ctx context.Context, id uint) error
	AssignRoleToUser(userID uint, roleID uint) error
	AppendPermissionsToRole(roleID uint, permissionNames []string, showDeleted *bool) error
	AppendPermissionsToRoleByIds(roleID uint, permissionIds []uint, showDeleted *bool) error
	RoleExists(id uint, showDeleted *bool) (bool, error)
	RoleExistsByName(name string, excludeID ...uint) (bool, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetAllRoles(include []string, showDeleted *bool, roleType *string, sortBy, sortOrder string) ([]RoleEntity, error) {
	var roles []RoleEntity
	query := r.db.Model(&RoleEntity{})

	if roleType != nil && *roleType != "" {
		query = query.Where(c.RoleRoleType+" = ?", *roleType)
	}

	filters := []string{c.RoleRoleName, c.RoleRoleType, c.RoleDescription}
	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, &filters); orderClause != "" {
		query = query.Order(orderClause)
	}

	if utils.ContainsString(include, c.RolePermissions) {
		query = query.Preload(c.RolePermissionsCapitalized)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Find(&roles).Error
	if err != nil {
		l.Logger.Error("Failed to fetch roles", "method", "GetAllRoles", "error", err, "include", include, "roleType", roleType, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return roles, err
}

func (r *roleRepository) GetRoleByID(id uint, include []string, showDeleted *bool) (*RoleEntity, error) {
	var role RoleEntity
	query := r.db.Model(&RoleEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if utils.ContainsString(include, c.RolePermissions) {
		query = query.Preload(c.RolePermissionsCapitalized)
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&role).Error; err != nil {
		l.Logger.Error("Failed to fetch role by ID", "method", "GetRoleByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetRoleWithPermissionsByType(roleType string) (*RoleEntity, error) {
	var role RoleEntity

	if err := r.db.Model(&RoleEntity{}).
		Preload(c.RolePermissionsCapitalized).
		Where(c.RoleRoleType+" = ?", roleType).
		First(&role).Error; err != nil {
		l.Logger.Error("Failed to fetch role by type", "method", "GetRoleWithPermissionsByType", "error", err, "roleType", roleType)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.WrapServerError("failed to fetch role by type", err)
		}

		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetRoleByType(roleType string, showDeleted *bool) (*RoleEntity, error) {
	var role RoleEntity

	query := r.db.Model(&RoleEntity{}).Where(c.RoleRoleType+" = ?", roleType)
	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Preload(c.RolePermissionsCapitalized).First(&role).Error; err != nil {
		l.Logger.Error("Failed to fetch role by type", "method", "GetRoleByType", "error", err, "roleType", roleType)
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) CreateRole(role *RoleEntity) (*RoleEntity, error) {
	err := r.db.Create(role).Error
	if err != nil {
		l.Logger.Error("Failed to create role", "method", "CreateRole", "error", err, "role", role)
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) UpdateRole(role *RoleEntity) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if err := tx.Model(role).Association(c.RolePermissionsCapitalized).Replace(role.Permissions); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		l.Logger.Error("Failed to update role", "method", "UpdateRole", "error", err, "role", role)
	}

	return err
}

func (r *roleRepository) UpdateRoleWithoutPermissions(role *RoleEntity) error {
	var updatedRole = RoleEntity{
		RoleName:    role.RoleName,
		RoleType:    role.RoleType,
		Description: role.Description,
	}
	updatedRole.UpdatedBy = role.UpdatedBy

	err := r.db.Model(&RoleEntity{}).Where(c.FieldID+" = ?", role.ID).Updates(updatedRole).Error
	if err != nil {
		l.Logger.Error("Failed to update role", "method", "UpdateRoleWithoutPermissions", "error", err, "role", role)
	}
	return err
}

func (r *roleRepository) DeleteRole(ctx context.Context, id uint) error {
	role := &RoleEntity{}
	role.DeletedAt = timeutil.GormNowUTC()
	userID, ok := ctx.Value(c.UserIDContextKey).(*uint)
	if ok {
		role.DeletedBy = userID
	}

	err := r.db.Omit(c.FieldUpdatedAt).Where(c.FieldID+" = ?", id).Updates(role).Error
	if err != nil {
		l.Logger.Error("Failed to delete role", "method", "DeleteRole", "error", err, "id", id)
	}

	return err
}

func (r *roleRepository) UndoDeletedRole(ctx context.Context, id uint) error {
	role := &RoleEntity{}
	role.DeletedAt = nil
	role.DeletedBy = nil
	userID, ok := ctx.Value(c.UserIDContextKey).(*uint)
	if ok {
		role.UpdatedBy = userID
	}

	err := r.db.Unscoped().Model(&RoleEntity{}).
		Select(c.FieldDeletedAt, c.FieldDeletedBy).
		Where(c.FieldID+" = ?", id).Updates(role).Error

	if err != nil {
		l.Logger.Error("Failed to undo deleted role", "method", "UndoDeletedRole", "error", err, "id", id)
	}

	return err
}

func (r *roleRepository) AssignRoleToUser(userID uint, roleID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var role RoleEntity
		if err := tx.Where(c.FieldID+" = ?", roleID).First(&role).Error; err != nil {
			return errors.New("role not found")
		}

		var userEmail *string
		err := tx.Table(c.TableUser).Where(c.FieldID+" = ?", userID).Select(c.UserEmail).Limit(1).Scan(&userEmail).Error

		if err != nil || userEmail == nil {
			return errors.New("user not found")
		}

		if *userEmail == config.GetConfig().AppConfig.SuperAdminEmail {
			return errors.New("You can not assign/change super admin user role")
		}

		if err := tx.Exec("DELETE FROM "+c.TableUserRole+" WHERE "+c.FieldUserID+" = ?", userID).Error; err != nil {
			return err
		}

		return tx.Exec("INSERT INTO "+c.TableUserRole+" ("+c.FieldUserID+", "+c.FieldRoleID+") VALUES (?, ?)", userID, roleID).Error
	})

	if err != nil {
		l.Logger.Error("Failed to assign roles to user", "method", "AssignRolesToUser", "error", err, "userID", userID, "roleID", roleID)
	}

	return err
}

func (r *roleRepository) AppendPermissionsToRole(roleID uint, permissionNames []string, showDeleted *bool) error {
	return appendPermissionsToRole(r.db, roleID, &permissionNames, nil, showDeleted)
}

func (r *roleRepository) AppendPermissionsToRoleByIds(roleID uint, permissionIds []uint, showDeleted *bool) error {
	return appendPermissionsToRole(r.db, roleID, nil, &permissionIds, showDeleted)
}

func (r *roleRepository) RoleExists(id uint, showDeleted *bool) (bool, error) {
	var role RoleEntity
	query := r.db.Model(&RoleEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&role).Error
	if err != nil {
		l.Logger.Error("Failed to check if role exists", "method", "RoleExists", "error", err, "id", id)
		return false, err
	}

	return role.ID != 0, nil
}

func (r *roleRepository) RoleExistsByName(name string, excludeID ...uint) (bool, error) {
	var role RoleEntity
	query := r.db.Model(&RoleEntity{}).Where(c.RoleRoleName+" = ?", name)

	if len(excludeID) > 0 {
		query = query.Where(c.FieldID+" != ?", excludeID[0])
	}

	err := query.Select(c.FieldID).Take(&role).Error
	if err != nil {
		l.Logger.Error("Failed to check if role exists by name", "method", "RoleExistsByName", "error", err, "name", name)
		return false, err
	}

	return role.ID != 0, nil
}

func appendPermissionsToRole(db *gorm.DB, roleID uint, permissionNames *[]string, permissionIds *[]uint, showDeleted *bool) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var role RoleEntity
		query := tx.Where(c.FieldID+" = ?", roleID)

		if showDeleted != nil && *showDeleted {
			query = query.Unscoped()
		}
		if err := query.Preload(c.RolePermissionsCapitalized).First(&role).Error; err != nil {
			return err
		}

		var permissions []permission.PermissionEntity
		if permissionIds != nil && len(*permissionIds) > 0 {
			query = tx.Where(c.FieldID+" IN ?", *permissionIds)
		} else {
			query = tx.Where(c.PermissionName+" IN ?", *permissionNames)
		}

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

		var newPermissions []permission.PermissionEntity
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
			return query.Association(c.RolePermissionsCapitalized).Append(newPermissions)
		}

		return nil
	})

	if err != nil {
		l.Logger.Error("Failed to add permissions to role", "method", "AppendPermissionsToRole", "error", err, "roleID", roleID, "permissionNames", permissionNames)
	}

	return err
}
