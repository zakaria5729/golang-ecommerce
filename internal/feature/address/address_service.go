package address

import (
	"errors"
	"fmt"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AddressService struct {
	repo *AddressRepository
}

func NewAddressService(repo *AddressRepository) *AddressService {
	return &AddressService{
		repo: repo,
	}
}

func (s *AddressService) GetAllAddressesByUser(userID uint, addressTypeFilter string, isDefaultFilter string, sortBy, sortOrder string) ([]Address, error) {
	addressType := utils.ParseStringPtr(addressTypeFilter)
	isDefault := utils.ParseBoolPtr(isDefaultFilter)

	addresses, err := s.repo.GetAllAddressesByUser(userID, addressType, isDefault, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}

	return addresses, nil
}

func (s *AddressService) GetAllAddressesPaginated(showDeleted *bool, userIdStr string, pageStr string, pageSizeStr string, addressTypeFilter string, isDefaultFilter string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	addressType := utils.ParseStringPtr(addressTypeFilter)
	isDefault := utils.ParseBoolPtr(isDefaultFilter)
	userID, _ := utils.ParseUint(userIdStr)

	addresses, total, err := s.repo.GetAllAddressesPaginated(showDeleted, userID, page, pageSize, addressType, isDefault, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}

	var addressPtrs []*Address
	for i := range addresses {
		addressPtrs = append(addressPtrs, &addresses[i])
	}

	return utils.BuildPaginatedResponse(addressPtrs, total, page, pageSize), nil
}

func (s *AddressService) GetAddressByID(id uint, showDeleted *bool) (*Address, error) {
	address, err := s.repo.GetAddressByID(id, nil, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("address not found: %w", err)
	}

	return address, nil
}

func (s *AddressService) CreateAddress(userID uint, req *Address) (*Address, error) {
	req.UserID = userID
	req.Sanitize()

	if req.AddressType == "" {
		req.AddressType = c.AddressTypeShipping
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

	if err := s.repo.CreateAddress(address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return address, nil
}

func (s *AddressService) UpdateAddress(id uint, userID uint, req *Address) (*Address, error) {
	req.Sanitize()

	existingAddress, err := s.repo.GetAddressByID(id, &userID, nil)
	if err != nil {
		return nil, fmt.Errorf("address not found: %w", err)
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

	if err := s.repo.UpdateAddress(existingAddress); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return existingAddress, nil
}

func (s *AddressService) DeleteAddress(id uint) error {
	exists, err := s.repo.AddressExists(id, nil, nil)
	if err != nil || !exists {
		return fmt.Errorf("address not found: %w", err)
	}

	if err := s.repo.DeleteAddress(id); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	return nil
}

func (s *AddressService) RemoveAddress(id uint, userID uint) error {
	exists, err := s.repo.AddressExists(id, &userID, nil)
	if err != nil || !exists {
		return fmt.Errorf("address not found: %w", err)
	}

	if err := s.repo.RemoveAddress(id, userID); err != nil {
		return fmt.Errorf("failed to remove address: %w", err)
	}

	return nil
}

func (s *AddressService) UndoDeleteAddress(id uint) error {
	showDeleted := true
	exists, err := s.repo.AddressExists(id, nil, &showDeleted)
	if err != nil || !exists {
		return fmt.Errorf("address not found: %w", err)
	}

	if err := s.repo.UndoDeleteAddress(id); err != nil {
		return fmt.Errorf("failed to undo delete address: %w", err)
	}

	return nil
}

func (s *AddressService) SetDefaultAddress(id uint, userID uint, addressType string) error {
	if addressType != c.AddressTypeShipping && addressType != c.AddressTypeBilling {
		return errors.New("invalid address type. Must be " + c.AddressTypeShipping + " or " + c.AddressTypeBilling)
	}

	exists, err := s.repo.AddressExists(id, &userID, nil)
	if err != nil || !exists {
		return fmt.Errorf("address not found: %w", err)
	}

	if err := s.repo.SetDefaultAddress(id, userID, addressType); err != nil {
		return fmt.Errorf("failed to set default address: %w", err)
	}

	return nil
}

func (s *AddressService) GetDefaultAddress(userID uint, addressType string) (*Address, error) {
	if addressType != c.AddressTypeShipping && addressType != c.AddressTypeBilling {
		return nil, errors.New("invalid address type. Must be " + c.AddressTypeShipping + " or " + c.AddressTypeBilling)
	}

	address, err := s.repo.GetDefaultAddress(userID, addressType)
	if err != nil {
		return nil, fmt.Errorf("failed to get default address: %w", err)
	}

	return address, nil
}
