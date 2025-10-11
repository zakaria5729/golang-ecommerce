package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/category"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type CategoryHandler struct {
	service *category.CategoryService
}

func NewCategoryHandler(service *category.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) GetAllCategoriesWithSubcategoriesPublic(w http.ResponseWriter, r *http.Request) {
	showDeleted := false
	categoryResponses, err := getAllCategoriesWithSubcategoriesData(r, h.service, &showDeleted)
	response.SendResponse(w, categoryResponses, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) GetAllCategoriesPublic(w http.ResponseWriter, r *http.Request) {
	categoryResponses, err := getAllCategoriesData(r, h.service, nil)
	response.SendResponse(w, categoryResponses, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) GetCategoryByIdPublic(w http.ResponseWriter, r *http.Request) {
	categoryResponse, err := getCategoryDataById(r, h.service, nil)
	response.SendResponse(w, categoryResponse, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) GetAllCategoriesPaginatedPublic(w http.ResponseWriter, r *http.Request) {
	paginatedResponse, err := getAllCategoriesPaginatedData(r, h.service, nil)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) GetAllCategoriesWithSubcategories(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	categoryResponses, err := getAllCategoriesWithSubcategoriesData(r, h.service, showDeleted)
	response.SendResponse(w, categoryResponses, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	categoryResponses, err := getAllCategoriesData(r, h.service, showDeleted)
	response.SendResponse(w, categoryResponses, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) GetAllCategoriesPaginated(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	paginatedResponse, err := getAllCategoriesPaginatedData(r, h.service, showDeleted)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	categoryResponse, err := getCategoryDataById(r, h.service, showDeleted)
	response.SendResponse(w, categoryResponse, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req category.CreateCategoryRequest
	if !utils.DecodeJSON(w, r, &req, "CreateCategory") {
		return
	}

	if validationErrors := validateCreateCategoryRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	categoryResponse, err := h.service.CreateCategory(r.Context(), &req)
	response.SendResponse(w, categoryResponse, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req category.UpdateCategoryRequest
	if !utils.DecodeJSON(w, r, &req, "UpdateCategory") {
		return
	}

	if validationErrors := validateUpdateCategoryRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	categoryResponse, err := h.service.UpdateCategory(r.Context(), *id, &req)
	response.SendResponse(w, categoryResponse, err, http.StatusInternalServerError)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCategory(r.Context(), *id)
	response.SendResponse(w, "Category deleted successfully", err, http.StatusInternalServerError)
}

func (h *CategoryHandler) UndoDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletedCategory(r.Context(), *id)
	response.SendResponse(w, "Undo category deleted successfully", err, http.StatusInternalServerError)
}

func getAllCategoriesData(r *http.Request, service *category.CategoryService, showDeleted *bool) ([]category.CategoryResponse, error) {
	q := r.URL.Query()
	parentIDFilter := q.Get(c.CategoryParentID)
	priorityLimitFilter := q.Get(c.CategoryPriorityLimit)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	return service.GetAllCategories(showDeleted, parentIDFilter, priorityLimitFilter, sortBy, sortOrder)
}

func getAllCategoriesWithSubcategoriesData(r *http.Request, service *category.CategoryService, showDeleted *bool) ([]category.CategorySubcategoriesResponse, error) {
	q := r.URL.Query()
	subcategoryDepthFilter := q.Get(c.SubcategoryDepth)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	return service.GetAllCategoriesWithSubcategories(showDeleted, subcategoryDepthFilter, sortBy, sortOrder)
}

func getCategoryDataById(r *http.Request, service *category.CategoryService, showDeleted *bool) (*category.CategoryResponse, error) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return nil, errors.New("invalid category ID")
	}

	return service.GetCategoryByID(*id, showDeleted)
}

func getAllCategoriesPaginatedData(r *http.Request, service *category.CategoryService, showDeleted *bool) (*models.PaginatedResponse, error) {
	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	parentIDFilter := q.Get(c.CategoryParentID)
	showPriorityFilter := q.Get(c.CategoryPriority)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	return service.GetAllCategoriesPaginated(showDeleted, parentIDFilter, pageStr, pageSizeStr, showPriorityFilter, sortBy, sortOrder)
}

func validateUpdateCategoryRequest(req *category.UpdateCategoryRequest) validator.ValidationErrors {
	return validateRequest(req.Title, req.SubTitle, req.ParentID, req.Priority)
}

func validateCreateCategoryRequest(req *category.CreateCategoryRequest) validator.ValidationErrors {
	return validateRequest(req.Title, req.SubTitle, req.ParentID, req.Priority)
}

func validateRequest(title string, subTitle *string, parentID *uint, priority *uint) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if title == "" || len(title) < 2 {
		errors.AddError("title", "Title must be at least 2 characters long")
	} else if len(title) > 100 {
		errors.AddError("title", "Title must not exceed 100 characters")
	}

	if subTitle != nil && *subTitle != "" && len(*subTitle) > 200 {
		errors.AddError("sub_title", "Sub title must not exceed 200 characters")
	}

	if parentID != nil && *parentID == 0 {
		errors.AddError("parent_id", "Parent ID must be a positive integer")
	}

	if priority != nil && *priority == 0 {
		errors.AddError("priority", "Priority must be a positive integer")
	}

	return errors
}
