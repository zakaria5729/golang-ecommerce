package handler

import (
	"errors"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/category"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
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

func (h *CategoryHandler) GetAllCategoriesWithSubcategoriesPublic(w http.ResponseWriter, r *http.Request) {
	isActive := true
	err, categoryResponses := h.getAllCategoriesWithSubcategoriesData(r, &isActive, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponses)
}

func (h *CategoryHandler) GetAllCategoriesPublic(w http.ResponseWriter, r *http.Request) {
	isActive := true
	err, categoryResponses := h.getAllCategoriesData(r, &isActive, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponses)
}

func (h *CategoryHandler) GetCategoryByIdPublic(w http.ResponseWriter, r *http.Request) {
	isActive := true
	err, categoryResponse := h.getCategoryDataById(r, &isActive, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponse)
}

func (h *CategoryHandler) GetAllCategoriesPaginatedPublic(w http.ResponseWriter, r *http.Request) {
	isActive := true
	err, paginatedResponse := h.getAllCategoriesPaginatedData(r, &isActive, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *CategoryHandler) GetAllCategoriesWithSubcategories(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	isActive := utils.ParseBoolPtr(q.Get(c.CategoryIsActive))
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))
	err, categoryResponses := h.getAllCategoriesWithSubcategoriesData(r, isActive, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponses)
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	isActive := utils.ParseBoolPtr(q.Get(c.CategoryIsActive))
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	err, categoryResponses := h.getAllCategoriesData(r, isActive, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponses)
}

func (h *CategoryHandler) GetAllCategoriesPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	isActive := utils.ParseBoolPtr(q.Get(c.CategoryIsActive))
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	err, paginatedResponse := h.getAllCategoriesPaginatedData(r, isActive, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	isActive := utils.ParseBoolPtr(q.Get(c.CategoryIsActive))
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	err, categoryResponse := h.getCategoryDataById(r, isActive, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponse)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req category.CreateCategoryRequest
	if !utils.DecodeJSON(w, r, &req, "CreateCategory") {
		return
	}

	if validationErrors := h.validateCreateCategoryRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	categoryResponse, err := h.useCase.CreateCategory(&req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponse, http.StatusCreated)
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

	if validationErrors := h.validateUpdateCategoryRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	categoryResponse, err := h.useCase.UpdateCategory(*id, &req)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponse)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
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

func (h *CategoryHandler) UndoDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.UndoDeletedCategory(*id); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Undo category deleted successfully")
}

func (h *CategoryHandler) ToggleCategoryIsActive(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	categoryResponse, err := h.useCase.ToggleCategoryIsActive(*id)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, categoryResponse)
}

func (h *CategoryHandler) getAllCategoriesData(r *http.Request, isActive *bool, showDeleted *bool) (error, []category.CategoryResponse) {
	q := r.URL.Query()
	includeStr := q.Get(c.Include)
	parentIDFilter := q.Get(c.CategoryParentID)
	priorityLimitFilter := q.Get(c.CategoryPriorityLimit)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	responses, err := h.useCase.GetAllCategories(isActive, showDeleted, includeStr, parentIDFilter, priorityLimitFilter, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories", "method", "getAllCategoriesData", "error", err, "include", includeStr, "parentID", parentIDFilter, "priorityLimit", priorityLimitFilter, "sortBy", sortBy, "sortOrder", sortOrder)
		return err, nil
	}

	return nil, responses
}

func (h *CategoryHandler) getAllCategoriesWithSubcategoriesData(r *http.Request, isActive *bool, showDeleted *bool) (error, []category.CategorySubcategoriesResponse) {
	q := r.URL.Query()
	includeStr := q.Get(c.Include)
	showPriorityFilter := q.Get(c.CategoryPriority)
	subcategoryDepthFilter := q.Get(c.SubcategoryDepth)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	responses, err := h.useCase.GetAllCategoriesWithSubcategories(isActive, showDeleted, includeStr, showPriorityFilter, subcategoryDepthFilter, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories", "method", "GetAllCategoriesWithSubcategories", "error", err, "include", includeStr, "showPriority", showPriorityFilter, "sortBy", sortBy, "sortOrder", sortOrder)
		return err, nil
	}

	return nil, responses
}

func (h *CategoryHandler) getCategoryDataById(r *http.Request, isActive *bool, showDeleted *bool) (error, *category.CategoryResponse) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid category ID"), nil
	}

	response, err := h.useCase.GetCategoryByID(*id, isActive, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to fetch category by ID", "method", "getCategoryDataById", "error", err, "id", id, "isActive", isActive, "showDeleted", showDeleted)
		return errors.New("category not found"), nil
	}

	return nil, response
}

func (h *CategoryHandler) getAllCategoriesPaginatedData(r *http.Request, isActive *bool, showDeleted *bool) (error, *models.PaginatedResponse) {
	q := r.URL.Query()
	includeStr := q.Get(c.Include)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	parentIDFilter := q.Get(c.CategoryParentID)
	showPriorityFilter := q.Get(c.CategoryPriority)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	paginatedResponse, err := h.useCase.GetAllCategoriesPaginated(isActive, showDeleted, includeStr, parentIDFilter, pageStr, pageSizeStr, showPriorityFilter, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch categories paginated", "method", "getAllCategoriesPaginatedData", "error", err, "include", includeStr, "parentID", parentIDFilter, "showPriority", showPriorityFilter, "sortBy", sortBy, "sortOrder", sortOrder)
		return err, nil
	}

	return nil, paginatedResponse
}

func (h *CategoryHandler) validateUpdateCategoryRequest(req *category.UpdateCategoryRequest) validator.ValidationErrors {
	return h.validateRequest(req.Title, req.SubTitle, req.ParentID, req.Priority)
}

func (h *CategoryHandler) validateCreateCategoryRequest(req *category.CreateCategoryRequest) validator.ValidationErrors {
	return h.validateRequest(req.Title, req.SubTitle, req.ParentID, req.Priority)
}

func (h *CategoryHandler) validateRequest(title string, subTitle *string, parentID *uint, priority *uint) validator.ValidationErrors {
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
