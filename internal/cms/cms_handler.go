package cms

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/cms/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type CmsHandler interface {
	GetByTag(w http.ResponseWriter, r *http.Request)
	CreatePageSection(w http.ResponseWriter, r *http.Request)
}

type cmsHandler struct {
	service CmsService
}

func NewCmsHandler(service CmsService) CmsHandler {
	return &cmsHandler{service: service}
}

func (h *cmsHandler) GetByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue(c.PageTag)
	if tag == "" {
		http.Error(w, "Tag is required", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	title := q.Get(c.PageTitle)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	includeSections := utils.ParseBoolPtr(q.Get(c.IncludeSections))
	showDeleted := false

	page, err := h.service.GetByTag(tag, &title, includeSections, sortBy, sortOrder, &showDeleted)
	response.SendApiResponse(w, page, err)
}

func (h *cmsHandler) CreatePageSection(w http.ResponseWriter, r *http.Request) {
	var req model.CmsPageRequest
	if !utils.DecodeJSON(w, r, &req, "CreatePageSection") {
		return
	}

	if validationErrors := validateCreatePageSectionRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	pageResponse, err := h.service.CreatePageSection(r.Context(), &req)
	response.SendApiResponse(w, pageResponse, err)
}

func validateCreatePageSectionRequest(req *model.CmsPageRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Tag, "tag"),
		validator.ValidateRequired(req.Name, "name"),
	)
}

// func (h *CMSHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get(":id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		http.Error(w, "Invalid page ID", http.StatusBadRequest)
// 		return
// 	}

// 	var page model.UpdatePageRequest
// 	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if _, err := h.service.UpdatePage(r.Context(), uint(id), &page); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(page)
// }

// func (h *CMSHandler) GetPage(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get(":id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		http.Error(w, "Invalid page ID", http.StatusBadRequest)
// 		return
// 	}

// 	page, err := h.service.GetPageByID(r.Context(), uint(id))
// 	if err != nil {
// 		http.Error(w, "Page not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(page)
// }

// func (h *CMSHandler) GetPageBySlug(w http.ResponseWriter, r *http.Request) {
// 	vars := r.URL.Query()
// 	slug := vars.Get("slug")

// 	page, err := h.service.GetPageBySlug(r.Context(), slug)
// 	if err != nil {
// 		http.Error(w, "Page not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(page)
// }

// func (h *CMSHandler) ListPages(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	page, _ := strconv.Atoi(query.Get("page"))
// 	if page == 0 {
// 		page = 1
// 	}

// 	pageSize, _ := strconv.Atoi(query.Get("page_size"))
// 	if pageSize == 0 {
// 		pageSize = 10
// 	}

// 	var isPublished *bool
// 	if isPublishedStr := query.Get("is_published"); isPublishedStr != "" {
// 		val := strings.ToLower(isPublishedStr) == "true"
// 		isPublished = &val
// 	}

// 	pages, total, err := h.service.ListPages(r.Context(), &model.ListPagesRequest{
// 		Page:        page,
// 		PageSize:    pageSize,
// 		IsPublished: isPublished,
// 	})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"data":  pages,
// 		"total": total,
// 		"page":  page,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func (h *CMSHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get(":id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		http.Error(w, "Invalid page ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.service.DeletePage(r.Context(), uint(id)); err != nil {
// 		http.Error(w, "Failed to delete page: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// // Section handlers
// func (h *CMSHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
// 	var section model.CreateSectionRequest
// 	if err := json.NewDecoder(r.Body).Decode(&section); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if _, err := h.service.CreateSection(r.Context(), &section); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(section)
// }

// func (h *CMSHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get(":id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		http.Error(w, "Invalid section ID", http.StatusBadRequest)
// 		return
// 	}

// 	var section model.UpdateSectionRequest
// 	if err := json.NewDecoder(r.Body).Decode(&section); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if _, err := h.service.UpdateSection(r.Context(), uint(id), &section); err != nil {
// 		http.Error(w, "Failed to update section: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(section)
// }

// func (h *CMSHandler) GetSection(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get(":id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		http.Error(w, "Invalid section ID", http.StatusBadRequest)
// 		return
// 	}

// 	section, err := h.service.GetSectionByID(r.Context(), uint(id))
// 	if err != nil {
// 		http.Error(w, "Section not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(section)
// }

// func (h *CMSHandler) GetSectionByIdentifier(w http.ResponseWriter, r *http.Request) {
// 	vars := r.URL.Query()
// 	identifier := vars.Get("identifier")

// 	section, err := h.service.GetSectionByIdentifier(r.Context(), identifier)
// 	if err != nil {
// 		http.Error(w, "Section not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(section)
// }

// func (h *CMSHandler) ListSections(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	page, _ := strconv.Atoi(query.Get("page"))
// 	if page == 0 {
// 		page = 1
// 	}

// 	pageSize, _ := strconv.Atoi(query.Get("page_size"))
// 	if pageSize == 0 {
// 		pageSize = 10
// 	}

// 	var isActive *bool
// 	if isActiveStr := query.Get("is_active"); isActiveStr != "" {
// 		val := strings.ToLower(isActiveStr) == "true"
// 		isActive = &val
// 	}

// 	sections, total, err := h.service.ListSections(r.Context(), &model.ListSectionsRequest{
// 		Page:      page,
// 		PageSize:  pageSize,
// 		IsActive:  isActive,
// 		SortBy:    query.Get("sort_by"),
// 		SortOrder: query.Get("sort_order"),
// 	})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"data":  sections,
// 		"total": total,
// 		"page":  page,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func (h *CMSHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get(":id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		http.Error(w, "Invalid section ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.service.DeleteSection(r.Context(), uint(id)); err != nil {
// 		http.Error(w, "Failed to delete section: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
