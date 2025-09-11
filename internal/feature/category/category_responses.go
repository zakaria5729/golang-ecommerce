package category

import "github.com/easy-comerce/backend/pkg/models"

type CategoryResponse struct {
	models.BaseModel
	Title    string  `json:"title"`
	SubTitle *string `json:"sub_title"`
	ImageURL *string `json:"image_url"`
	ParentID *uint   `json:"parent_id"`
	IsActive bool    `json:"is_active"`
	Priority *uint   `json:"priority"`
}
