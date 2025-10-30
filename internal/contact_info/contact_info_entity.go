package contact_info

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ContactInfoEntity struct {
	base.BaseEntity
	Name    string `gorm:"not null; column:name"`
	Email   string `gorm:"not null; column:email"`
	Phone   string `gorm:"not null; column:phone"`
	Message string `gorm:"not null; column:message"`
	Type    string `gorm:"not null; column:type"`
}

func (ContactInfoEntity) TableName() string {
	return constants.TableContactInfo
}

func (c *ContactInfoEntity) Sanitize() {
	if c.Name != "" {
		c.Name = utils.Trim(c.Name)
	}
	if c.Email != "" {
		c.Email = utils.Trim(c.Email)
	}
	if c.Phone != "" {
		c.Phone = utils.Trim(c.Phone)
	}
	if c.Message != "" {
		c.Message = utils.Trim(c.Message)
	}
	if c.Type != "" {
		c.Type = utils.Trim(c.Type)
	}
}
