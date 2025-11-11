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
	GetPageByTag(w http.ResponseWriter, r *http.Request)
	GetPageByID(w http.ResponseWriter, r *http.Request)
	GetPagesPaginated(w http.ResponseWriter, r *http.Request)
	CreatePageSection(w http.ResponseWriter, r *http.Request)
	UpdatePageSection(w http.ResponseWriter, r *http.Request)
	DeletePageSection(w http.ResponseWriter, r *http.Request)
	UndoDeletePageSection(w http.ResponseWriter, r *http.Request)
}

type cmsHandler struct {
	service CmsService
}

func NewCmsHandler(service CmsService) CmsHandler {
	return &cmsHandler{service: service}
}

func (h *cmsHandler) GetPageByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue(c.PageTag)
	if tag == "" {
		http.Error(w, "Tag is required", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	title := q.Get(c.SectionTitle)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	includeSections := utils.ParseBoolPtr(q.Get(c.IncludeSections))
	showDeleted := false

	page, err := h.service.GetPageByTag(tag, &title, includeSections, sortBy, sortOrder, &showDeleted)
	response.SendApiResponse(w, page, err)
}

func (h *cmsHandler) GetPageByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	title := q.Get(c.SectionTitle)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))
	includeSections := utils.ParseBoolPtr(q.Get(c.IncludeSections))

	page, err := h.service.GetPageByID(*id, &title, includeSections, sortBy, sortOrder, showDeleted)
	response.SendApiResponse(w, page, err)
}

func (h *cmsHandler) GetPagesPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	title := q.Get(c.SectionTitle)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))
	includeSections := utils.ParseBoolPtr(q.Get(c.IncludeSections))

	page, err := h.service.GetPagesPaginated(pageStr, pageSizeStr, sortBy, sortOrder, showDeleted, &title, includeSections)
	response.SendApiResponse(w, page, err)
}

func (h *cmsHandler) CreatePageSection(w http.ResponseWriter, r *http.Request) {
	var req model.CmsPageRequest
	if !utils.DecodeJSON(w, r, &req, "CreatePageSection") {
		return
	}

	if validationErrors := validatePageSectionRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	pageResponse, err := h.service.CreatePageSection(r.Context(), &req)
	response.SendApiResponse(w, pageResponse, err)
}

func (h *cmsHandler) UpdatePageSection(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	var req model.CmsPageRequest
	if !utils.DecodeJSON(w, r, &req, "UpdatePageSection") {
		return
	}

	if validationErrors := validatePageSectionRequest(&req); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	pageResponse, err := h.service.UpdatePageSection(r.Context(), *id, &req)
	response.SendApiResponse(w, pageResponse, err)
}

func (h *cmsHandler) DeletePageSection(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeletePageSection(r.Context(), *id)
	response.SendApiResponse(w, "Page deleted successfully", err)
}

func (h *cmsHandler) UndoDeletePageSection(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	err = h.service.UndoDeletePageSection(r.Context(), *id)
	response.SendApiResponse(w, "Page deleted undo successfully", err)
}

func validatePageSectionRequest(req *model.CmsPageRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Tag, "tag"),
		validator.ValidateRequired(req.Name, "name"),
	)
}
