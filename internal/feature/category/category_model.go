package category

import (
	"time"
)

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey; not null; column:id"`
	Title     string    `json:"title" gorm:"not null; column:title"`
	SubTitle  *string   `json:"sub_title" gorm:"column:sub_title"`
	ImageURL  *string   `json:"image_url" gorm:"column:image_url"`
	ParentID  *uint     `json:"parent_id" gorm:"column:parent_id"`
	IsActive  bool      `json:"is_active" gorm:"default:true; column:is_active;"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

const (
	CategoryID        = "id"
	CategoryTitle     = "title"
	CategorySubTitle  = "sub_title"
	CategoryImageURL  = "image_url"
	CategoryParentID  = "parent_id"
	CategoryIsActive  = "is_active"
	CategoryCreatedAt = "created_at"
	CategoryUpdatedAt = "updated_at"
)
