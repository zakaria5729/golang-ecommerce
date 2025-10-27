package category

import (
	"errors"
	"net/http"
	"strconv"

	m "github.com/easy-comerce/backend/internal/category/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type CategoryHandler interface {
	GetAllCategoriesWithSubcategoriesPublic(w http.ResponseWriter, r *http.Request)
	GetAllCategoriesPublic(w http.ResponseWriter, r *http.Request)
	GetCategoryByIdPublic(w http.ResponseWriter, r *http.Request)
	GetAllCategoriesPaginatedPublic(w http.ResponseWriter, r *http.Request)
	GetAllCategoriesWithSubcategories(w http.ResponseWriter, r *http.Request)
	GetAllCategoriesPaginated(w http.ResponseWriter, r *http.Request)
	GetAllCategories(w http.ResponseWriter, r *http.Request)
	GetCategoryByID(w http.ResponseWriter, r *http.Request)
	CreateCategory(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	UndoDeleteCategory(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) CategoryHandler {
	return &categoryHandler{
		service: service,
	}
}

func (h *categoryHandler) GetAllCategoriesWithSubcategoriesPublic(w http.ResponseWriter, r *http.Request) {
	showDeleted := false
	categoryResponses, err := getAllCategoriesWithSubcategoriesData(r, h.service, &showDeleted)
	response.SendResponse(w, categoryResponses, err, http.StatusInternalServerError)
}

func (h *categoryHandler) GetAllCategoriesPublic(w http.ResponseWriter, r *http.Request) {
	categoryResponses, err := getAllCategoriesData(r, h.service, nil)
	response.SendResponse(w, categoryResponses, err, http.StatusInternalServerError)
}

func (h *categoryHandler) GetCategoryByIdPublic(w http.ResponseWriter, r *http.Request) {
	categoryResponse, err := getCategoryDataById(r, h.service, nil)
	response.SendResponse(w, categoryResponse, err, http.StatusInternalServerError)
}

func (h *categoryHandler) GetAllCategoriesPaginatedPublic(w http.ResponseWriter, r *http.Request) {
	paginatedResponse, err := getAllCategoriesPaginatedData(r, h.service, nil)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *categoryHandler) GetAllCategoriesWithSubcategories(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	subcategoryDepthStr := r.URL.Query().Get(c.SubcategoryDepth)
	sortBy := r.URL.Query().Get(c.SortBy)
	sortOrder := r.URL.Query().Get(c.SortOrder)

	var depthPtr *int
	if subcategoryDepthStr != "" {
		depth, _ := strconv.Atoi(subcategoryDepthStr)
		depthPtr = &depth
	}

	categories, err := h.service.GetAllCategoriesWithSubcategories(showDeleted, depthPtr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch categories with subcategories", http.StatusInternalServerError)
		return
	}

	response.SendResponse(w, categories, nil, http.StatusOK)
}

func (h *categoryHandler) GetAllCategoriesPaginated(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	paginatedResponse, err := getAllCategoriesPaginatedData(r, h.service, showDeleted)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *categoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	categoryResponses, err := getAllCategoriesData(r, h.service, showDeleted)
	response.SendResponse(w, categoryResponses, err, http.StatusInternalServerError)
}

func (h *categoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	categoryResponse, err := getCategoryDataById(r, h.service, showDeleted)
	response.SendResponse(w, categoryResponse, err, http.StatusInternalServerError)
}

func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req m.CreateCategoryRequest
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

func (h *categoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req m.UpdateCategoryRequest
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

func (h *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCategory(r.Context(), *id)
	response.SendResponse(w, "Category deleted successfully", err, http.StatusInternalServerError)
}

func (h *categoryHandler) UndoDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeleteCategory(r.Context(), *id)
	if err != nil {
		response.SendErrorJSON(w, "Failed to undo delete category", http.StatusInternalServerError)
		return
	}

	response.SendResponse(w, "Category undeleted successfully", nil, http.StatusOK)
}

func getAllCategoriesData(r *http.Request, service CategoryService, showDeleted *bool) ([]m.CategoryResponse, error) {
	q := r.URL.Query()
	parentID, _ := utils.ParseUint(q.Get(c.CategoryParentID))
	priorityLimit, _ := utils.ParseInt(q.Get(c.CategoryPriorityLimit))
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	return service.GetAllCategories(showDeleted, parentID, priorityLimit, sortBy, sortOrder)
}

func getAllCategoriesWithSubcategoriesData(r *http.Request, service CategoryService, showDeleted *bool) ([]m.CategorySubcategoriesResponse, error) {
	q := r.URL.Query()
	subcategoryDepthStr := q.Get(c.SubcategoryDepth)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	var depthPtr *int
	if subcategoryDepthStr != "" {
		depth, _ := strconv.Atoi(subcategoryDepthStr)
		depthPtr = &depth
	}

	return service.GetAllCategoriesWithSubcategories(showDeleted, depthPtr, sortBy, sortOrder)
}

func getCategoryDataById(r *http.Request, service CategoryService, showDeleted *bool) (*m.CategoryResponse, error) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return nil, errors.New("invalid category ID")
	}

	return service.GetCategoryByID(*id, showDeleted)
}

func getAllCategoriesPaginatedData(r *http.Request, service CategoryService, showDeleted *bool) (*response.PaginatedResponse, error) {
	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	parentID, _ := utils.ParseUint(q.Get(c.CategoryParentID))
	priorityLimit, _ := utils.ParseInt(q.Get(c.CategoryPriorityLimit))
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	return service.GetAllCategoriesPaginated(showDeleted, parentID, pageStr, pageSizeStr, priorityLimit, sortBy, sortOrder)
}

func validateUpdateCategoryRequest(req *m.UpdateCategoryRequest) validator.ValidationErrors {
	return validateRequest(req.Title, req.SubTitle, req.ParentID, req.Priority)
}

func validateCreateCategoryRequest(req *m.CreateCategoryRequest) validator.ValidationErrors {
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
