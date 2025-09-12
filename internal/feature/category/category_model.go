package category

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Category struct {
	models.BaseModel
	SubTitle *string `json:"sub_title,omitempty" gorm:"column:sub_title"`
	ImageURL *string `json:"image_url,omitempty" gorm:"-"`
	ParentID *uint   `json:"parent_id,omitempty" gorm:"column:parent_id"`
	Priority *uint   `json:"priority,omitempty" gorm:"column:priority"`
	PathKey  *string `json:"path_key,omitempty" gorm:"column:path_key"`
	Title    string  `json:"title" gorm:"not null; column:title"`
	IsActive bool    `json:"is_active" gorm:"column:is_active"`
}

func (Category) TableName() string {
	return constants.TableCategory
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
	CategoryTitle    = "title"
	CategorySubTitle = "sub_title"
	CategoryImageURL = "image_url"
	CategoryParentID = "parent_id"
	CategoryIsActive = "is_active"
	CategoryPriority = "priority"
)

func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		BaseModel: c.BaseModel,
		Title:     c.Title,
		SubTitle:  c.SubTitle,
		ParentID:  c.ParentID,
		IsActive:  c.IsActive,
		Priority:  c.Priority,
		ImageURL:  utils.BuildFullImageURL(c.PathKey),
	}
}
