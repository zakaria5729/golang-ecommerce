package category

import (
	"time"
)

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	SubTitle  *string   `json:"sub_title,omitempty"`
	ImageURL  *string   `json:"image_url,omitempty"`
	ParentID  *uint     `json:"parent_id,omitempty"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
