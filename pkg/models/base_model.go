package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint            `json:"id" gorm:"primarykey; column:id; default:null"`
	CreatedAt *time.Time      `json:"created_at,omitempty" gorm:"column:created_at; default:null"`
	CreatedBy *uint           `json:"created_by,omitempty" gorm:"column:created_by; default:null"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" gorm:"column:updated_at; default:null"`
	UpdatedBy *uint           `json:"updated_by,omitempty" gorm:"column:updated_by; default:null"`
	DeletedBy *uint           `json:"deleted_by,omitempty" gorm:"column:deleted_by; default:null"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index; column:deleted_at; default:null"`
}
