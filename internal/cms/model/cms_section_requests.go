package model

type CmsSectionRequest struct {
	ID          *uint    `json:"id"`
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Link        string   `json:"link"`
	ImageUrls   []string `json:"image_urls"`
}
