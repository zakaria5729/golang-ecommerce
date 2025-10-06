package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint            `json:"id" gorm:"primarykey; column:id"`
	CreatedAt *time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	CreatedBy *uint           `json:"created_by,omitempty" gorm:"column:created_by"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" gorm:"column:updated_at"`
	UpdatedBy *uint           `json:"updated_by,omitempty" gorm:"column:updated_by"`
	DeletedBy *uint           `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index; column:deleted_at"`
}
