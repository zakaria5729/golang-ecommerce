package address

import (
	"github.com/easy-comerce/backend/pkg/base"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AddressEntity struct {
	base.AuditEntity
	UserID      uint    `gorm:"not null; column:user_id"`
	State       *string `gorm:"column:state"`
	ZipCode     *string `gorm:"column:zip_code"`
	Street      string  `gorm:"not null; column:street"`
	City        string  `gorm:"not null; column:city"`
	Country     string  `gorm:"not null; column:country"`
	AddressType string  `gorm:"column:address_type"`
	IsDefault   bool    `gorm:"column:is_default"`
}

func (AddressEntity) TableName() string {
	return constants.TableAddress
}

func (a *AddressEntity) Sanitize() {
	if a.Street != "" {
		a.Street = utils.Trim(a.Street)
	}
	if a.City != "" {
		a.City = utils.Trim(a.City)
	}
	if a.State != nil && *a.State != "" {
		sanitized := utils.Trim(*a.State)
		a.State = &sanitized
	}
	if a.ZipCode != nil && *a.ZipCode != "" {
		sanitized := utils.Trim(*a.ZipCode)
		a.ZipCode = &sanitized
	}
	if a.Country != "" {
		a.Country = utils.Trim(a.Country)
	}
	if a.AddressType != "" {
		a.AddressType = utils.Trim(a.AddressType)
	}
}
