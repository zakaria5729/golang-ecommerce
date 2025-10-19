package model

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"gorm.io/gorm"
)

type PermissionUserStatus struct {
	ID           uint            `gorm:"primarykey"`
	DeletedAt    *gorm.DeletedAt `gorm:"index; column:deleted_at"`
	Banned       bool            `gorm:"column:banned"`
	Verified     bool            `gorm:"column:verified"`
	RefreshToken *string         `gorm:"column:refresh_token"`
}

func (PermissionUserStatus) TableName() string {
	return constants.TableUser
}
