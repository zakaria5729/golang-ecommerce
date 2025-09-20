package address

import (
	"strings"

	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/constants"
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

// **REQUIRED
func (r *AddressRepository) GetAllAddressesByUser(userID uint, include []string, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]Address, error) {
	var addresses []Address

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))
	query = query.Where(c.AddressUserID+" = ?", userID)

	if addressType != nil && *addressType != "" {
		query = query.Where(c.AddressAddressType+" = ?", *addressType)
	}

	if isDefault != nil {
		query = query.Where(c.AddressIsDefault+" = ?", *isDefault)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	err := query.Find(&addresses).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch addresses", "method", "GetAllAddressesByUser", "error", err, "userID", userID, "include", include, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return addresses, err
}

// **REQUIRED
func (r *AddressRepository) GetAllAddressesPaginated(userID *uint, include []string, page, pageSize int, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]Address, int, error) {
	var addresses []Address
	var total int64

	selectFields := r.getSelectableFields(include)
	query := r.db.Select(strings.Join(selectFields, ", "))

	if userID != nil {
		query = query.Where(c.AddressUserID+" = ?", *userID)
	}

	if addressType != nil && *addressType != "" {
		query = query.Where(c.AddressAddressType+" = ?", *addressType)
	}

	if isDefault != nil {
		query = query.Where(c.AddressIsDefault+" = ?", *isDefault)
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Model(&Address{}).Count(&total).Error; err != nil {
		logger.Logger.Error("Failed to count addresses", "method", "GetAllAddressesPaginatedByUser", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&addresses).Error
	if err != nil {
		logger.Logger.Error("Failed to fetch addresses paginated", "method", "GetAllAddressesPaginatedByUser", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return addresses, int(total), err
}

// **REQUIRED
func (r *AddressRepository) GetAddressByID(id uint, userID *uint) (*Address, error) {
	var address Address
	query := r.db.Where(c.FieldID+" = ?", id)
	if userID != nil {
		query = query.Where(c.AddressUserID+" = ?", *userID)
	}

	if err := query.First(&address).Error; err != nil {
		logger.Logger.Error("Failed to fetch address by ID", "method", "GetAddressByID", "error", err, "id", id, "userID", userID)
		return nil, err
	}

	return &address, nil
}

// **REQUIRED
func (r *AddressRepository) CreateAddress(address *Address) error {
	err := r.db.Create(address).Error
	if err != nil {
		logger.Logger.Error("Failed to create address", "method", "CreateAddress", "error", err, "address", address)
	}
	return err
}

// **REQUIRED
func (r *AddressRepository) UpdateAddress(address *Address) error {
	err := r.db.Model(&Address{}).Where(c.FieldID+" = ?", address.ID).Updates(address).Error
	if err != nil {
		logger.Logger.Error("Failed to update address", "method", "UpdateAddress", "error", err, "address", address)
	}
	return err
}

func (r *AddressRepository) DeleteAddress(id uint, userID uint) error {
	err := r.db.Where(c.FieldID+" = ? AND "+c.AddressUserID+" = ?", id, userID).Delete(&Address{}).Error
	if err != nil {
		logger.Logger.Error("Failed to delete address", "method", "DeleteAddress", "error", err, "id", id, "userID", userID)
	}
	return err
}

// **REQUIRED
func (r *AddressRepository) SetDefaultAddress(id uint, userID uint, addressType string) error {
	err := r.db.Where(c.FieldID+" = ? AND "+c.AddressUserID+" = ?", id, userID).Update(c.AddressIsDefault, true).Error
	if err != nil {
		logger.Logger.Error("Failed to delete address", "method", "DeleteAddress", "error", err, "id", id, "userID", userID)
	}
	return err
}

// **REQUIRED
func (r *AddressRepository) AddressExists(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Address{}).Where(c.FieldID+" = ? AND "+c.AddressUserID+" = ?", id, userID).Count(&count).Error
	if err != nil {
		logger.Logger.Error("Failed to check if address exists", "method", "AddressExists", "error", err, "id", id, "userID", userID)
	}
	return count > 0, err
}

// **REQUIRED
func (r *AddressRepository) GetDefaultAddress(userID uint, addressType string) (*Address, error) {
	var address Address

	err := r.db.Where(c.AddressUserID+" = ? AND "+c.AddressAddressType+" = ? AND "+c.AddressIsDefault+" = ?", userID, addressType, true).First(&address).Error
	if err != nil {
		logger.Logger.Error("Failed to get default address", "method", "GetDefaultAddress", "error", err, "userID", userID, "addressType", addressType)
		return nil, err
	}

	return &address, nil
}

func (r *AddressRepository) getSelectableFields(include []string) []string {
	defaultFields := []string{c.FieldID, c.AddressUserID, c.AddressStreet, c.AddressCity, c.AddressCountry, c.AddressAddressType, c.FieldCreatedAt, c.FieldUpdatedAt}
	optionalFields := []string{c.AddressState, c.AddressZipCode, c.AddressIsDefault}
	return utils.BuildSelectFields(defaultFields, optionalFields, include)
}
