package auth

import (
	"strings"
	"time"

	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	models.BaseModel
	Email                string     `json:"email" gorm:"uniqueIndex;not null;column:email"`
	Password             string     `json:"-" gorm:"not null;column:password"`
	Name                 string     `json:"name" gorm:"not null;column:name"`
	Verified             bool       `json:"verified" gorm:"default:false;column:verified"`
	Banned               bool       `json:"banned" gorm:"default:false;column:banned"`
	PasswordResetToken   *string    `json:"-" gorm:"column:password_reset_token"`
	PasswordResetExpires *time.Time `json:"-" gorm:"column:password_reset_expires"`
	RefreshToken         *string    `json:"-" gorm:"column:refresh_token"`
	RefreshTokenExpires  *time.Time `json:"-" gorm:"column:refresh_token_expires"`
	LastLoginAt          *time.Time `json:"last_login_at,omitempty" gorm:"column:last_login_at"`
	Roles                []Role     `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

const (
	UserEmail       = "email"
	UserName        = "name"
	UserVerified    = "verified"
	UserBanned      = "banned"
	UserLastLoginAt = "last_login_at"
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
