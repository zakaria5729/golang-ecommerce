package model

import "github.com/easy-comerce/backend/pkg/base"

type CmsSectionResponse struct {
	base.AuditEntity
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Link        string   `json:"link"`
	Content     string   `json:"content"`
	PageId      uint     `json:"page_id"`
	ImageUrls   []string `json:"image_urls"`
}
