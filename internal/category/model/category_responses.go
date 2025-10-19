package category

import "github.com/easy-comerce/backend/pkg/base"

type CategoryResponse struct {
	base.BaseEntity
	SubTitle *string `json:"sub_title"`
	ImageURL *string `json:"image_url"`
	ParentID *uint   `json:"parent_id"`
	Priority *uint   `json:"priority"`
	Title    string  `json:"title"`
}

type CategorySubcategoriesResponse struct {
	Category      CategoryResponse                `json:"category"`
	Subcategories []CategorySubcategoriesResponse `json:"subcategories"`
}
