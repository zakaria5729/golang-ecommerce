package user

import (
	"errors"
	"strings"
	"time"

	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/user/model"
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type UserEntity struct {
	base.AuditEntity
	PasswordResetExpires *time.Time        `gorm:"column:password_reset_expires"`
	VerificationExpires  *time.Time        `gorm:"column:verification_expires"`
	RefreshTokenExpires  *time.Time        `gorm:"column:refresh_token_expires"`
	LastLoginAt          *time.Time        `gorm:"column:last_login_at"`
	PasswordResetToken   *string           `gorm:"column:password_reset_token"`
	VerificationToken    *string           `gorm:"column:verification_token"`
	RefreshToken         *string           `gorm:"column:refresh_token"`
	PathKey              *string           `gorm:"column:path_key"`
	Roles                []role.RoleEntity `gorm:"many2many:user_roles;joinForeignKey:user_id;joinReferences:role_id"`
	Email                string            `gorm:"not null; column:email"`
	Password             string            `json:"-" gorm:"not null; column:password"`
	Name                 string            `gorm:"not null; column:name"`
	Verified             bool              `gorm:"default:false; column:verified"`
	Banned               bool              `gorm:"default:false; column:banned"`
	PurchaseCount        uint              `gorm:"default:0; column:purchase_count"`
	TotalSpent           uint              `gorm:"default:0; column:total_spent"`
}

func (UserEntity) TableName() string {
	return constants.TableUser
}

func (u *UserEntity) Sanitize() {
	if u.Email != "" {
		u.Email = utils.Trim(strings.ToLower(u.Email))
	}
	if u.Name != "" {
		u.Name = utils.Trim(u.Name)
	}
}

func (u *UserEntity) HashPassword() error {
	err, hashedPassword := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	if hashedPassword == "" {
		return errors.New("hashed password is empty")
	}
	u.Password = hashedPassword
	return nil
}

func (u *UserEntity) CheckPassword(password string) bool {
	return utils.CheckPassword(password, u.Password)
}

func (u *UserEntity) ToResponse() *model.UserResponse {
	return &model.UserResponse{
		BaseEntity:    u.BaseEntity,
		LastLoginAt:   u.LastLoginAt,
		ImageURL:      utils.BuildFullImageURL(config.GetStorageDomain(), u.PathKey),
		Roles:         u.Roles,
		Email:         u.Email,
		Name:          u.Name,
		Verified:      u.Verified,
		Banned:        u.Banned,
		PurchaseCount: u.PurchaseCount,
		TotalSpent:    float64(u.TotalSpent / 100),
	}
}
