package entity

import (
	m "github.com/easy-comerce/backend/internal/cms/model"
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
)

type CmsPageEntity struct {
	base.AuditEntity
	Tag         string `json:"tag" gorm:"column:tag; not null"`
	Name        string `json:"name" gorm:"column:name; not null"`
	Description string `json:"description" gorm:"column:description"`
}

func (CmsPageEntity) TableName() string {
	return constants.TableCmsPage
}

func (p *CmsPageEntity) ToResponse() *m.CmsPageResponse {
	return &m.CmsPageResponse{
		AuditEntity: p.AuditEntity,
		Tag:         p.Tag,
		Name:        p.Name,
		Description: p.Description,
	}
}
