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

// **REQUIRED
func (r *UserRepository) GetAllUsersPaginated(include []string, page int, pageSize int, sortBy, sortOrder string) ([]User, int, error) {
	var users []User
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
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

// **REQUIRED
func (r *UserRepository) GetUserByID(id uint, include []string) (*User, error) {
	var user User

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if utils.ContainsString(include, constants.UserRoles) {
		query = query.Preload(constants.UserRolesCapitalized)
	}
	if utils.ContainsString(include, constants.UserPermissions) {
		query = query.Preload(constants.UserRolesPermissionsCapitalized)
	}

	if err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", id).First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &user, nil
}

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
func (r *UserRepository) GetUserIdByEmail(email string) (*uint, error) {
	var user User

	if err := r.db.Select(constants.FieldID).Where(constants.UserEmail+" = ?", email).First(&user).Error; err != nil {
		logger.Logger.Error("Failed to fetch userID by email", "method", "GetUserIdByEmail", "error", err, "email", email)
		return nil, err
	}

	return &user.ID, nil
}

// **REQUIRED
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

// **REQUIRED
func (r *UserRepository) CreateUser(user *User) (*User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		logger.Logger.Error("Failed to create user", "method", "CreateUser", "error", err, "user", user)
		return nil, err
	}
	return user, nil
}

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
func (r *UserRepository) UpdateUserPassword(userID uint, password string) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Update(constants.UserPassword, password).Error
	if err != nil {
		logger.Logger.Error("Failed to update user password", "method", "UpdateUserPassword", "error", err, "userID", userID)
	}
	return err
}

// **REQUIRED
func (r *UserRepository) DeleteUser(id uint) error {
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", id).Update(constants.FieldDeletedAt, timeutil.NowUTC()).Error
	if err != nil {
		logger.Logger.Error("Failed to soft delete user", "method", "SoftDeleteUser", "error", err, "id", id)
	}
	return err
}

// **REQUIRED
func (r *UserRepository) UserExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", id).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists", "method", "UserExists", "error", err, "id", id)
	}
	return count > 0, err
}

// **REQUIRED
func (r *UserRepository) IsUserExists(email string) (bool, error) {
	var count int64

	err := r.db.Model(&User{}).Where(constants.UserEmail+" = ?", email).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email", "method", "IsUserExists", "error", err, "email", email)
	}
	return count > 0, err
}

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
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

// func (r *UserRepository) UpdateLastLogin(userID uint) (time.Time, error) {
// 	now := timeutil.NowUTC()
// 	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Update(UserLastLoginAt, now).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to update last login", "method", "UpdateLastLogin", "error", err, "userID", userID)
// 	}
// 	return now, err
// }

// func (r *UserRepository) GetUserByPasswordResetToken(token string, include []string) (*User, error) {
// 	var user User

// 	selectFields := r.getLoginSelectableFields(include)
// 	query := r.db.Select(strings.Join(selectFields, ", "))

// 	currentTime := timeutil.NowUTC()
// 	if err := query.Preload("Roles", func(db *gorm.DB) *gorm.DB {
// 		return db.Select("id, role_name, role_type")
// 	}).Preload("Roles.Permissions", func(db *gorm.DB) *gorm.DB {
// 		return db.Select("id, name")
// 	}).Where("password_reset_token = ? AND password_reset_expires > ?", token, currentTime).First(&user).Error; err != nil {
// 		logger.Logger.Error("Failed to fetch user by password reset token", "method", "GetUserByPasswordResetToken", "error", err, "token", token, "include", include, "currentTime", currentTime)
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (r *UserRepository) ClearPasswordResetToken(userID uint) error {
// 	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]interface{}{
// 		"password_reset_token":   nil,
// 		"password_reset_expires": nil,
// 	}).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to clear password reset token", "method", "ClearPasswordResetToken", "error", err, "userID", userID)
// 	}
// 	return err
// }

// func (r *UserRepository) AssignRolesToUser(userID uint, roleIDs []uint) error {
// 	err := r.db.Transaction(func(tx *gorm.DB) error {
// 		var roles []interface{}
// 		if err := tx.Select("id, role_name, role_type, description, created_at, updated_at").Where(constants.FieldID+" IN ?", roleIDs).Find(&roles).Error; err != nil {
// 			return err
// 		}

// 		var user User
// 		if err := tx.Select("id").First(&user, userID).Error; err != nil {
// 			return err
// 		}

// 		return tx.Model(&user).Association("Roles").Replace(roles)
// 	})
// 	if err != nil {
// 		logger.Logger.Error("Failed to assign roles to user", "method", "AssignRolesToUser", "error", err, "userID", userID, "roleIDs", roleIDs)
// 	}
// 	return err
// }

// func (r *UserRepository) ClearRefreshToken(userID uint) error {
// 	err := r.db.Model(&User{}).Where(constants.FieldID+" = ?", userID).Updates(map[string]any{
// 		"refresh_token":         nil,
// 		"refresh_token_expires": nil,
// 	}).Error
// 	if err != nil {
// 		logger.Logger.Error("Failed to clear refresh token", "method", "ClearRefreshToken", "error", err, "userID", userID)
// 	}
// 	return err
// }

// func (r *UserRepository) GetPasswordByUserID(userID uint) (*User, error) {
// 	var user User
// 	query := r.db.Select(UserPassword).
// 		Where(constants.FieldID+" = ?", userID)

// 	if err := query.First(&user).Error; err != nil {
// 		logger.Logger.Error("Failed to fetch UserPassword by id", "method", "GetPasswordByUserID", "error", err, "userID", userID)
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (r *UserRepository) UpdateUser(user *User) (*User, error) {
// 	err := r.db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Save(user).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Model(user).Association("Roles").Replace(user.Roles); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		logger.Logger.Error("Failed to update user", "method", "UpdateUser", "error", err, "user", user)
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (r *UserRepository) getLoginSelectableFields(include []string) []string {
// 	defaultFields := []string{constants.FieldID, UserEmail, UserName, UserVerified, UserBanned}
// 	optionalFields := []string{UserPassword}
// 	return utils.BuildSelectFields(defaultFields, optionalFields, include)
// }

// // ToSharedInterface converts UserRepository to shared.UserRepositoryInterface
// func (r *UserRepository) ToSharedInterface() shared.UserRepositoryInterface {
// 	return &sharedUserRepository{repo: r}
// }

// sharedUserRepository wraps UserRepository to implement shared.UserRepositoryInterface
// type sharedUserRepository struct {
// 	repo *UserRepository
// }

// func (s *sharedUserRepository) GetUserByID(id uint, include []string) (*shared.User, error) {
// 	user, err := s.repo.GetUserByID(id, include)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &shared.User{
// 		ID:    user.ID,
// 		Email: user.Email,
// 		Name:  user.Name,
// 	}, nil
// }

// func (s *sharedUserRepository) AssignRolesToUser(userID uint, roleIDs []uint) error {
// 	return s.repo.AssignRolesToUser(userID, roleIDs)
// }

// // IsUserBanned checks if a user is banned (lightweight query - only checks banned status)
// func (r *UserRepository) IsUserBanned(userID uint) (bool, error) {
// 	var banned bool
// 	err := r.db.Model(&User{}).
// 		Select("banned").
// 		Where("id = ? AND deleted_at IS NULL", userID).
// 		Scan(&banned).Error

// 	if err != nil {
// 		logger.Logger.Error("Failed to check if user is banned", "method", "IsUserBanned", "error", err, "userID", userID)
// 		return false, err
// 	}

// 	return banned, nil
// }

// GetUserStatus gets user status (banned, verified) from users table
// func (r *UserRepository) GetUserStatus(userID uint) (banned bool, verified bool, err error) {
// 	var result struct {
// 		Banned   bool `gorm:"column:banned"`
// 		Verified bool `gorm:"column:verified"`
// 	}

// 	err = r.db.Model(&User{}).
// 		Select("banned, verified").
// 		Where("id = ? AND deleted_at IS NULL", userID).
// 		Scan(&result).Error

// 	if err != nil {
// 		logger.Logger.Error("Failed to get user status from users table", "method", "GetUserStatus", "error", err, "userID", userID)
// 		return false, false, err
// 	}

// 	return result.Banned, result.Verified, nil
// }
