package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	userPkg "github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
)

type UserHandler struct {
	userUseCase *userPkg.UserUseCase
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userUseCase: userPkg.NewUserUseCase(),
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	include := r.URL.Query().Get("include")
	profile, err := h.userUseCase.GetUserProfile(user.ID, include)
	if err != nil {
		logger.Logger.Error("Get profile failed", "method", "GetProfile", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, profile)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	include := r.URL.Query().Get("include")
	profile, err := h.userUseCase.GetUserProfile(user.ID, include)
	if err != nil {
		logger.Logger.Error("Get me failed", "method", "GetMe", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, profile)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req userPkg.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Failed to decode update profile request", "method", "UpdateProfile", "error", err)
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the ID for update
	req.ID = user.ID

	updatedUser, err := h.userUseCase.UpdateUserProfile(user.ID, &req)
	if err != nil {
		logger.Logger.Error("Update profile failed", "method", "UpdateProfile", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendSuccessJSON(w, updatedUser)
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

	users, err := h.userUseCase.GetAllUsers(include, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Get all users failed", "method", "GetAllUsers", "error", err)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, users)
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
	user, err := h.userUseCase.GetUserByID(uint(id), include)
	if err != nil {
		logger.Logger.Error("Get user by ID failed", "method", "GetUserByID", "error", err, "id", id)
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.SendErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
