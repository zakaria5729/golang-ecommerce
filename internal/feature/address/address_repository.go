package address

import (
	"errors"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{
		db: db.GetDB(),
	}
}

func (r *AddressRepository) GetAllAddresses(userID uint, include []string, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]Address, error) {
	var addresses []Address

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(AddressUserID+" = ?", userID)

	if addressType != nil && *addressType != "" {
		query = query.Where(AddressAddressType+" = ?", *addressType)
	}

	if isDefault != nil {
		query = query.Where(AddressIsDefault+" = ?", *isDefault)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&addresses).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch addresses", "method", "GetAllAddresses", "error", err, "userID", userID, "include", include, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return addresses, err
}

func (r *AddressRepository) GetAllAddressesPaginated(userID uint, include []string, page, pageSize int, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]Address, int, error) {
	var addresses []Address
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(AddressUserID+" = ?", userID)

	if addressType != nil && *addressType != "" {
		query = query.Where(AddressAddressType+" = ?", *addressType)
	}

	if isDefault != nil {
		query = query.Where(AddressIsDefault+" = ?", *isDefault)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Model(&Address{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count addresses", "method", "GetAllAddressesPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&addresses).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch addresses paginated", "method", "GetAllAddressesPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return addresses, int(total), err
}

func (r *AddressRepository) GetAddressByID(id uint, userID uint, include []string) (*Address, error) {
	var address Address

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(AddressUserID+" = ?", userID)

	if err := query.Where(constants.FieldID+" = ?", id).First(&address).Error; err != nil {
		logger.Logger.Error("Failed to fetch address by ID", "method", "GetAddressByID", "error", err, "id", id, "userID", userID, "include", include)
		return nil, err
	}

	return &address, nil
}

func (r *AddressRepository) CreateAddress(address *Address) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// If this address is being set as default, unset other default addresses of the same type
		if address.IsDefault {
			if err := tx.Model(&Address{}).
				Where(AddressUserID+" = ? AND "+AddressAddressType+" = ? AND "+AddressIsDefault+" = ?",
					address.UserID, address.AddressType, true).
				Update(AddressIsDefault, false).Error; err != nil {
				return err
			}
		}

		return tx.Create(address).Error
	})
	if err != nil {
		logger.Logger.Error("Failed to create address", "method", "CreateAddress", "error", err, "address", address)
	}
	return err
}

func (r *AddressRepository) UpdateAddress(address *Address) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var currentAddress Address
		if err := tx.First(&currentAddress, address.ID).Error; err != nil {
			return err
		}

		// If this address is being set as default, unset other default addresses of the same type
		if address.IsDefault && !currentAddress.IsDefault {
			if err := tx.Model(&Address{}).
				Where(AddressUserID+" = ? AND "+AddressAddressType+" = ? AND "+AddressIsDefault+" = ? AND "+constants.FieldID+" != ?",
					address.UserID, address.AddressType, true, address.ID).
				Update(AddressIsDefault, false).Error; err != nil {
				return err
			}
		}

		return tx.Save(address).Error
	})
	if err != nil {
		logger.Logger.Error("Failed to update address", "method", "UpdateAddress", "error", err, "address", address)
	}
	return err
}

func (r *AddressRepository) DeleteAddress(id uint, userID uint) error {
	err := r.db.Where(constants.FieldID+" = ? AND "+AddressUserID+" = ?", id, userID).Delete(&Address{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete address", "method", "DeleteAddress", "error", err, "id", id, "userID", userID)
	}
	return err
}

func (r *AddressRepository) SetDefaultAddress(id uint, userID uint, addressType string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// First, unset all default addresses of this type for the user
		if err := tx.Model(&Address{}).
			Where(AddressUserID+" = ? AND "+AddressAddressType+" = ? AND "+AddressIsDefault+" = ?",
				userID, addressType, true).
			Update(AddressIsDefault, false).Error; err != nil {
			return err
		}

		// Then set the specified address as default
		if err := tx.Model(&Address{}).
			Where(constants.FieldID+" = ? AND "+AddressUserID+" = ? AND "+AddressAddressType+" = ?",
				id, userID, addressType).
			Update(AddressIsDefault, true).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Logger.Error("Failed to set default address", "method", "SetDefaultAddress", "error", err, "id", id, "userID", userID, "addressType", addressType)
	}
	return err
}

func (r *AddressRepository) AddressExists(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Address{}).Where(constants.FieldID+" = ? AND "+AddressUserID+" = ?", id, userID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if address exists", "method", "AddressExists", "error", err, "id", id, "userID", userID)
	}
	return count > 0, err
}

func (r *AddressRepository) GetDefaultAddress(userID uint, addressType string) (*Address, error) {
	var address Address
	err := r.db.Where(AddressUserID+" = ? AND "+AddressAddressType+" = ? AND "+AddressIsDefault+" = ?",
		userID, addressType, true).First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Info("No default address found", "method", "GetDefaultAddress", "userID", userID, "addressType", addressType)
			return nil, nil
		}
		logger.Logger.Error("Failed to get default address", "method", "GetDefaultAddress", "error", err, "userID", userID, "addressType", addressType)
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{constants.FieldID, AddressUserID, AddressStreet, AddressCity, AddressCountry, AddressAddressType, constants.FieldCreatedAt, constants.FieldUpdatedAt}
	optionalFields := []string{AddressState, AddressZipCode, AddressIsDefault}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
