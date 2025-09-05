package handler

import (
	"encoding/json"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/category"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
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
		response.SendErrorJSON(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categories)
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
		response.SendErrorJSON(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	include := r.URL.Query().Get(constants.Include)
	category, err := h.useCase.GetCategoryByID(*id, include)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SendSuccessJSON(w, category)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req category.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorJSON(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateCategoryRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	category, err := h.useCase.CreateCategory(&req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, category, http.StatusCreated)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req category.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorJSON(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validateCategoryRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	category, err := h.useCase.UpdateCategory(*id, &req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteCategory(*id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Category deleted successfully")
}

func (h *CategoryHandler) ToggleCategoryStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.useCase.ToggleCategoryStatus(*id)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, category)
}

func (h *CategoryHandler) validateCategoryRequest(req *category.Category) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if req.Title == "" {
		errors.AddError("title", "Title is required")
	} else {
		if len(req.Title) < 2 {
			errors.AddError("title", "Title must be at least 2 characters long")
		}
		if len(req.Title) > 100 {
			errors.AddError("title", "Title must not exceed 100 characters")
		}
	}

	if req.SubTitle != nil && *req.SubTitle != "" {
		if len(*req.SubTitle) > 200 {
			errors.AddError("sub_title", "Sub title must not exceed 200 characters")
		}
	}

	if req.ImageURL != nil && *req.ImageURL != "" {
		urlErrors := validator.ValidateURL(*req.ImageURL, "image_url")
		errors = append(errors, urlErrors...)
	}

	if req.ParentID != nil && *req.ParentID == 0 {
		errors.AddError("parent_id", "Parent ID must be a positive integer")
	}

	if req.Priority != nil && *req.Priority == 0 {
		errors.AddError("priority", "Priority must be a positive integer")
	}

	return errors
}
