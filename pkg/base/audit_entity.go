package base

import (
	"gorm.io/gorm"
)

type AuditEntity struct {
	BaseEntity
	CreatedBy *uint           `json:"created_by,omitempty" gorm:"column:created_by; default:null"`
	UpdatedBy *uint           `json:"updated_by,omitempty" gorm:"column:updated_by; default:null"`
	DeletedBy *uint           `json:"deleted_by,omitempty" gorm:"column:deleted_by; default:null"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index; column:deleted_at; default:null"`
}
