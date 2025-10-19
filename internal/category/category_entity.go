package category

import (
	m "github.com/easy-comerce/backend/internal/category/model"
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type CategoryEntity struct {
	base.BaseEntity
	SubTitle *string `gorm:"column:sub_title"`
	ParentID *uint   `gorm:"column:parent_id"`
	Priority *uint   `gorm:"column:priority"`
	PathKey  *string `gorm:"column:path_key"`
	Title    string  `gorm:"not null; column:title"`
}

func (CategoryEntity) TableName() string {
	return constants.TableCategory
}

func (c *CategoryEntity) Sanitize() {
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

func (c *CategoryEntity) ToResponse() *m.CategoryResponse {
	return &m.CategoryResponse{
		BaseEntity: c.BaseEntity,
		Title:      c.Title,
		SubTitle:   c.SubTitle,
		ParentID:   c.ParentID,
		Priority:   c.Priority,
		ImageURL:   utils.BuildFullImageURL(config.GetStorageDomain(), c.PathKey),
	}
}
