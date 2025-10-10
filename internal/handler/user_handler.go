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
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	response.SendSuccessJSON(w, currentUser.ToResponse())
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req user.UpdateProfileRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateProfile") {
		return
	}

	validationErrors := h.validateUpdateProfileRequest(&req)
	if validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err = h.service.UpdateProfile(*userID, &req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendCommonResponseJSON(w, "Profile updated successfully")
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

	if validationErrors := h.validateChangePasswordRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err := h.service.ChangePassword(currentUser.ID, currentUser.Password, &req); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendCommonResponseJSON(w, "Password changed successfully")
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if !utils.DecodeJSON(w, r, &req, "CreateUser") {
		return
	}

	if validationErrors := h.validateCreateUserRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	userResponse, err := h.service.CreateUser(r.Context(), &req)
	response.SendResponse(w, userResponse, err, http.StatusInternalServerError)

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authUserID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || authUserID == nil || *authUserID == 0 {
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

	if validationErrors := h.validateUpdateUserRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err = h.service.UpdateUser(*authUserID, *updateUserID, &req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendCommonResponseJSON(w, "Update user info successfully")
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

	if err := h.service.DeleteUser(r.Context(), *id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "User deleted successfully")
}

func (h *UserHandler) UndoDeletedUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.service.UndoDeletedUser(r.Context(), *id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "Undo user deleted successfully")
}

func (h *UserHandler) validateUpdateProfileRequest(req *user.UpdateProfileRequest) validator.ValidationErrors {
	errors := validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
	)

	return errors
}

func (h *UserHandler) validateCreateUserRequest(req *user.CreateUserRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateRequired(req.Password, "password"),
		validator.ValidateMinLength(req.Password, "password", 8),
		validator.ValidatePositiveInteger(req.RoleID, "role_id"),
	)
}

func (h *UserHandler) validateUpdateUserRequest(req *user.UpdateUserRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
	)
}

func (h *UserHandler) validateChangePasswordRequest(req *user.ChangePasswordRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.CurrentPassword, "current_password"),
		validator.ValidatePassword(req.NewPassword, "new_password"),
	)
}
