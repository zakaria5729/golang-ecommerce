package user

import (
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

	if err := r.db.Model(&User{}).Where(c.FieldID+" = ?", id).First(&user).Error; err != nil {
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

func (r *UserRepository) GetAuthUserStatusByID(id uint) (bool, bool, error) {
	var user User

	if err := r.db.Model(&User{}).
		Select(c.UserBanned, c.UserVerified).
		Where(c.FieldID+" = ?", id).
		First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id)
		return false, false, err
	}

	return user.Banned, user.Verified, nil
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
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserPassword:             password,
		c.UserPasswordResetToken:   nil,
		c.UserPasswordResetExpires: nil,
		c.UserVerified:             true,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to reset user password", "method", "ResetPassword", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateUserInfo(userID uint, user *User) error {
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserPassword: user.Password,
		c.UserName:     user.Name,
		c.UserBanned:   user.Banned,
		c.UserVerified: user.Verified,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to update user password and clear reset password token", "method", "UpdatePasswordAndClearResetPasswordToken", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateNameAndPathKey(userID uint, name string, pathKey *string) error {
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserName:    name,
		c.UserPathKey: pathKey,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to update user name and path key", "method", "UpdateNameAndPathKey", "error", err, "userID", userID)
	}
	return err
}

func (r *UserRepository) UpdateUserPassword(userID uint, password string) error {
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Update(c.UserPassword, password).Error
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
	err := r.db.Unscoped().Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil).Error
	if err != nil {
		logger.Logger.Error("Failed to undo deleted user", "method", "UndoDeletedUser", "error", err, "id", id)
	}
	return err
}

func (r *UserRepository) UserExists(id uint, showDeleted *bool) (bool, error) {
	var count int64
	query := r.db.Model(&User{}).Where(c.FieldID+" = ?", id)

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

	err := r.db.Model(&User{}).Where(c.UserEmail+" = ?", email).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email", "method", "IsUserExists", "error", err, "email", email)
	}
	return count > 0, err
}

func (r *UserRepository) UserExistsByEmailAndRoleId(email string, roleID uint, showDeleted *bool) (bool, error) {
	var count int64
	query := r.db.Model(&User{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Joins("JOIN "+c.TableUserRole+" ON user_roles.user_id = users.id").
		Where("users.email = ? AND user_roles.role_id = ?", email, roleID).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email and role id", "method", "UserExistsByEmailAndRoleId", "error", err, "email", email, "roleID", roleID)
	}

	return count > 0, err
}

func (r *UserRepository) SetPasswordResetToken(userID uint, token string, expiresAt time.Time) error {
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserPasswordResetToken:   token,
		c.UserPasswordResetExpires: expiresAt,
	}).Error
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
	err = r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserRefreshToken:        token,
		c.UserRefreshTokenExpires: expiresAt,
		c.UserLastLoginAt:         lastLoginAt,
	}).Error
	if err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "SetRefreshToken", "error", err, "userID", userID)
	}
	return lastLoginAt, err
}

func (r *UserRepository) SetRefreshToken(userID uint, token *string, expiresAt *time.Time) error {
	err := r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserRefreshToken:        token,
		c.UserRefreshTokenExpires: expiresAt,
	}).Error
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
	var user User
	err := r.db.Select(c.UserPurchaseCount, c.UserTotalSpent).Where(c.FieldID+" = ?", userID).First(&user).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "UpdateUserPurchaseCountAndTotalSpent", "error", err, "id", userID)
		return err
	}

	err = r.db.Model(&User{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserPurchaseCount: user.PurchaseCount + purchaseCount,
		c.UserTotalSpent:    user.TotalSpent + totalSpent,
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
	defaultFields := []string{c.FieldID, c.UserEmail, c.UserName, c.UserVerified, c.UserBanned, c.FieldCreatedAt, c.FieldUpdatedAt}
	optionalFields := []string{c.UserLastLoginAt, c.UserPassword}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
