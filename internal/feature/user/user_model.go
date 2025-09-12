package user

import (
	"strings"
	"time"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	models.BaseModel
	PasswordResetExpires *time.Time  `gorm:"column:password_reset_expires"`
	RefreshTokenExpires  *time.Time  `gorm:"column:refresh_token_expires"`
	LastLoginAt          *time.Time  `gorm:"column:last_login_at"`
	PasswordResetToken   *string     `gorm:"column:password_reset_token"`
	RefreshToken         *string     `gorm:"column:refresh_token"`
	PathKey              *string     `gorm:"column:path_key"`
	Roles                []role.Role `gorm:"many2many:user_roles;"`
	Email                string      `gorm:"uniqueIndex;not null;column:email"`
	Password             string      `json:"-" gorm:"not null;column:password"`
	Name                 string      `gorm:"not null;column:name"`
	Verified             bool        `gorm:"default:false;column:verified"`
	Banned               bool        `gorm:"default:false;column:banned"`
}

func (User) TableName() string {
	return constants.TableUser
}

const (
	UserEmail                = "email"
	UserName                 = "name"
	UserVerified             = "verified"
	UserBanned               = "banned"
	UserLastLoginAt          = "last_login_at"
	UserRoles                = "roles"
	UserPermissions          = "permissions"
	UserPassword             = "password"
	UserPathKey              = "path_key"
	UserImageURL             = "image_url"
	UserRefreshToken         = "refresh_token"
	UserRefreshTokenExpires  = "refresh_token_expires"
	UserPasswordResetToken   = "password_reset_token"
	UserPasswordResetExpires = "password_reset_expires"
)

func (u *User) Sanitize() {
	if u.Email != "" {
		u.Email = utils.Trim(strings.ToLower(u.Email))
	}
	if u.Name != "" {
		u.Name = utils.Trim(u.Name)
	}
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		BaseModel:   u.BaseModel,
		Email:       u.Email,
		Name:        u.Name,
		Verified:    u.Verified,
		Banned:      u.Banned,
		LastLoginAt: u.LastLoginAt,
		ImageURL:    utils.BuildFullImageURL(u.PathKey),
		Roles:       u.Roles,
	}
}
