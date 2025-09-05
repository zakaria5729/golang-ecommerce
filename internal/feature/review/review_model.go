package review

import (
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Review struct {
	models.BaseModel
	ProductID uint    `json:"product_id" gorm:"not null; column:product_id"`
	UserID    uint    `json:"user_id" gorm:"not null; column:user_id"`
	Rating    int     `json:"rating" gorm:"not null; column:rating"`
	Comment   *string `json:"comment,omitempty" gorm:"column:comment"`
}

func (r *Review) Sanitize() {
	if r.Comment != nil && *r.Comment != "" {
		sanitized := utils.Trim(*r.Comment)
		r.Comment = &sanitized
	}
}

const (
	ReviewProductID = "product_id"
	ReviewUserID    = "user_id"
	ReviewRating    = "rating"
	ReviewComment   = "comment"
)

const (
	MinRating = 1
	MaxRating = 5
)
