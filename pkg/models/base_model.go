package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint            `json:"id" gorm:"primarykey; column:id"`
	CreatedAt *time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index; column:deleted_at"`
}
