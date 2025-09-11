package category

type CreateCategoryRequest struct {
	Title    string  `json:"title" validate:"required,min=2,max=100"`
	SubTitle *string `json:"sub_title,omitempty" validate:"omitempty,max=200"`
	ParentID *uint   `json:"parent_id,omitempty" validate:"omitempty,min=1"`
	Priority *uint   `json:"priority,omitempty" validate:"omitempty,min=1"`
}

type UpdateCategoryRequest struct {
	Title    string  `json:"title,omitempty" validate:"omitempty,min=2,max=100"`
	SubTitle *string `json:"sub_title,omitempty" validate:"omitempty,max=200"`
	ParentID *uint   `json:"parent_id,omitempty" validate:"omitempty,min=1"`
	Priority *uint   `json:"priority,omitempty" validate:"omitempty,min=1"`
	IsActive *bool   `json:"is_active,omitempty"`
}
