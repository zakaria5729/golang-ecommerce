package model

type CmsPageRequest struct {
	Tag         string               `json:"tag" validate:"required"`
	Name        string               `json:"name" validate:"required"`
	Description string               `json:"description"`
	Sections    *[]CmsSectionRequest `json:"sections"`
}
