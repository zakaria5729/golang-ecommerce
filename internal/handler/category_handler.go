package handler

import (
	"encoding/json"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/category"
	"github.com/easy-comerce/backend/pkg/constants"
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
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	parentIDFilter := q.Get(category.CategoryParentID)
	showPriorityFilter := q.Get(category.CategoryPriority)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	categories, err := h.useCase.GetAllCategories(includeStr, parentIDFilter, showPriorityFilter, sortBy, sortOrder)
	if err != nil {
		response.JSONError(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	response.JSON(w, categories)
}

func (h *CategoryHandler) GetAllCategoriesPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	parentIDFilter := q.Get(category.CategoryParentID)
	showPriorityFilter := q.Get(category.CategoryPriority)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetAllCategoriesPaginated(includeStr, parentIDFilter, pageStr, pageSizeStr, showPriorityFilter, sortBy, sortOrder)
	if err != nil {
		response.JSONError(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	response.JSON(w, paginatedResponse)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	category, err := h.useCase.GetCategoryByID(*id, include)
	if err != nil {
		response.JSONError(w, "Category not found", http.StatusNotFound)
		return
	}

	response.JSON(w, category)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req category.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	category, err := h.useCase.CreateCategory(&req)
	if err != nil {
		response.JSONError(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	response.JSONCreated(w, category)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req category.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	category, err := h.useCase.UpdateCategory(*id, &req)
	if err != nil {
		response.JSONError(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	response.JSON(w, category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteCategory(*id); err != nil {
		response.JSONError(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	response.JSONNoContent(w)
}

func (h *CategoryHandler) ToggleCategoryStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.useCase.ToggleCategoryStatus(*id)
	if err != nil {
		response.JSONError(w, "Failed to toggle category status", http.StatusInternalServerError)
		return
	}

	response.JSON(w, category)
}
