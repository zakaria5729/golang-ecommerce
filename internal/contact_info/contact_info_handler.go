package contact_info

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/contact_info/model"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ContactInfoHandler interface {
	PostNewsLetter(w http.ResponseWriter, r *http.Request)
	CreateContactUs(w http.ResponseWriter, r *http.Request)
}

type contactInfoHandler struct {
	service ContactInfoService
}

func NewContactInfoHandler(service ContactInfoService) ContactInfoHandler {
	return &contactInfoHandler{
		service: service,
	}
}

func (h *contactInfoHandler) PostNewsLetter(w http.ResponseWriter, r *http.Request) {
	var req model.PostNewsLetterRequest
	if !utils.DecodeJSON(w, r, &req, "PostNewsLetter") {
		return
	}

	if validationErrors := validatePostNewsLetterRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.PostNewsLetter(req.Email)
	response.SendApiResponse(w, "NewsLetter added successfully", err)
}

func (h *contactInfoHandler) CreateContactUs(w http.ResponseWriter, r *http.Request) {
	var req model.CreateContactInfoRequest
	if !utils.DecodeJSON(w, r, &req, "CreateContactInfo") {
		return
	}

	if validationErrors := validateCreateContactInfoRequest(&req); len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err := h.service.CreateContactUs(&req)
	response.SendApiResponse(w, "Contact info added successfully", err)
}

func validatePostNewsLetterRequest(req *model.PostNewsLetterRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateEmail(req.Email, "email"),
	)
}

func validateCreateContactInfoRequest(req *model.CreateContactInfoRequest) validator.ValidationErrors {
	return validator.MergeValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMinLength(req.Name, "name", 2),
		validator.ValidateRequired(req.Email, "email"),
		validator.ValidateEmail(req.Email, "email"),
		validator.ValidateRequired(req.Message, "message"),
		validator.ValidateMinLength(req.Message, "message", 2),
	)
}
