package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	m "github.com/easy-comerce/backend/internal/cms/model"
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
)

type ImageArray []string

type CmsSectionEntity struct {
	base.AuditEntity
	Title       string     `json:"title" gorm:"column:title; not null"`
	Description string     `json:"description" gorm:"column:description"`
	Link        string     `json:"link" gorm:"column:link"`
	Content     string     `json:"content" gorm:"column:content; not null"`
	PageId      uint       `json:"page_id" gorm:"column:page_id; not null"`
	ImageUrls   ImageArray `json:"image_urls" gorm:"column:image_urls;type:jsonb;default:'[]'"`
}

func (CmsSectionEntity) TableName() string {
	return constants.TableCmsSection
}

func (p *CmsSectionEntity) ToResponse() *m.CmsSectionResponse {
	return &m.CmsSectionResponse{
		AuditEntity: p.AuditEntity,
		Title:       p.Title,
		Description: p.Description,
		Link:        p.Link,
		Content:     p.Content,
		PageId:      p.PageId,
		ImageUrls:   p.ImageUrls,
	}
}

func (p *CmsSectionEntity) SetImageUrls(urls []string) {
	p.ImageUrls = ImageArray(urls)
}

func (s *ImageArray) Scan(value any) error {
	if value == nil {
		*s = ImageArray([]string{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			return fmt.Errorf("unsupported type: %T", value)
		}
	}

	return json.Unmarshal(bytes, s)
}

func (s ImageArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(s)
}
