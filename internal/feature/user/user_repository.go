package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/constants"
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

	if utils.ContainsString(include, c.UserRoles) {
		query = query.Preload(c.UserRolesCapitalized)
	}

	if utils.ContainsString(include, c.UserPermissions) {
		query = query.Preload(c.UserRolesPermissionsCapitalized)
	}

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

	if utils.ContainsString(include, c.UserRoles) {
		query = query.Preload(c.UserRolesCapitalized)
	}

	if utils.ContainsString(include, c.UserPermissions) {
		query = query.Preload(c.UserRolesPermissionsCapitalized)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if err := query.Where(c.FieldID+" = ?", id).First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAuthUserByID(id uint, includeRoles bool, includePermissions bool) (*User, error) {
	var user User
	query := r.db.Model(&User{}).Where(c.FieldID+" = ?", id)

	if includeRoles {
		query = query.Preload(c.UserRolesCapitalized)
	}

	if includePermissions {
		query = query.Preload(c.UserRolesPermissionsCapitalized)
	}

	if err := query.First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id, "includeRoles", includeRoles, "includePermissions", includePermissions)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAuthUserStatusByID(id uint) (bool, bool, *string, error) {
	var user User

	if err := r.db.Model(&User{}).
		Select(c.UserBanned, c.UserVerified, c.UserRefreshToken).
		Where(c.FieldID+" = ?", id).
		First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id)
		return false, false, nil, err
	}

	return user.Banned, user.Verified, user.RefreshToken, nil
}

func (r *UserRepository) GetUserIdByEmail(email string) (*uint, error) {
	var user User

	if err := r.db.Select(c.FieldID).Where(c.UserEmail+" = ?", email).First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch userID by email", "method", "GetUserIdByEmail", "error", err, "email", email)
		return nil, err
	}

	return &user.ID, nil
}

func (r *UserRepository) GetFullUserByEmail(email string) (*User, error) {
	var user User

	if err := r.db.Model(&User{}).Where(c.UserEmail+" = ?", email).
		Preload(c.UserRolesCapitalized).
		Preload(c.UserRolesPermissionsCapitalized).
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

func (r *UserRepository) ResetPassword(userID uint, password string) error {
	user := &User{
		Password:             password,
		PasswordResetToken:   nil,
		PasswordResetExpires: nil,
		Verified:             true,
	}

	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).
		Select(c.UserPassword, c.UserPasswordResetToken, c.UserPasswordResetExpires, c.UserVerified).
		Updates(user).Error
	if err != nil {
		logger.Logger.Error("Failed to reset user password", "method", "ResetPassword", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateUserInfo(userID uint, user *User) error {
	updatedUser := &User{
		Name:     user.Name,
		Banned:   user.Banned,
		Verified: user.Verified,
	}

	err := r.db.Model(&User{}).
		Select(c.UserName, c.UserBanned, c.UserVerified).
		Where(c.FieldID+" = ?", userID).Updates(updatedUser).Error
	if err != nil {
		logger.Logger.Error("Failed to update user password and clear reset password token", "method", "UpdatePasswordAndClearResetPasswordToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateNameAndPathKey(userID uint, name string, pathKey *string) error {
	user := &User{
		Name:    name,
		PathKey: pathKey,
	}

	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		logger.Logger.Error("Failed to update user name and path key", "method", "UpdateNameAndPathKey", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateUserPassword(userID uint, password string, req *ChangePasswordRequest) error {
	user := &User{
		Password:             password,
		PasswordResetToken:   nil,
		PasswordResetExpires: nil,
	}

	if !user.CheckPassword(req.CurrentPassword) {
		logger.Logger.Error("Invalid current password", "method", "ChangePassword", "userID", user.ID)
		return errors.New("invalid current password")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash new password", "method", "ChangePassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to process new password: %w", err)
	}

	err := r.db.Model(&User{}).
		Select(c.UserPassword, c.UserPasswordResetToken, c.UserPasswordResetExpires).
		Where(c.FieldID+" = ?", userID).Updates(&user).Error
	if err != nil {
		logger.Logger.Error("Failed to update user password", "method", "UpdateUserPassword", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) DeleteUser(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&User{}).Error
	if err != nil {
		logger.Logger.Error("Failed to soft delete user", "method", "SoftDeleteUser", "error", err, "id", id)
	}
	return err
}

func (r *UserRepository) UndoDeletedUser(id uint) error {
	err := r.db.Unscoped().Model(&User{}).Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil).Error
	if err != nil {
		logger.Logger.Error("Failed to undo deleted user", "method", "UndoDeletedUser", "error", err, "id", id)
	}
	return err
}

func (r *UserRepository) UserExists(id uint, showDeleted *bool) (bool, error) {
	var user User
	query := r.db.Model(&User{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&user).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists", "method", "UserExists", "error", err, "id", id)
		return false, err
	}

	return user.ID != 0, nil
}

func (r *UserRepository) GetUserEmail(id uint, showDeleted *bool) (*string, error) {
	var user User
	query := r.db.Model(&User{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.UserEmail).Take(&user).Error
	if err != nil {
		logger.Logger.Error("Failed to get user email", "method", "GetUserEmail", "error", err, "id", id)
		return nil, err
	}

	return &user.Email, nil
}

func (r *UserRepository) IsUserExists(email string) (bool, error) {
	var user User

	err := r.db.Model(&User{}).Where(c.UserEmail+" = ?", email).Select(c.FieldID).Take(&user).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email", "method", "IsUserExists", "error", err, "email", email)
		return false, err
	}

	return user.ID != 0, nil
}

func (r *UserRepository) UserExistsByEmailAndRoleId(email string, roleID uint, showDeleted *bool) (bool, error) {
	var user User
	query := r.db.Model(&User{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Joins("JOIN "+c.TableUserRole+" ON user_roles.user_id = users.id").
		Where("users.email = ? AND user_roles.role_id = ?", email, roleID).
		Select(c.FieldID).
		Take(&user).Error

	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email and role id", "method", "UserExistsByEmailAndRoleId", "error", err, "email", email, "roleID", roleID)
		return false, err
	}

	return user.ID != 0, nil
}

func (r *UserRepository) SetPasswordResetToken(userID uint, token string, expiresAt time.Time) error {
	user := &User{
		PasswordResetToken:   &token,
		PasswordResetExpires: &expiresAt,
	}

	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		logger.Logger.Error("Failed to set password reset token", "method", "SetPasswordResetToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) GetUserByResetPasswordToken(token string) (*User, error) {
	var user User
	err := r.db.Where(c.UserPasswordResetToken+" = ?", token).
		Where(c.UserPasswordResetExpires+" > ?", timeutil.NowUTC()).
		First(&user).Error
	if err != nil {
		logger.Logger.Error("Failed to check if reset password token is valid and not expired", "method", "IsValidResetPasswordToken", "error", err, "token", token)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SetRefreshTokenAndLastLoginAt(userID uint, token string, expiresAt time.Time) (lastLoginAt time.Time, err error) {
	lastLoginAt = timeutil.NowUTC()
	user := &User{
		RefreshToken:        &token,
		RefreshTokenExpires: &expiresAt,
		LastLoginAt:         &lastLoginAt,
	}

	err = r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "SetRefreshToken", "error", err, "userID", userID)
	}
	return lastLoginAt, err
}

func (r *UserRepository) SetRefreshToken(userID uint, token *string, expiresAt *time.Time) error {
	user := &User{
		RefreshToken:        token,
		RefreshTokenExpires: expiresAt,
	}

	err := r.db.Model(&User{}).Select(c.UserRefreshToken, c.UserRefreshTokenExpires).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "SetRefreshToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) GetUserByRefreshToken(token string) (*User, error) {
	var user User

	if err := r.db.Model(&User{}).
		Preload(c.UserRolesCapitalized).
		Where(c.UserRefreshToken+" = ? AND "+c.UserRefreshTokenExpires+" > ?", token, timeutil.NowUTC()).
		First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by refresh token", "method", "GetUserByRefreshToken", "error", err, "token", token)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserPurchaseCountAndTotalSpent(userID uint, purchaseCount uint, totalSpent uint) error {
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserPurchaseCount: gorm.Expr(c.UserPurchaseCount+" + ?", purchaseCount),
		c.UserTotalSpent:    gorm.Expr(c.UserTotalSpent+" + ?", totalSpent),
	}).Error

	if err != nil {
		logger.Logger.Error("Failed to update user purchase count and total spent", "method", "UpdateUserPurchaseCountAndTotalSpent", "error", err, "userID", userID, "purchaseCount", purchaseCount, "totalSpent", totalSpent)
	}
	return err
}

func (r *UserRepository) UpdateUserRole(userID uint, roleID uint) error {
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Update("Roles", roleID).Error
	if err != nil {
		logger.Logger.Error("Failed to update user role", "method", "UpdateUserRole", "error", err, "userID", userID, "roleID", roleID)
	}
	return err
}

func (r *UserRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{c.FieldID, c.UserEmail, c.UserName, c.UserVerified, c.UserBanned, c.UserLastLoginAt, c.UserPurchaseCount, c.UserTotalSpent, c.FieldCreatedAt, c.FieldUpdatedAt}
	optionalFields := []string{c.UserPassword, c.UserRolesCapitalized, c.UserRolesPermissionsCapitalized}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
