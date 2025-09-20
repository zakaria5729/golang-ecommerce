package user

import (
	"strings"
	"time"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: db.GetDB(),
	}
}

func (r *UserRepository) GetAllUsersPaginated(include []string, showDeleted *bool, page int, pageSize int, sortBy, sortOrder string) ([]User, int, error) {
	var users []User
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Model(&User{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count users", "method", "GetAllUsersPaginated", "error", err, "include", include, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&users).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch users paginated", "method", "GetAllUsersPaginated", "error", err, "include", include, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return users, int(total), err
}

func (r *UserRepository) GetUserByID(id uint, include []string, showDeleted *bool) (*User, error) {
	var user User

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if utils.ContainsString(include, constants.UserRoles) {
		query = query.Preload(constants.UserRolesCapitalized)
	}
	if utils.ContainsString(include, constants.UserPermissions) {
		query = query.Preload(constants.UserRolesPermissionsCapitalized)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", id).First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAuthUserByID(id uint, includeRoles bool, includePermissions bool) (*User, error) {
	var user User
	query := r.db.Model(&User{}).Where(constants.FieldID+" = ?", id)

	if includeRoles {
		query = query.Preload(constants.UserRolesCapitalized)
	}

	if includePermissions {
		query = query.Preload(constants.UserRolesPermissionsCapitalized)
	}

	if err := query.First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id, "includeRoles", includeRoles, "includePermissions", includePermissions)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAuthUserStatusByID(id uint) (bool, bool, error) {
	var user User

	if err := r.db.Model(&User{}).
		Select(constants.UserBanned, constants.UserVerified).
		Where(constants.FieldID+" = ?", id).
		First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id)
		return false, false, err
	}

	return user.Banned, user.Verified, nil
}

func (r *UserRepository) GetUserIdByEmail(email string) (*uint, error) {
	var user User

	if err := r.db.Select(constants.FieldID).Where(constants.UserEmail+" = ?", email).First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch userID by email", "method", "GetUserIdByEmail", "error", err, "email", email)
		return nil, err
	}

	return &user.ID, nil
}

func (r *UserRepository) GetFullUserByEmail(email string) (*User, error) {
	var user User

	if err := r.db.Model(&User{}).Where(constants.UserEmail+" = ?", email).
		Preload(constants.UserRolesCapitalized).
		Preload(constants.UserRolesPermissionsCapitalized).
		First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by email", "method", "GetUserByEmail", "error", err, "email", email)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user *User) (*User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		logger.Logger.Error("Failed to create user", "method", "CreateUser", "error", err, "user", user)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdatePasswordAndClearResetPasswordToken(userID uint, password string) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]any{
		constants.UserPassword:             password,
		constants.UserPasswordResetToken:   nil,
		constants.UserPasswordResetExpires: nil,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to update user password and clear reset password token", "method", "UpdatePasswordAndClearResetPasswordToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateUserInfo(userID uint, user *User) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]any{
		constants.UserPassword: user.Password,
		constants.UserName:     user.Name,
		constants.UserBanned:   user.Banned,
		constants.UserVerified: user.Verified,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to update user password and clear reset password token", "method", "UpdatePasswordAndClearResetPasswordToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateNameAndPathKey(userID uint, name string, pathKey *string) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]any{
		constants.UserName:    name,
		constants.UserPathKey: pathKey,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to update user name and path key", "method", "UpdateNameAndPathKey", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateUserPassword(userID uint, password string) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Update(constants.UserPassword, password).Error
	if err != nil {
		logger.Logger.Error("Failed to update user password", "method", "UpdateUserPassword", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) DeleteUser(id uint) error {
	err := r.db.Where(constants.FieldID+" = ?", id).Delete(&User{}).Error
	if err != nil {
		logger.Logger.Error("Failed to soft delete user", "method", "SoftDeleteUser", "error", err, "id", id)
	}
	return err
}

func (r *UserRepository) UndoDeletedUser(id uint) error {
	err := r.db.Unscoped().Where(constants.FieldID+" = ?", id).Update(constants.FieldDeletedAt, nil).Error
	if err != nil {
		logger.Logger.Error("Failed to undo deleted user", "method", "UndoDeletedUser", "error", err, "id", id)
	}
	return err
}

func (r *UserRepository) UserExists(id uint, showDeleted *bool) (bool, error) {
	var count int64
	query := r.db.Model(&User{}).Where(constants.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists", "method", "UserExists", "error", err, "id", id)
	}
	return count > 0, err
}

func (r *UserRepository) IsUserExists(email string) (bool, error) {
	var count int64

	err := r.db.Model(&User{}).Where(constants.UserEmail+" = ?", email).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email", "method", "IsUserExists", "error", err, "email", email)
	}
	return count > 0, err
}

func (r *UserRepository) UserExistsByEmail(email string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&User{}).Where(constants.UserEmail+" = ?", email)

	if excludeID != nil {
		query = query.Where(constants.FieldID+" != ?", *excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email", "method", "UserExistsByEmail", "error", err, "email", email, "excludeID", excludeID)
	}
	return count > 0, err
}

func (r *UserRepository) SetPasswordResetToken(userID uint, token string, expiresAt time.Time) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]any{
		constants.UserPasswordResetToken:   token,
		constants.UserPasswordResetExpires: expiresAt,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to set password reset token", "method", "SetPasswordResetToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) GetUserByResetPasswordToken(token string) (*User, error) {
	var user User
	err := r.db.Where(constants.UserPasswordResetToken+" = ?", token).
		Where(constants.UserPasswordResetExpires+" > ?", timeutil.NowUTC()).
		First(&user).Error
	if err != nil {
		logger.Logger.Error("Failed to check if reset password token is valid and not expired", "method", "IsValidResetPasswordToken", "error", err, "token", token)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SetRefreshTokenAndLastLoginAt(userID uint, token string, expiresAt time.Time) (lastLoginAt time.Time, err error) {
	lastLoginAt = timeutil.NowUTC()
	err = r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]any{
		constants.UserRefreshToken:        token,
		constants.UserRefreshTokenExpires: expiresAt,
		constants.UserLastLoginAt:         lastLoginAt,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "SetRefreshToken", "error", err, "userID", userID)
	}
	return lastLoginAt, err
}

func (r *UserRepository) SetRefreshToken(userID uint, token *string, expiresAt *time.Time) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]any{
		constants.UserRefreshToken:        token,
		constants.UserRefreshTokenExpires: expiresAt,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "SetRefreshToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) GetUserByRefreshToken(token string) (*User, error) {
	var user User

	if err := r.db.Model(&User{}).
		Preload(constants.UserRolesCapitalized).
		Where(constants.UserRefreshToken+" = ? AND "+constants.UserRefreshTokenExpires+" > ?", token, timeutil.NowUTC()).
		First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by refresh token", "method", "GetUserByRefreshToken", "error", err, "token", token)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, constants.UserEmail, constants.UserName, constants.UserVerified, constants.UserBanned, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{constants.UserLastLoginAt, constants.UserPassword}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
