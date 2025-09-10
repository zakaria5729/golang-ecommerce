package permission

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"gorm.io/gorm"
)

type PermissionUserStatus struct {
	ID        uint            `gorm:"primarykey"`
	Banned    bool            `gorm:"column:banned"`
	Verified  bool            `gorm:"column:verified"`
	DeletedAt *gorm.DeletedAt `gorm:"index; column:deleted_at"`
}

func (PermissionUserStatus) TableName() string {
	return constants.TableUser
}
