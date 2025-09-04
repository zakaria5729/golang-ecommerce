package address

import (
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type Address struct {
	models.BaseModel
	UserID      uint    `json:"user_id" gorm:"not null; column:user_id"`
	Street      string  `json:"street" gorm:"not null; column:street"`
	City        string  `json:"city" gorm:"not null; column:city"`
	State       *string `json:"state,omitempty" gorm:"column:state"`
	ZipCode     *string `json:"zip_code,omitempty" gorm:"column:zip_code"`
	Country     string  `json:"country" gorm:"not null; column:country"`
	IsDefault   bool    `json:"is_default,omitempty" gorm:"column:is_default"`
	AddressType string  `json:"address_type,omitempty" gorm:"column:address_type"`
}

func (a *Address) Sanitize() {
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

const (
	AddressUserID      = "user_id"
	AddressStreet      = "street"
	AddressCity        = "city"
	AddressState       = "state"
	AddressZipCode     = "zip_code"
	AddressCountry     = "country"
	AddressIsDefault   = "is_default"
	AddressAddressType = "address_type"
)

const (
	AddressTypeShipping = "shipping"
	AddressTypeBilling  = "billing"
)
