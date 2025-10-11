package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/user"
	c "github.com/easy-comerce/backend/pkg/constants"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type UserHandler struct {
	service *user.UserService
}

func NewUserHandler(service *user.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	currentUser, err := m.GetUserFromContext(r.Context())
	response.SendResponse(w, currentUser, err, http.StatusInternalServerError)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req user.UpdateProfileRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateProfile") {
		return
	}

	validationErrors := validateUpdateProfileRequest(&req)
	if validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err = h.service.UpdateProfile(*userID, &req)
	response.SendResponse(w, "Profile updated successfully", err, http.StatusInternalServerError)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	currentUser, err := m.GetUserFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req user.ChangePasswordRequest
	if !utils.DecodeJSON(w, r, &req, "ChangePassword") {
		return
	}

	if validationErrors := validateChangePasswordRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err = h.service.ChangePassword(currentUser.ID, currentUser.Password, &req)
	response.SendResponse(w, "Password changed successfully", err, http.StatusInternalServerError)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if !utils.DecodeJSON(w, r, &req, "CreateUser") {
		return
	}

	if validationErrors := validateCreateUserRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	userResponse, err := h.service.CreateUser(r.Context(), &req)
	response.SendResponse(w, userResponse, err, http.StatusInternalServerError)

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authUserID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	updateUserID, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || updateUserID == nil || *updateUserID == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req user.UpdateUserRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateUser") {
		return
	}

	if validationErrors := validateUpdateUserRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err = h.service.UpdateUser(*authUserID, *updateUserID, &req)
	response.SendResponse(w, "Update user info successfully", err, http.StatusInternalServerError)
}

func (h *UserHandler) GetAllUsersPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(c.Include)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := q.Get(c.ShowDeleted)

	userResponses, err := h.service.GetAllUsersPaginated(includeStr, showDeleted, pageStr, pageSizeStr, sortBy, sortOrder)
	response.SendResponse(w, userResponses, err, http.StatusInternalServerError)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(c.Include)
	showDeleted := r.URL.Query().Get(c.ShowDeleted)
	userResponse, err := h.service.GetUserByID(*id, include, showDeleted)
	response.SendResponse(w, userResponse, err, http.StatusInternalServerError)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), *id)
	response.SendResponse(w, "User deleted successfully", err, http.StatusInternalServerError)
}

func (h *UserHandler) UndoDeletedUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedUser(r.Context(), *id)
	response.SendResponse(w, "Undo user deleted successfully", err, http.StatusInternalServerError)
}

func validateUpdateProfileRequest(req *user.UpdateProfileRequest) validator.ValidationErrors {
	errors := validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
	)

	return errors
}

func validateCreateUserRequest(req *user.CreateUserRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateRequired(req.Password, "password"),
		validator.ValidateMinLength(req.Password, "password", 8),
		validator.ValidatePositiveInteger(req.RoleID, "role_id"),
	)
}

func validateUpdateUserRequest(req *user.UpdateUserRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
	)
}

func validateChangePasswordRequest(req *user.ChangePasswordRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.CurrentPassword, "current_password"),
		validator.ValidatePassword(req.NewPassword, "new_password"),
	)
}
