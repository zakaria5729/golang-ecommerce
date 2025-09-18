package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
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

// **REQUIRED
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	currentUser, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	response.SendSuccessJSON(w, currentUser.ToResponse())
}

// **REQUIRED
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
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

// **REQUIRED
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if !utils.DecodeJSON(w, r, &req, "CreateUser") {
		return
	}

	if validationErrors := h.validateCreateUserRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	userResponse, err := h.userUseCase.CreateUser(&req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, userResponse)

}

// **REQUIRED
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
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

	err = h.userUseCase.UpdateUser(*id, &req)
	if err != nil {
		logger.Logger.Error("Update user failed", "method", "UpdateUser", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendCommonResponseJSON(w, "Update user info successfully")
}

// **REQUIRED
func (h *UserHandler) GetAllUsersPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	userResponses, err := h.userUseCase.GetAllUsersPaginated(includeStr, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all users failed", "method", "GetAllUsers", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, userResponses)
}

// **REQUIRED
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	userResponse, err := h.userUseCase.GetUserByID(*id, include)
	if err != nil {
		logger.Logger.Error("Get user by ID failed", "method", "GetUserByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, userResponse)
}

// **REQUIRED
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.DeleteUser(*id); err != nil {
		logger.Logger.Error("Delete user failed", "method", "DeleteUser", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "User deleted successfully")
}

func (h *UserHandler) validateUpdateProfileRequest(req *user.UpdateProfileRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateRequired(*req.PathKey, "path_key"),
	)
}

func (h *UserHandler) validateCreateUserRequest(req *user.CreateUserRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateRequiredBool(req.Banned, "banned"),
		validator.ValidateRequiredBool(req.Verified, "verified"),
		validator.ValidateRequired(req.Password, "password"),
		validator.ValidateMinLength(req.Password, "password", 8),
		validator.ValidatePositiveInteger(req.RoleID, "role_id"),
	)
}

func (h *UserHandler) validateUpdateUserRequest(req *user.UpdateUserRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateRequiredBool(req.Banned, "banned"),
		validator.ValidateRequiredBool(req.Verified, "verified"),
		validator.ValidateRequired(req.Password, "password"),
		validator.ValidateMinLength(req.Password, "password", 8),
	)
}
