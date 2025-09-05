package handler

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/easy-comerce/backend/internal/feature/address"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AddressHandler struct {
	useCase *address.AddressUseCase
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{
		useCase: address.NewAddressUseCase(),
	}
}

func (h *AddressHandler) GetAllAddresses(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	addressTypeFilter := q.Get(address.AddressAddressType)
	isDefaultFilter := q.Get(address.AddressIsDefault)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	addresses, err := h.useCase.GetAllAddresses(userID, includeStr, addressTypeFilter, isDefaultFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, addresses)
}

func (h *AddressHandler) GetAllAddressesPaginated(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	addressTypeFilter := q.Get(address.AddressAddressType)
	isDefaultFilter := q.Get(address.AddressIsDefault)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetAllAddressesPaginated(userID, includeStr, pageStr, pageSizeStr, addressTypeFilter, isDefaultFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *AddressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	address, err := h.useCase.GetAddressByID(*id, userID, include)
	if err != nil {
		response.SendErrorJSON(w, "Address not found", http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, address)
}

func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	var req address.Address
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorJSON(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateAddressRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	address, err := h.useCase.CreateAddress(userID, &req)
	if err != nil {
		response.SendErrorJSON(w, "Failed to create address", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, address, http.StatusCreated)
}

func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	var req address.Address
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorJSON(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateAddressRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	address, err := h.useCase.UpdateAddress(*id, userID, &req)
	if err != nil {
		response.SendErrorJSON(w, "Failed to update address", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, address)
}

func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteAddress(*id, userID); err != nil {
		response.SendErrorJSON(w, "Failed to delete address")
		return
	}

	response.SendDeleteJSON(w, "Address deleted successfully")
}

func (h *AddressHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	addressType := r.URL.Query().Get("type")
	if addressType == "" {
		response.SendErrorJSON(w, "Address type is required", http.StatusBadRequest)
		return
	}

	address, err := h.useCase.SetDefaultAddress(*id, userID, addressType)
	if err != nil {
		response.SendErrorJSON(w, "Failed to set default address", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, address)
}

func (h *AddressHandler) GetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	addressType := r.URL.Query().Get("type")
	if addressType == "" {
		response.SendErrorJSON(w, "Address type is required", http.StatusBadRequest)
		return
	}

	address, err := h.useCase.GetDefaultAddress(userID, addressType)
	if err != nil {
		response.SendErrorJSON(w, "Failed to get default address", http.StatusInternalServerError)
		return
	}

	if address == nil {
		response.SendErrorJSON(w, "No default address found", http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, address)
}

func (h *AddressHandler) validateAddressRequest(req *address.Address) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if req.Street == "" {
		errors.AddError("street", "Street is required")
	} else {
		if len(req.Street) < 5 {
			errors.AddError("street", "Street must be at least 5 characters long")
		}
		if len(req.Street) > 200 {
			errors.AddError("street", "Street must not exceed 200 characters")
		}
	}

	if req.City == "" {
		errors.AddError("city", "City is required")
	} else {
		if len(req.City) < 2 {
			errors.AddError("city", "City must be at least 2 characters long")
		}
		if len(req.City) > 100 {
			errors.AddError("city", "City must not exceed 100 characters")
		}
	}

	if req.Country == "" {
		errors.AddError("country", "Country is required")
	} else {
		if len(req.Country) < 2 {
			errors.AddError("country", "Country must be at least 2 characters long")
		}
		if len(req.Country) > 100 {
			errors.AddError("country", "Country must not exceed 100 characters")
		}
	}

	if req.State != nil && *req.State != "" {
		if len(*req.State) < 2 {
			errors.AddError("state", "State must be at least 2 characters long")
		}
		if len(*req.State) > 100 {
			errors.AddError("state", "State must not exceed 100 characters")
		}
	}

	if req.ZipCode != nil && *req.ZipCode != "" {
		if len(*req.ZipCode) < 3 {
			errors.AddError("zip_code", "Zip code must be at least 3 characters long")
		}
		if len(*req.ZipCode) > 20 {
			errors.AddError("zip_code", "Zip code must not exceed 20 characters")
		}
	}

	if req.AddressType != "" {
		validTypes := []string{"billing", "shipping", "both"}
		isValid := slices.Contains(validTypes, req.AddressType)
		if !isValid {
			errors.AddError("address_type", "Address type must be one of: billing, shipping, both")
		}
	}

	return errors
}

func (h *AddressHandler) getUserID() uint {
	// TODO: Get from JWT token or session in the future
	return uint(1)
}
