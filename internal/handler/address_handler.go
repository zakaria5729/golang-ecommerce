package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/address"
	c "github.com/easy-comerce/backend/pkg/constants"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type AddressHandler struct {
	service *address.AddressService
}

func NewAddressHandler(service *address.AddressService) *AddressHandler {
	return &AddressHandler{
		service: service,
	}
}

func (h *AddressHandler) GetAllAddressesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	addressTypeFilter := q.Get(c.AddressAddressType)
	isDefaultFilter := q.Get(c.AddressIsDefault)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	addresses, err := h.service.GetAllAddressesByUser(*userID, addressTypeFilter, isDefaultFilter, sortBy, sortOrder)
	response.SendResponse(w, addresses, err, http.StatusInternalServerError)
}

func (h *AddressHandler) GetAllAddressesPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	userIdStr := q.Get(c.AddressUserID)
	addressTypeFilter := q.Get(c.AddressAddressType)
	isDefaultFilter := q.Get(c.AddressIsDefault)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	paginatedResponse, err := h.service.GetAllAddressesPaginated(showDeleted, userIdStr, pageStr, pageSizeStr, addressTypeFilter, isDefaultFilter, sortBy, sortOrder)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *AddressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	address, err := h.service.GetAddressByID(*id, showDeleted)
	response.SendResponse(w, address, err, http.StatusInternalServerError)
}

func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req address.Address
	if !utils.DecodeJSON(w, r, &req, "CreateAddress") {
		return
	}

	if validationErrors := validateAddressRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	address, err := h.service.CreateAddress(*userID, &req)
	response.SendResponse(w, address, err, http.StatusInternalServerError)
}

func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	var req address.Address
	if !utils.DecodeJSON(w, r, &req, "UpdateAddress") {
		return
	}

	if validationErrors := validateAddressRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	address, err := h.service.UpdateAddress(*id, *userID, &req)
	response.SendResponse(w, address, err, http.StatusInternalServerError)
}

func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAddress(*id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Address deleted successfully")
}

func (h *AddressHandler) RemoveAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveAddress(*id, *userID); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Address removed successfully")
}

func (h *AddressHandler) UndoDeleteAddress(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	if err := h.service.UndoDeleteAddress(*id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Undo delete address successfully")
}

func (h *AddressHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	addressType := r.URL.Query().Get(c.AddressAddressType)
	if addressType == "" {
		response.SendErrorJSON(w, "Address type is required", http.StatusBadRequest)
		return
	}

	if err := h.service.SetDefaultAddress(*id, *userID, addressType); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, "Address set as default successfully")
}

func (h *AddressHandler) GetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	addressType := r.URL.Query().Get(c.AddressAddressType)
	if addressType == "" {
		response.SendErrorJSON(w, "Address type is required", http.StatusBadRequest)
		return
	}

	address, err := h.service.GetDefaultAddress(*userID, addressType)
	response.SendResponse(w, address, err, http.StatusInternalServerError)
}

func validateAddressRequest(req *address.Address) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if req.Street == "" || len(req.Street) < 5 {
		errors.AddError("street", "Street must be at least 5 characters long")
	} else if len(req.Street) > 200 {
		errors.AddError("street", "Street must not exceed 200 characters")
	}

	if req.City == "" || len(req.City) < 2 {
		errors.AddError("city", "City must be at least 2 characters long")
	} else if len(req.City) > 100 {
		errors.AddError("city", "City must not exceed 100 characters")
	}

	if req.Country == "" || len(req.Country) < 2 {
		errors.AddError("country", "Country must be at least 2 characters long")
	} else if len(req.Country) > 100 {
		errors.AddError("country", "Country must not exceed 100 characters")
	}

	if req.State != nil && *req.State != "" {
		if len(*req.State) < 2 {
			errors.AddError("state", "State must be at least 2 characters long")
		} else if len(*req.State) > 100 {
			errors.AddError("state", "State must not exceed 100 characters")
		}
	}

	if req.ZipCode != nil && *req.ZipCode != "" {
		if len(*req.ZipCode) < 3 {
			errors.AddError("zip_code", "Zip code must be at least 3 characters long")
		} else if len(*req.ZipCode) > 20 {
			errors.AddError("zip_code", "Zip code must not exceed 20 characters")
		}
	}

	if req.AddressType != "" {
		if req.AddressType != c.AddressTypeBilling && req.AddressType != c.AddressTypeShipping {
			errors.AddError("address_type", "Address type must be one of: "+c.AddressTypeBilling+", "+c.AddressTypeShipping)
		}
	}

	return errors
}
