package address

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AddressUseCase struct {
	repo *AddressRepository
}

func NewAddressUseCase() *AddressUseCase {
	return &AddressUseCase{
		repo: NewAddressRepository(),
	}
}

func (uc *AddressUseCase) GetAllAddresses(userID uint, includeStr string, addressTypeFilter string, isDefaultFilter string, sortBy, sortOrder string) ([]Address, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	addressType := utils.ParseStringPtr(addressTypeFilter)
	isDefault := utils.ParseBoolPtr(isDefaultFilter)

	addresses, err := uc.repo.GetAllAddresses(userID, include, addressType, isDefault, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch addresses", "method", "GetAllAddresses", "error", err, "userID", userID, "include", include, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}

	return addresses, nil
}

func (uc *AddressUseCase) GetAllAddressesPaginated(userID uint, includeStr string, pageStr string, pageSizeStr string, addressTypeFilter string, isDefaultFilter string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	include := utils.ParseCommaSeparatedString(includeStr)
	addressType := utils.ParseStringPtr(addressTypeFilter)
	isDefault := utils.ParseBoolPtr(isDefaultFilter)

	addresses, total, err := uc.repo.GetAllAddressesPaginated(userID, include, page, pageSize, addressType, isDefault, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch addresses paginated", "method", "GetAllAddressesPaginated", "error", err, "userID", userID, "include", include, "page", page, "pageSize", pageSize, "addressType", addressType, "isDefault", isDefault, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}

	var addressPtrs []*Address
	for i := range addresses {
		addressPtrs = append(addressPtrs, &addresses[i])
	}

	return utils.BuildPaginatedResponse(addressPtrs, total, page, pageSize), nil
}

func (uc *AddressUseCase) GetAddressByID(id uint, userID uint, includeStr string) (*Address, error) {
	if id == 0 {
		logger.Logger.Error("Invalid address ID provided", "method", "GetAddressByID", "id", id)
		return nil, errors.New("invalid address ID")
	}

	if userID == 0 {
		logger.Logger.Error("Invalid user ID provided", "method", "GetAddressByID", "userID", userID)
		return nil, errors.New("invalid user ID")
	}

	include := utils.ParseCommaSeparatedString(includeStr)
	address, err := uc.repo.GetAddressByID(id, userID, include)
	if err != nil {
		logger.Logger.Error("Address not found", "method", "GetAddressByID", "error", err, "id", id, "userID", userID, "include", include)
		return nil, fmt.Errorf("address not found: %w", err)
	}

	return address, nil
}

func (uc *AddressUseCase) CreateAddress(userID uint, req *Address) (*Address, error) {
	if userID == 0 {
		logger.Logger.Error("Invalid user ID provided", "method", "CreateAddress", "userID", userID)
		return nil, errors.New("invalid user ID")
	}

	req.Sanitize()
	req.UserID = userID

	if req.AddressType == "" {
		req.AddressType = AddressTypeShipping
	}

	if req.AddressType != AddressTypeShipping && req.AddressType != AddressTypeBilling {
		logger.Logger.Error("Invalid address type", "method", "CreateAddress", "addressType", req.AddressType)
		return nil, errors.New("invalid address type. Must be 'shipping' or 'billing'")
	}

	address := &Address{
		UserID:      req.UserID,
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		ZipCode:     req.ZipCode,
		Country:     req.Country,
		IsDefault:   req.IsDefault,
		AddressType: req.AddressType,
	}

	if err := uc.repo.CreateAddress(address); err != nil {
		logger.Logger.Error("Failed to create address", "method", "CreateAddress", "error", err, "address", address)
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return address, nil
}

func (uc *AddressUseCase) UpdateAddress(id uint, userID uint, req *Address) (*Address, error) {
	if id == 0 {
		logger.Logger.Error("Invalid address ID", "method", "UpdateAddress", "id", id)
		return nil, errors.New("invalid address ID")
	}

	if userID == 0 {
		logger.Logger.Error("Invalid user ID", "method", "UpdateAddress", "userID", userID)
		return nil, errors.New("invalid user ID")
	}

	req.Sanitize()

	existingAddress, err := uc.repo.GetAddressByID(id, userID, nil)
	if err != nil {
		logger.Logger.Error("Address not found", "method", "UpdateAddress", "error", err, "id", id, "userID", userID)
		return nil, fmt.Errorf("address not found: %w", err)
	}

	if req.AddressType != "" && req.AddressType != AddressTypeShipping && req.AddressType != AddressTypeBilling {
		logger.Logger.Error("Invalid address type", "method", "UpdateAddress", "addressType", req.AddressType)
		return nil, errors.New("invalid address type. Must be 'shipping' or 'billing'")
	}

	if req.Street != "" {
		existingAddress.Street = req.Street
	}
	if req.City != "" {
		existingAddress.City = req.City
	}
	if req.State != nil {
		existingAddress.State = req.State
	}
	if req.ZipCode != nil {
		existingAddress.ZipCode = req.ZipCode
	}
	if req.Country != "" {
		existingAddress.Country = req.Country
	}
	if req.AddressType != "" {
		existingAddress.AddressType = req.AddressType
	}
	if req.IsDefault != existingAddress.IsDefault {
		existingAddress.IsDefault = req.IsDefault
	}

	if err := uc.repo.UpdateAddress(existingAddress); err != nil {
		logger.Logger.Error("Failed to update address", "method", "UpdateAddress", "error", err, "address", existingAddress)
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return existingAddress, nil
}

func (uc *AddressUseCase) DeleteAddress(id uint, userID uint) error {
	if id == 0 {
		logger.Logger.Error("Invalid address ID", "method", "DeleteAddress", "id", id)
		return errors.New("invalid address ID")
	}

	if userID == 0 {
		logger.Logger.Error("Invalid user ID", "method", "DeleteAddress", "userID", userID)
		return errors.New("invalid user ID")
	}

	_, err := uc.repo.GetAddressByID(id, userID, nil)
	if err != nil {
		logger.Logger.Error("Address not found", "method", "DeleteAddress", "error", err, "id", id, "userID", userID)
		return fmt.Errorf("address not found: %w", err)
	}

	if err := uc.repo.DeleteAddress(id, userID); err != nil {
		logger.Logger.Error("Failed to delete address", "method", "DeleteAddress", "error", err, "id", id, "userID", userID)
		return fmt.Errorf("failed to delete address: %w", err)
	}

	return nil
}

func (uc *AddressUseCase) SetDefaultAddress(id uint, userID uint, addressType string) (*Address, error) {
	if id == 0 {
		logger.Logger.Error("Invalid address ID", "method", "SetDefaultAddress", "id", id)
		return nil, errors.New("invalid address ID")
	}

	if userID == 0 {
		logger.Logger.Error("Invalid user ID", "method", "SetDefaultAddress", "userID", userID)
		return nil, errors.New("invalid user ID")
	}

	if addressType != AddressTypeShipping && addressType != AddressTypeBilling {
		logger.Logger.Error("Invalid address type", "method", "SetDefaultAddress", "addressType", addressType)
		return nil, errors.New("invalid address type. Must be 'shipping' or 'billing'")
	}

	_, err := uc.repo.GetAddressByID(id, userID, nil)
	if err != nil {
		logger.Logger.Error("Address not found", "method", "SetDefaultAddress", "error", err, "id", id, "userID", userID)
		return nil, fmt.Errorf("address not found: %w", err)
	}

	if err := uc.repo.SetDefaultAddress(id, userID, addressType); err != nil {
		logger.Logger.Error("Failed to set default address", "method", "SetDefaultAddress", "error", err, "id", id, "userID", userID, "addressType", addressType)
		return nil, fmt.Errorf("failed to set default address: %w", err)
	}

	address, err := uc.repo.GetAddressByID(id, userID, nil)
	if err != nil {
		logger.Logger.Error("Failed to fetch updated address", "method", "SetDefaultAddress", "error", err, "id", id, "userID", userID)
		return nil, fmt.Errorf("failed to fetch updated address: %w", err)
	}

	return address, nil
}

func (uc *AddressUseCase) GetDefaultAddress(userID uint, addressType string) (*Address, error) {
	if userID == 0 {
		logger.Logger.Error("Invalid user ID", "method", "GetDefaultAddress", "userID", userID)
		return nil, errors.New("invalid user ID")
	}

	if addressType != AddressTypeShipping && addressType != AddressTypeBilling {
		logger.Logger.Error("Invalid address type", "method", "GetDefaultAddress", "addressType", addressType)
		return nil, errors.New("invalid address type. Must be 'shipping' or 'billing'")
	}

	address, err := uc.repo.GetDefaultAddress(userID, addressType)
	if err != nil {
		logger.Logger.Error("Failed to get default address", "method", "GetDefaultAddress", "error", err, "userID", userID, "addressType", addressType)
		return nil, fmt.Errorf("failed to get default address: %w", err)
	}

	return address, nil
}
