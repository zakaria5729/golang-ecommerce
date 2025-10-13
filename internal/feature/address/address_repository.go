package address

import (
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{
		db: db,
	}
}

func (r *AddressRepository) GetAllAddressesByUser(userID uint, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]Address, error) {
	var addresses []Address
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

func (r *AddressRepository) GetAllAddressesPaginated(showDeleted *bool, userID *uint, page, pageSize int, addressType *string, isDefault *bool, sortBy, sortOrder string) ([]Address, int, error) {
	var addresses []Address
	var total int64

	query := r.db.Model(&Address{})
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

func (r *AddressRepository) GetAddressByID(id uint, userID *uint, showDeleted *bool) (*Address, error) {
	var address Address
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

func (r *AddressRepository) CreateAddress(address *Address) error {
	err := r.db.Create(address).Error
	if err != nil {
		l.Logger.Error("❌ Failed to create address", "method", "CreateAddress", "error", err, "address", address)
	}
	return err
}

func (r *AddressRepository) UpdateAddress(address *Address) error {
	err := r.db.Model(&Address{}).Where(c.FieldID+" = ?", address.ID).Updates(address).Error
	if err != nil {
		l.Logger.Error("❌ Failed to update address", "method", "UpdateAddress", "error", err, "address", address)
	}
	return err
}

func (r *AddressRepository) DeleteAddress(id uint) error {
	err := r.db.Where(c.FieldID+" = ?", id).Delete(&Address{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete address", "method", "DeleteAddress", "error", err, "id", id)
	}
	return err
}

func (r *AddressRepository) RemoveAddress(id uint, userID uint) error {
	err := r.db.Where(c.FieldID+" = ? AND "+c.AddressUserID+" = ?", id, userID).Delete(&Address{}).Error
	if err != nil {
		l.Logger.Error("❌ Failed to remove address", "method", "RemoveAddress", "error", err, "id", id, "userID", userID)
	}
	return err
}

func (r *AddressRepository) UndoDeleteAddress(id uint) error {
	err := r.db.Unscoped().Where(c.FieldID+" = ?", id).Update(c.FieldDeletedAt, nil).Error
	if err != nil {
		l.Logger.Error("❌ Failed to undo delete address", "method", "UndoDeleteAddress", "error", err, "id", id)
	}
	return err
}

func (r *AddressRepository) SetDefaultAddress(id uint, userID uint, addressType string) error {
	err := r.db.Where(c.FieldID+" = ? AND "+c.AddressUserID+" = ?", id, userID).Update(c.AddressIsDefault, true).Error
	if err != nil {
		l.Logger.Error("❌ Failed to delete address", "method", "DeleteAddress", "error", err, "id", id, "userID", userID)
	}
	return err
}

func (r *AddressRepository) AddressExists(id uint, userID *uint, showDeleted *bool) (bool, error) {
	query := r.db.Model(&Address{}).Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}
	if userID != nil {
		query = query.Where(c.AddressUserID+" = ?", *userID)
	}

	var address Address
	err := query.Select(c.FieldID).Take(&address).Error
	if err != nil || address.ID == 0 {
		l.Logger.Error("❌ Failed to check if address exists", "method", "AddressExists", "error", err, "id", id, "userID", userID)
		return false, err
	}

	return true, nil
}

func (r *AddressRepository) GetDefaultAddress(userID uint, addressType string) (*Address, error) {
	var address Address

	err := r.db.Where(c.AddressUserID+" = ? AND "+c.AddressAddressType+" = ? AND "+c.AddressIsDefault+" = ?", userID, addressType, true).First(&address).Error
	if err != nil {
		l.Logger.Error("❌ Failed to get default address", "method", "GetDefaultAddress", "error", err, "userID", userID, "addressType", addressType)
		return nil, err
	}

	return &address, nil
}
