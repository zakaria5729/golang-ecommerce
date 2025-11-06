package model

import "github.com/easy-comerce/backend/pkg/base"

type CmsPageResponse struct {
	base.AuditEntity
	Tag         string               `json:"tag"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Sections    []CmsSectionResponse `json:"sections"`
}
