package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/easy-comerce/backend/internal/user/model"
	e "github.com/easy-comerce/backend/pkg/app_error"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsersPaginated(include []string, showDeleted *bool, page int, pageSize int, sortBy, sortOrder string) ([]UserEntity, int, error)
	GetUserByID(id uint, include []string, showDeleted *bool) (*UserEntity, error)
	GetAuthUserByID(id uint, includeRoles bool, includePermissions bool) (*UserEntity, error)
	GetAuthUserStatusByID(id uint) (banned bool, verified bool, refreshToken *string, err error)
	GetUserIdByEmail(email string) (*uint, error)
	GetFullUserByEmail(email string) (*UserEntity, error)
	CreateUser(user *UserEntity) (*UserEntity, error)
	ResetPassword(userID uint, password string) error
	UpdateUserInfo(userID uint, user *UserEntity) error
	UpdateNameAndPathKey(userID uint, name string, pathKey *string) error
	UpdateUserPassword(userID uint, password string, req *model.ChangePasswordRequest) error
	DeleteUser(ctx context.Context, id uint) error
	UndoDeletedUser(ctx context.Context, id uint) error
	UserExists(id uint, showDeleted *bool) (bool, error)
	GetUserEmail(id uint, showDeleted *bool) (*string, error)
	GetUserIdAndVerifiedByEmail(email string, showDeleted *bool) (*uint, bool, error)
	IsUserExists(email string) (bool, error)
	getUserPasswordByID(id uint) (string, error)
	UserExistsByEmailAndRoleId(email string, roleID uint, showDeleted *bool) (bool, error)
	SetPasswordResetToken(userID uint, token string, expiresAt time.Time) error
	GetUserByResetPasswordToken(token string) (*UserEntity, error)
	SetRefreshTokenAndLastLoginAt(userID uint, token string, expiresAt time.Time) (lastLoginAt time.Time, err error)
	SetRefreshToken(userID uint, token *string, expiresAt *time.Time) error
	SetVerifiedAndVerificationToken(userID uint, verified bool, token *string, expiresAt *time.Time) error
	GetUserByVerificationToken(token string) (*uint, error)
	GetUserByRefreshToken(token string) (*UserEntity, error)
	UpdateUserPurchaseCountAndTotalSpent(userID uint, purchaseCount uint, totalSpent uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetAllUsersPaginated(include []string, showDeleted *bool, page int, pageSize int, sortBy, sortOrder string) ([]UserEntity, int, error) {
	var users []UserEntity
	var total int64

	selectFields := getSelectableFields(include)
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

	if err := query.Model(&UserEntity{}).Count(&total).Error; err != nil {
		l.Logger.Error("❌ Failed to count users", "method", "GetAllUsersPaginated", "error", err, "include", include, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&users).Error
	if err != nil {
		l.Logger.Error("❌ Failed to fetch users paginated", "method", "GetAllUsersPaginated", "error", err, "include", include, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return users, int(total), err
}

func (r *userRepository) GetUserByID(id uint, include []string, showDeleted *bool) (*UserEntity, error) {
	var user UserEntity

	selectFields := getSelectableFields(include)
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
		l.Logger.Error("❌ Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id, "include", include)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAuthUserByID(id uint, includeRoles bool, includePermissions bool) (*UserEntity, error) {
	var user UserEntity
	query := r.db.Model(&UserEntity{}).
		Select(getDefaultFields()).
		Where(c.FieldID+" = ?", id)

	if includeRoles {
		query = query.Preload(c.UserRolesCapitalized)
	}

	if includePermissions {
		query = query.Preload(c.UserRolesPermissionsCapitalized)
	}

	if err := query.First(&user).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id, "includeRoles", includeRoles, "includePermissions", includePermissions)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAuthUserStatusByID(id uint) (bool, bool, *string, error) {
	var user UserEntity

	if err := r.db.Model(&UserEntity{}).
		Select(c.UserBanned, c.UserVerified, c.UserRefreshToken).
		Where(c.FieldID+" = ?", id).
		First(&user).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch user by ID", "method", "GetUserByID", "error", err, "id", id)
		return false, false, nil, err
	}

	return user.Banned, user.Verified, user.RefreshToken, nil
}

func (r *userRepository) GetUserIdByEmail(email string) (*uint, error) {
	var user UserEntity

	if err := r.db.Select(c.FieldID).Where(c.UserEmail+" = ?", email).First(&user).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch userID by email", "method", "GetUserIdByEmail", "error", err, "email", email)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.WrapServerError("Failed to fetch userID by email", err)
		}

		return nil, err
	}

	return &user.ID, nil
}

func (r *userRepository) GetFullUserByEmail(email string) (*UserEntity, error) {
	var user UserEntity

	if err := r.db.Model(&UserEntity{}).Where(c.UserEmail+" = ?", email).
		Preload(c.UserRolesCapitalized).
		Preload(c.UserRolesPermissionsCapitalized).
		First(&user).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch user by email", "method", "GetUserByEmail", "error", err, "email", email)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.WrapServerError("Failed to fetch user by email", err)
		}

		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *UserEntity) (*UserEntity, error) {
	err := r.db.Create(user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to create user", "method", "CreateUser", "error", err, "user", user)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.WrapServerError("Failed to create user", err)
		}

		return nil, err
	}
	return user, nil
}

func (r *userRepository) ResetPassword(userID uint, password string) error {
	user := &UserEntity{
		Password:             password,
		PasswordResetToken:   nil,
		PasswordResetExpires: nil,
		Verified:             true,
	}

	err := r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", userID).
		Select(c.UserPassword, c.UserPasswordResetToken, c.UserPasswordResetExpires, c.UserVerified).
		Updates(user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to reset user password", "method", "ResetPassword", "error", err, "userID", userID)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return e.WrapServerError("Failed to reset user password", err)
		}
	}
	return err
}

func (r *userRepository) UpdateUserInfo(userID uint, user *UserEntity) error {
	updatedUser := &UserEntity{
		Name:     user.Name,
		Banned:   user.Banned,
		Verified: user.Verified,
	}

	err := r.db.Model(&UserEntity{}).
		Select(c.UserName, c.UserBanned, c.UserVerified).
		Where(c.FieldID+" = ?", userID).Updates(updatedUser).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update user password and clear reset password token", "method", "UpdatePasswordAndClearResetPasswordToken", "error", err, "userID", userID)
	}
	return err
}

func (r *userRepository) UpdateNameAndPathKey(userID uint, name string, pathKey *string) error {
	user := &UserEntity{
		Name:    name,
		PathKey: pathKey,
	}
	user.UpdatedBy = &userID
	user.CreatedBy = &userID

	err := r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update user name and path key", "method", "UpdateNameAndPathKey", "error", err, "userID", userID)
	}
	return err
}

func (r *userRepository) UpdateUserPassword(userID uint, password string, req *model.ChangePasswordRequest) error {
	user := &UserEntity{
		Password:             password,
		PasswordResetToken:   nil,
		PasswordResetExpires: nil,
	}

	if !user.CheckPassword(req.CurrentPassword) {
		l.Logger.Error("❌ Invalid current password", "method", "ChangePassword", "userID", user.ID)
		return errors.New("invalid current password")
	}

	user.UpdatedBy = &userID
	user.Password = req.NewPassword

	if err := user.HashPassword(); err != nil {
		l.Logger.Error("❌ Failed to hash new password", "method", "ChangePassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to process new password: %w", err)
	}

	err := r.db.Model(&UserEntity{}).
		Select(c.UserPasswordResetToken, c.UserPasswordResetExpires).
		Where(c.FieldID+" = ?", userID).Updates(&user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update user password", "method", "UpdateUserPassword", "error", err, "userID", userID)
	}
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	user := &UserEntity{}
	user.DeletedAt = timeutil.GormNowUTC()
	userID, _ := cu.GetUserIDFromContext(ctx)
	user.CreatedBy = userID

	err := r.db.Omit(c.FieldUpdatedAt).Where(c.FieldID+" = ?", id).Updates(user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to soft delete user", "method", "SoftDeleteUser", "error", err, "id", id)
	}
	return err
}

func (r *userRepository) UndoDeletedUser(ctx context.Context, id uint) error {
	user := &UserEntity{}
	user.DeletedAt = nil
	user.DeletedBy = nil
	userID, _ := cu.GetUserIDFromContext(ctx)
	user.CreatedBy = userID

	err := r.db.Unscoped().Model(&UserEntity{}).
		Select(c.FieldDeletedAt, c.FieldDeletedBy).
		Where(c.FieldID+" = ?", id).Updates(user).Error

	if err != nil {
		l.Logger.Error("❌ Failed to undo deleted user", "method", "UndoDeletedUser", "error", err, "id", id)
	}
	return err
}

func (r *userRepository) UserExists(id uint, showDeleted *bool) (bool, error) {
	var user UserEntity
	query := r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID).Take(&user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if user exists", "method", "UserExists", "error", err, "id", id)
		return false, err
	}

	return user.ID != 0, nil
}

func (r *userRepository) GetUserEmail(id uint, showDeleted *bool) (*string, error) {
	var user UserEntity
	query := r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.UserEmail).Take(&user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to get user email", "method", "GetUserEmail", "error", err, "id", id)
		return nil, err
	}

	return &user.Email, nil
}

func (r *userRepository) GetUserIdAndVerifiedByEmail(email string, showDeleted *bool) (*uint, bool, error) {
	var user UserEntity
	query := r.db.Model(&UserEntity{}).Where(c.UserEmail+" = ?", email)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Select(c.FieldID, c.UserEmail, c.UserVerified).Take(&user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to get user email", "method", "GetUserEmail", "error", err, "email", email)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, e.WrapServerError("Failed to get user email", err)
		}

		return nil, false, err
	}

	return &user.ID, user.Verified, nil
}

func (r *userRepository) IsUserExists(email string) (bool, error) {
	var user UserEntity

	err := r.db.Model(&UserEntity{}).Where(c.UserEmail+" = ?", email).Select(c.FieldID).Take(&user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if user exists by email", "method", "IsUserExists", "error", err, "email", email)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, e.WrapServerError("Failed to check if user exists by email", err)
		}

		return false, err
	}

	return user.ID != 0, nil
}

func (r *userRepository) getUserPasswordByID(id uint) (string, error) {
	var user UserEntity

	err := r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", id).Select(c.UserPassword).Take(&user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to get user password by ID", "method", "getUserPasswordByID", "error", err, "id", id)
		return "", err
	}

	return user.Password, nil
}

func (r *userRepository) UserExistsByEmailAndRoleId(email string, roleID uint, showDeleted *bool) (bool, error) {
	var user UserEntity
	query := r.db.Model(&UserEntity{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Joins("JOIN "+c.TableUserRole+" ON user_roles.user_id = users.id").
		Where("users.email = ? AND user_roles.role_id = ?", email, roleID).
		Select(c.FieldID).
		Take(&user).Error

	if err != nil {
		l.Logger.Error("❌ Failed to check if user exists by email and role id", "method", "UserExistsByEmailAndRoleId", "error", err, "email", email, "roleID", roleID)
		return false, err
	}

	return user.ID != 0, nil
}

func (r *userRepository) SetPasswordResetToken(userID uint, token string, expiresAt time.Time) error {
	user := &UserEntity{
		PasswordResetToken:   &token,
		PasswordResetExpires: &expiresAt,
	}

	err := r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to set password reset token", "method", "SetPasswordResetToken", "error", err, "userID", userID)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return e.WrapServerError("failed to set password reset token", err)
		}
	}
	return err
}

func (r *userRepository) GetUserByResetPasswordToken(token string) (*UserEntity, error) {
	var user UserEntity
	err := r.db.Where(c.UserPasswordResetToken+" = ? AND "+c.UserPasswordResetExpires+" > ?", token, timeutil.NowUTC()).
		Select(c.FieldID, c.UserEmail, c.UserPasswordResetExpires).
		First(&user).Error
	if err != nil {
		l.Logger.Error("❌ Failed to check if reset password token is valid and not expired", "method", "IsValidResetPasswordToken", "error", err, "token", token)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.WrapServerError("failed to check if reset password token is valid and not expired", err)
		}

		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SetRefreshTokenAndLastLoginAt(userID uint, token string, expiresAt time.Time) (lastLoginAt time.Time, err error) {
	lastLoginAt = timeutil.NowUTC()
	user := &UserEntity{
		RefreshToken:        &token,
		RefreshTokenExpires: &expiresAt,
		LastLoginAt:         &lastLoginAt,
	}

	err = r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		l.Logger.Error("Failed to set refresh token", "method", "SetRefreshToken", "error", err, "userID", userID)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return lastLoginAt, e.WrapServerError("failed to set refresh token", err)
		}
	}
	return lastLoginAt, err
}

func (r *userRepository) SetRefreshToken(userID uint, token *string, expiresAt *time.Time) error {
	user := &UserEntity{
		RefreshToken:        token,
		RefreshTokenExpires: expiresAt,
	}

	err := r.db.Model(&UserEntity{}).Select(c.UserRefreshToken, c.UserRefreshTokenExpires).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		l.Logger.Error("Failed to set refresh token", "method", "SetRefreshToken", "error", err, "userID", userID)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return e.WrapServerError("failed to set refresh token", err)
		}
	}

	return err
}

func (r *userRepository) SetVerifiedAndVerificationToken(userID uint, verified bool, token *string, expiresAt *time.Time) error {
	user := &UserEntity{
		VerificationToken:   token,
		VerificationExpires: expiresAt,
		Verified:            verified,
	}

	err := r.db.Model(&UserEntity{}).Select(c.UserVerificationToken, c.UserVerificationExpires, c.UserVerified).Where(c.FieldID+" = ?", userID).Updates(user).Error
	if err != nil {
		l.Logger.Error("Failed to set verified and verification token", "method", "SetVerifiedAndVerificationToken", "error", err, "userID", userID, "verified", verified)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return e.WrapServerError("failed to set verified and verification token", err)
		}
	}

	return errors.New("failed to set verified and verification token")
}

func (r *userRepository) GetUserByVerificationToken(token string) (*uint, error) {
	if token == "" {
		return nil, errors.New("verification token is required")
	}

	var user UserEntity
	expiry := timeutil.AddHoursUTC(c.VerificationTokenExpiryHours)

	if err := r.db.Model(&UserEntity{}).
		Select(c.FieldID).
		Where(c.UserVerificationToken+" = ? AND "+c.UserVerificationExpires+" BETWEEN ? AND ?", token, timeutil.NowUTC(), expiry).
		First(&user).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch user by verification token", "method", "GetUserByVerificationToken", "error", err, "token", token)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.WrapServerError("failed to fetch user by verification token", err)
		}

		return nil, errors.New("invalid or expired verification token")
	}

	return &user.ID, nil
}

func (r *userRepository) GetUserByRefreshToken(token string) (*UserEntity, error) {
	if token == "" {
		return nil, errors.New("refresh token is required")
	}

	var user UserEntity
	expiry := timeutil.AddHoursUTC(c.RefreshTokenExpiryHours)

	if err := r.db.Model(&UserEntity{}).Preload(c.UserRolesCapitalized).
		Where(c.UserRefreshToken+" = ? AND "+c.UserRefreshTokenExpires+" BETWEEN ? AND ?", token, timeutil.NowUTC(), expiry).
		First(&user).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch user by refresh token", "method", "GetUserByRefreshToken", "error", err, "token", token)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewServerError("Failed to fetch user by refresh token")
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUserPurchaseCountAndTotalSpent(userID uint, purchaseCount uint, totalSpent uint) error {
	err := r.db.Model(&UserEntity{}).Where(c.FieldID+" = ?", userID).Updates(map[string]any{
		c.UserPurchaseCount: gorm.Expr(c.UserPurchaseCount+" + ?", purchaseCount),
		c.UserTotalSpent:    gorm.Expr(c.UserTotalSpent+" + ?", totalSpent),
	}).Error

	if err != nil {
		l.Logger.Error("❌ Failed to update user purchase count and total spent", "method", "UpdateUserPurchaseCountAndTotalSpent", "error", err, "userID", userID, "purchaseCount", purchaseCount, "totalSpent", totalSpent)
	}
	return err
}

func getSelectableFields(include []string) []string {
	defaultFields := []string{c.FieldID, c.UserEmail, c.UserName, c.UserVerified, c.UserBanned, c.UserLastLoginAt, c.UserPurchaseCount, c.UserTotalSpent, c.FieldCreatedAt, c.FieldUpdatedAt}
	optionalFields := []string{c.UserPassword}
	// optionalFields := []string{c.UserPassword, c.UserRolesCapitalized, c.UserRolesPermissionsCapitalized}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}

func getDefaultFields() []string {
	return []string{
		c.FieldID,
		c.UserEmail,
		c.UserName,
		c.UserVerified,
		c.UserBanned,
		c.UserPurchaseCount,
		c.UserTotalSpent,
		c.UserPathKey,
	}
}

// type BaseModel struct {
// 	UpdatedAt *time.Time      `json:"updated_at,omitempty" gorm:"column:updated_at; default:null"`
// 	UpdatedBy *uint           `json:"updated_by,omitempty" gorm:"column:updated_by; default:null"`
// 	DeletedBy *uint           `json:"deleted_by,omitempty" gorm:"column:deleted_by; default:null"`
// 	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index; column:deleted_at; default:null"`
// }

// type User struct {
// 	PathKey       *string     `gorm:"column:path_key"`
// 	Roles         []role.Role `gorm:"many2many:user_roles;"`
// 	Email         string      `gorm:"not null; column:email"`
// 	Name          string      `gorm:"not null; column:name"`
// 	Verified      bool        `gorm:"default:false; column:verified"`
// 	Banned        bool        `gorm:"default:false; column:banned"`
// 	PurchaseCount uint        `gorm:"default:0; column:purchase_count"`
// 	TotalSpent    uint        `gorm:"default:0; column:total_spent"`
// }
