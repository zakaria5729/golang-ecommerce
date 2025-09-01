package handler

import (
	"net/http"
	"strconv"
	"strings"

	category "github.com/easy-comerce/backend/internal/feature/category"
	"github.com/easy-comerce/backend/pkg/response"
)

type CategoryHandler struct {
	useCase *category.CategoryUseCase
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		useCase: category.NewCategoryUseCase(),
	}
}

// GetAllCategories handles GET /api/v1/categories
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	include := r.URL.Query().Get("include")
	parentID := r.URL.Query().Get("parent_id")
	isActive := r.URL.Query().Get("is_active")

	// Get categories
	categories, err := h.useCase.GetAllCategories(include, parentID, isActive)
	if err != nil {
		response.JSONError(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	// Return the categories directly (database already filtered fields)
	response.JSON(w, categories)
}

// GetCategoryByID handles GET /api/v1/categories/{id}
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	idStr := pathParts[4] // /api/v1/categories/{id}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Get query parameters
	include := r.URL.Query().Get("include")

	// Get category
	category, err := h.useCase.GetCategoryByID(uint(id), include)
	if err != nil {
		response.JSONError(w, "Category not found", http.StatusNotFound)
		return
	}

	// Return the category directly (database already filtered fields)
	response.JSON(w, category)
}
