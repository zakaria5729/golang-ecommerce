package base

import (
	"time"
)

type BaseEntity struct {
	ID        uint       `json:"id" gorm:"primarykey; column:id; default:null"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at; default:null"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at; default:null"`
}
