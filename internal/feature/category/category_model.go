package category

import (
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title     string          `json:"title,omitempty" gorm:"not null; column:title"`
	SubTitle  *string         `json:"sub_title,omitempty" gorm:"column:sub_title"`
	ImageURL  *string         `json:"image_url,omitempty" gorm:"column:image_url"`
	ParentID  *uint           `json:"parent_id,omitempty" gorm:"column:parent_id"`
	IsActive  bool            `json:"is_active,omitempty" gorm:"default:true; column:is_active;"`
	DeletedAt *gorm.DeletedAt `json:"-" gorm:"-"`
}

func (c *Category) Sanitize() {
	if c.Title != "" {
		c.Title = utils.Trim(c.Title)
	}
	if c.SubTitle != nil && *c.SubTitle != "" {
		sanitized := utils.Trim(*c.SubTitle)
		c.SubTitle = &sanitized
	}
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
