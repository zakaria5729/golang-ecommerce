package models

import (
	"time"
)

type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey; column:id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}
