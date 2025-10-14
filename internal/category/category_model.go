package category

import (
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Category struct {
	models.BaseModel
	SubTitle *string `gorm:"column:sub_title"`
	ParentID *uint   `gorm:"column:parent_id"`
	Priority *uint   `gorm:"column:priority"`
	PathKey  *string `gorm:"column:path_key"`
	Title    string  `gorm:"not null; column:title"`
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
	if c.PathKey != nil && *c.PathKey == "" {
		c.PathKey = nil
	}
}

func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		BaseModel: c.BaseModel,
		Title:     c.Title,
		SubTitle:  c.SubTitle,
		ParentID:  c.ParentID,
		Priority:  c.Priority,
		ImageURL:  utils.BuildFullImageURL(config.GetStorageDomain(), c.PathKey),
	}
}
