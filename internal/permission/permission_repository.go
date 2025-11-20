package permission

import (
	"errors"

	"github.com/easy-comerce/backend/internal/permission/model"
	"github.com/easy-comerce/backend/pkg/base"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	base.BaseRepository[PermissionEntity]
	HasPermission(userID uint, permission string, sqlComment ...string) (bool, error)
	HasAnyPermission(userID uint, permissions []string, sqlComment ...string) (bool, error)
	GetUserStatusAndPermission(userID uint, permission string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
	GetUserStatusAndAnyPermission(userID uint, permissions []string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error)
}

type permissionRepository struct {
	base.BaseRepository[PermissionEntity]
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		base.NewBaseRepository[PermissionEntity](db),
		db,
	}
}

func (r *permissionRepository) HasPermission(userID uint, permission string, sqlComment ...string) (bool, error) {
	if permission == "" {
		return false, nil
	}

	var exists bool
	subQuery := buildPermissionJoinQuery(r.db).
		Select("1").
		Where("u."+c.FieldID+" = ? AND p."+c.PermissionName+" = ?", userID, permission).
		Limit(1)

	err := r.db.Raw("SELECT EXISTS(?)", subQuery).Scan(&exists).Error
	if err != nil {
		l.Error("Failed to check user permissions", "method", "HasPermission", "error", err, "userID", userID, "permission", permission)
	}

	return exists, err
}

func (r *permissionRepository) HasAnyPermission(userID uint, permissions []string, sqlComment ...string) (bool, error) {
	if len(permissions) == 0 {
		return false, nil
	}

	var exists bool
	subQuery := buildPermissionJoinQuery(r.db).
		Select("1").
		Where("u."+c.FieldID+" = ? AND p."+c.PermissionName+" IN ?", userID, permissions).
		Limit(1)

	err := r.db.Raw("SELECT EXISTS(?)", subQuery).Scan(&exists).Error
	if err != nil {
		l.Error("Failed to check user permissions", "method", "HasAnyPermission", "error", err, "userID", userID, "permissions", permissions)
	}

	return exists, err
}

func (r *permissionRepository) GetUserStatusAndPermission(userID uint, permission string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	var userStatus model.PermissionUserStatus
	query := r.db.Select(c.UserBanned, c.UserVerified, c.UserRefreshToken).Where("id = ?", userID)

	err = query.First(&userStatus).Error
	if err != nil {
		l.Error("Failed to get user status", "method", "GetUserStatusAndPermission", "error", err, "userID", userID)
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
		l.Error("Failed to check user permission", "method", "GetUserStatusAndPermission", "error", err, "userID", userID, "permission", permission)
		return false, false, nil, false, err
	}

	return userStatus.Banned, userStatus.Verified, userStatus.RefreshToken, hasPermission, nil
}

func (r *permissionRepository) GetUserStatusAndAnyPermission(userID uint, permissions []string, sqlComment ...string) (banned bool, verified bool, refreshToken *string, hasPermission bool, err error) {
	var userStatus model.PermissionUserStatus
	query := r.db.Select(c.UserBanned, c.UserVerified, c.UserRefreshToken).Where("id = ?", userID)
	// query = appendSqlComment(query, sqlComment...)

	err = query.First(&userStatus).Error
	if err != nil {
		l.Error("Failed to get user status", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID, "sqlComment", sqlComment)
		return false, false, nil, false, err
	}

	hasPermission, err = r.HasAnyPermission(userID, permissions)
	if err != nil {
		l.Error("Failed to check user permissions", "method", "GetUserStatusAndAnyPermission", "error", err, "userID", userID, "permissions", permissions, "sqlComment", sqlComment)
		return false, false, nil, false, err
	}

	return userStatus.Banned, userStatus.Verified, userStatus.RefreshToken, hasPermission, nil
}

func buildPermissionJoinQuery(db *gorm.DB) *gorm.DB {
	return db.Table(c.TableUser + " u").
		Joins("INNER JOIN " + c.TableUserRole + " ur ON u.id = ur.user_id").
		Joins("INNER JOIN " + c.TableRole + " r ON ur.role_id = r.id AND r.deleted_at IS NULL").
		Joins("INNER JOIN " + c.TableRolePermission + " rp ON r.id = rp.role_id").
		Joins("INNER JOIN " + c.TablePermission + " p ON rp.permission_id = p.id AND p.deleted_at IS NULL").
		Where("u.deleted_at IS NULL")
}

// func appendSqlComment(db *gorm.DB, sqlComment ...string) *gorm.DB {
// 	if len(sqlComment) > 0 {
// 		db = db.Clauses(hints.Comment("sql-comment", strings.Join(sqlComment, ", ")))
// 	}
// 	return db
// }

// func appendSqlComment(db *gorm.DB, sqlComment ...string) *gorm.DB {
// 	if len(sqlComment) > 0 {
// 		comment := strings.Join(sqlComment, ", ")
// 		return db.Clauses(clause.Expr{
// 			SQL: fmt.Sprintf("/* %s */", comment),
// 		})
// 	}
// 	return db
// }
