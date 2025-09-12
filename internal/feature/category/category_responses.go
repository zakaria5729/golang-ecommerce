package category

import "github.com/easy-comerce/backend/pkg/models"

type CategoryResponse struct {
	models.BaseModel
	SubTitle *string `json:"sub_title"`
	ImageURL *string `json:"image_url"`
	ParentID *uint   `json:"parent_id"`
	Priority *uint   `json:"priority"`
	Title    string  `json:"title"`
	IsActive bool    `json:"is_active"`
}
