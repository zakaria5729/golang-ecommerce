package handler

import (
	"encoding/json"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/address"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AddressHandler struct {
	useCase *address.AddressUseCase
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{
		useCase: address.NewAddressUseCase(),
	}
}

func (h *AddressHandler) getUserID() uint {
	// TODO: Get from JWT token or session in the future
	return uint(1)
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
		response.JSONError(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	response.JSON(w, addresses)
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
		response.JSONError(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	response.JSON(w, paginatedResponse)
}

func (h *AddressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	address, err := h.useCase.GetAddressByID(*id, userID, include)
	if err != nil {
		response.JSONError(w, "Address not found", http.StatusNotFound)
		return
	}

	response.JSON(w, address)
}

func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	var req address.Address
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	address, err := h.useCase.CreateAddress(userID, &req)
	if err != nil {
		response.JSONError(w, "Failed to create address", http.StatusInternalServerError)
		return
	}

	response.JSONCreated(w, address)
}

func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	var req address.Address
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	address, err := h.useCase.UpdateAddress(*id, userID, &req)
	if err != nil {
		response.JSONError(w, "Failed to update address", http.StatusInternalServerError)
		return
	}

	response.JSON(w, address)
}

func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteAddress(*id, userID); err != nil {
		response.JSONError(w, "Failed to delete address", http.StatusInternalServerError)
		return
	}

	response.JSONNoContent(w)
}

func (h *AddressHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	addressType := r.URL.Query().Get("type")
	if addressType == "" {
		response.JSONError(w, "Address type is required", http.StatusBadRequest)
		return
	}

	address, err := h.useCase.SetDefaultAddress(*id, userID, addressType)
	if err != nil {
		response.JSONError(w, "Failed to set default address", http.StatusInternalServerError)
		return
	}

	response.JSON(w, address)
}

func (h *AddressHandler) GetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID()

	addressType := r.URL.Query().Get("type")
	if addressType == "" {
		response.JSONError(w, "Address type is required", http.StatusBadRequest)
		return
	}

	address, err := h.useCase.GetDefaultAddress(userID, addressType)
	if err != nil {
		response.JSONError(w, "Failed to get default address", http.StatusInternalServerError)
		return
	}

	if address == nil {
		response.JSONError(w, "No default address found", http.StatusNotFound)
		return
	}

	response.JSON(w, address)
}
