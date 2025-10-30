package review

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ReviewEntity struct {
	base.AuditEntity
	ProductID uint    `gorm:"not null; column:product_id"`
	UserID    uint    `gorm:"not null; column:user_id"`
	Comment   *string `gorm:"column:comment"`
	Rating    int     `gorm:"not null; column:rating"`
}

func (ReviewEntity) TableName() string {
	return constants.TableReview
}

func (r *ReviewEntity) Sanitize() {
	if r.Comment != nil && *r.Comment != "" {
		sanitized := utils.Trim(*r.Comment)
		r.Comment = &sanitized
	}
}
