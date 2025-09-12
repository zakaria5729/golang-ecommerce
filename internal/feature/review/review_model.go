package review

import (
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Review struct {
	models.BaseModel
	ProductID uint    `gorm:"not null; column:product_id"`
	UserID    uint    `gorm:"not null; column:user_id"`
	Comment   *string `gorm:"column:comment"`
	Rating    int     `gorm:"not null; column:rating"`
}

func (Review) TableName() string {
	return constants.TableReview
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
