package handler

import (
	"encoding/json"
	"net/http"

	category "github.com/easy-comerce/backend/internal/feature/category"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type CategoryHandler struct {
	useCase *category.CategoryUseCase
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		useCase: category.NewCategoryUseCase(),
	}
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	include := r.URL.Query().Get("include")
	parentID := r.URL.Query().Get("parent_id")
	isActive := r.URL.Query().Get("is_active")
	categories, err := h.useCase.GetAllCategories(include, parentID, isActive)
	if err != nil {
		response.JSONError(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	response.JSON(w, categories)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get("include")
	category, err := h.useCase.GetCategoryByID(*id, include)
	if err != nil {
		response.JSONError(w, "Category not found", http.StatusNotFound)
		return
	}

	response.JSON(w, category)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat category.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		response.JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreateCategory(&cat); err != nil {
		response.JSONError(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	response.JSON(w, cat)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var cat category.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		response.JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	cat.ID = *id
	if err := h.useCase.UpdateCategory(&cat); err != nil {
		response.JSONError(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	response.JSON(w, cat)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := utils.ParseUint(idStr)
	if err != nil || id == nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteCategory(*id); err != nil {
		response.JSONError(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
