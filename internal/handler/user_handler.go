package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type UserHandler struct {
	userUseCase *user.UserUseCase
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userUseCase: user.NewUserUseCase(),
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

	err = h.userUseCase.UpdateProfile(*userID, &req)
	if err != nil {
		logger.Logger.Error("Update profile failed", "method", "UpdateProfile", "error", err, "userID", *userID)
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

	if err := h.userUseCase.ChangePassword(currentUser.ID, currentUser.Password, &req); err != nil {
		logger.Logger.Error("Change password failed", "method", "ChangePassword", "error", err, "userID", currentUser.ID)
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

	userResponse, err := h.userUseCase.CreateUser(r.Context(), &req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, userResponse)

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authUserID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || authUserID == nil || *authUserID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	updateUserID, err := utils.ParseUint(r.PathValue(constants.FieldID))
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

	err = h.userUseCase.UpdateUser(*authUserID, *updateUserID, &req)
	if err != nil {
		logger.Logger.Error("Update user failed", "method", "UpdateUser", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendCommonResponseJSON(w, "Update user info successfully")
}

func (h *UserHandler) GetAllUsersPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)
	showDeleted := q.Get(constants.ShowDeleted)

	userResponses, err := h.userUseCase.GetAllUsersPaginated(includeStr, showDeleted, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all users failed", "method", "GetAllUsers", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, userResponses)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	showDeleted := r.URL.Query().Get(constants.ShowDeleted)
	userResponse, err := h.userUseCase.GetUserByID(*id, include, showDeleted)
	if err != nil {
		logger.Logger.Error("Get user by ID failed", "method", "GetUserByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, userResponse)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.DeleteUser(r.Context(), *id); err != nil {
		logger.Logger.Error("Delete user failed", "method", "DeleteUser", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "User deleted successfully")
}

func (h *UserHandler) UndoDeletedUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.UndoDeletedUser(r.Context(), *id); err != nil {
		logger.Logger.Error("Undo delete user failed", "method", "UndoDeletedUser", "error", err, "id", id)
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
