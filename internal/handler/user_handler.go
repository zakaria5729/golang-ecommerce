package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
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
	currentUser, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	include := r.URL.Query().Get("include")
	profileResponse, err := h.userUseCase.GetUserProfile(currentUser.ID, include)
	if err != nil {
		logger.Logger.Error("Get profile failed", "method", "GetProfile", "error", err, "userID", currentUser.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, profileResponse)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	currentUser, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	include := r.URL.Query().Get("include")
	profileResponse, err := h.userUseCase.GetUserProfile(currentUser.ID, include)
	if err != nil {
		logger.Logger.Error("Get me failed", "method", "GetMe", "error", err, "userID", currentUser.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, profileResponse)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	currentUser, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req user.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode update profile request", "method", "UpdateProfile", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedUserResponse, err := h.userUseCase.UpdateUserProfile(currentUser.ID, &req)
	if err != nil {
		logger.Logger.Error("Update profile failed", "method", "UpdateProfile", "error", err, "userID", currentUser.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, updatedUserResponse)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	include := r.URL.Query().Get("include")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")

	// Parse pagination parameters
	var pageInt, pageSizeInt int
	var err error
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			pageInt = 1
		}
	} else {
		pageInt = 1
	}

	if pageSize != "" {
		pageSizeInt, err = strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
			pageSizeInt = 10
		}
	} else {
		pageSizeInt = 10
	}

	userResponses, err := h.userUseCase.GetAllUsers(include, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all users failed", "method", "GetAllUsers", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, userResponses)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.SendErrorJSON(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get("include")
	userResponse, err := h.userUseCase.GetUserByID(uint(id), include)
	if err != nil {
		logger.Logger.Error("Get user by ID failed", "method", "GetUserByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, userResponse)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.SendErrorJSON(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.DeleteUser(uint(id)); err != nil {
		logger.Logger.Error("Delete user failed", "method", "DeleteUser", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendDeleteJSON(w, "User deleted successfully")
}
