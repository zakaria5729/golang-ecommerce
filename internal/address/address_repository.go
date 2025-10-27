package address

import (
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AddressRepository interface {
	GetAllAddressesByUser(userID uint, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]AddressEntity, error)
	GetAllAddressesPaginated(showDeleted *bool, userID *uint, page, pageSize int, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]AddressEntity, int, error)
	GetAddressByID(id uint, userID *uint, showDeleted *bool) (*AddressEntity, error)
	CreateAddress(address *AddressEntity) error
	UpdateAddress(address *AddressEntity) error
	DeleteAddress(id uint) error
	RemoveAddress(id uint, userID uint) error
	UndoDeleteAddress(id uint) error
	SetDefaultAddress(id uint, userID uint, addressType string) error
	AddressExists(id uint, userID *uint, showDeleted *bool) (bool, error)
	GetDefaultAddress(userID uint, addressType string) (*AddressEntity, error)
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{
		db: db,
	}
}

func (r *addressRepository) GetAllAddressesByUser(userID uint, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]AddressEntity, error) {
	var addresses []AddressEntity
	query := r.db.Where(c.AddressUserID+" = ?", userID)

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
		l.Logger.Error("❌ Failed to fetch addresses", "method", "GetAllAddressesByUser", "error", err, "userID", userID, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
	}
	return addresses, err
}

func (r *addressRepository) GetAllAddressesPaginated(showDeleted *bool, userID *uint, page, pageSize int, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]AddressEntity, int, error) {
	var addresses []AddressEntity
	var total int64

	query := r.db.Model(&AddressEntity{})
	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

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

	if err := query.Count(&total).Error; err != nil {
		l.Logger.Error("❌ Failed to count addresses", "method", "GetAllAddressesPaginatedByUser", "error", err, "userID", userID, "page", page, "pageSize", pageSize, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&addresses).Error
	if err != nil {
		l.Logger.Error("❌ Failed to fetch addresses paginated", "method", "GetAllAddressesPaginatedByUser", "error", err, "userID", userID, "page", page, "pageSize", pageSize, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return addresses, int(total), err
}

func (r *addressRepository) GetAddressByID(id uint, userID *uint, showDeleted *bool) (*AddressEntity, error) {
	var address AddressEntity
	query := r.db.Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if userID != nil {
		query = query.Where(c.AddressUserID+" = ?", *userID)
	}

	if err := query.First(&address).Error; err != nil {
		l.Logger.Error("❌ Failed to fetch address by ID", "method", "GetAddressByID", "error", err, "id", id, "userID", userID)
		return nil, err
	}

	return &address, nil
}

func (r *addressRepository) CreateAddress(address *AddressEntity) error {
	err := r.db.Create(address).Error
	if err != nil {
		l.Logger.Error("❌ Failed to create address", "method", "CreateAddress", "error", err, "address", address)
	}
	return err
}

func (r *addressRepository) UpdateAddress(address *AddressEntity) error {
	err := r.db.Model(&AddressEntity{}).Where(c.FieldID+" = ?", address.ID).Updates(address).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update address", "method", "UpdateAddress", "error", err, "address", address)
	}
	return err
}

func (r *addressRepository) DeleteAddress(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&AddressEntity{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete address", "method", "DeleteAddress", "error", err, "id", id)
	}
	return err
}

func (r *addressRepository) RemoveAddress(id uint, userID uint) error {
	err := r.db.Where(c.FieldID+" = ? AND "+c.AddressUserID+" = ?", id, userID).Delete(&AddressEntity{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to remove address", "method", "RemoveAddress", "error", err, "id", id, "userID", userID)
	}
	return err
}

func (r *addressRepository) UndoDeleteAddress(id uint) error {
	err := r.db.Unscoped().Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil).Error
	if err != nil {
		l.Logger.Error("❌ Failed to undo delete address", "method", "UndoDeleteAddress", "error", err, "id", id)
	}
	return err
}

func (r *addressRepository) SetDefaultAddress(id uint, userID uint, addressType string) error {
	err := r.db.Where(c.FieldID+" = ? AND "+c.AddressUserID+" = ?", id, userID).Update(c.AddressIsDefault, true).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete address", "method", "DeleteAddress", "error", err, "id", id, "userID", userID)
	}
	return err
}

func (r *addressRepository) AddressExists(id uint, userID *uint, showDeleted *bool) (bool, error) {
	query := r.db.Model(&AddressEntity{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}
	if userID != nil {
		query = query.Where(c.AddressUserID+" = ?", *userID)
	}

	var address AddressEntity
	err := query.Select(c.FieldID).Take(&address).Error
	if err != nil || address.ID == 0 {
		l.Logger.Error("❌ Failed to check if address exists", "method", "AddressExists", "error", err, "id", id, "userID", userID)
		return false, err
	}

	return true, nil
}

func (r *addressRepository) GetDefaultAddress(userID uint, addressType string) (*AddressEntity, error) {
	var address AddressEntity

	err := r.db.Where(c.AddressUserID+" = ? AND "+c.AddressAddressType+" = ? AND "+c.AddressIsDefault+" = ?", userID, addressType, true).First(&address).Error
	if err != nil {
		l.Logger.Error("❌ Failed to get default address", "method", "GetDefaultAddress", "error", err, "userID", userID, "addressType", addressType)
		return nil, err
	}

	return &address, nil
}
