package models

import (
	"time"

	c "github.com/easy-comerce/backend/pkg/constants"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint            `json:"id" gorm:"primarykey; column:id"`
	CreatedAt *time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" gorm:"column:updated_at"`
	UpdatedBy *uint           `json:"updated_by,omitempty" gorm:"column:updated_by"`
	DeletedBy *uint           `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index; column:deleted_at"`
}

func (b *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	if userID, ok := tx.Statement.Context.Value(c.UserIDContextKey).(uint); ok {
		tx.Statement.SetColumn(c.FieldDeletedBy, userID)
	}
	return
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	if userID, ok := tx.Statement.Context.Value(c.UserIDContextKey).(uint); ok {
		tx.Statement.SetColumn(c.FieldUpdatedBy, userID)
	}
	return
}
