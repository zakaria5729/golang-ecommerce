package base

import (
	"time"
)

type BaseEntity struct {
	ID        uint       `json:"id" gorm:"primarykey; column:id;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at; autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at; autoUpdateTime"`
}
