package system

import (
	"net/http"

	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SystemHandler struct {
	service *SystemService
}

func NewSystemHandler(service *SystemService) *SystemHandler {
	return &SystemHandler{
		service: service,
	}
}

func (h *SystemHandler) SystemHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req SystemHealthRequest
	if !utils.DecodeJSON(w, r, &req, "SystemHealthCheck") {
		return
	}

	if req.HealthToken != nil && *req.HealthToken == c.AppHealthCheckToken {
		healthResponse := h.service.SystemHealthCheck(r.Context())
		response.SendResponse(w, healthResponse, nil, http.StatusOK)
		return
	}

	response.SendErrorJSON(w, "Invalid health token", http.StatusForbidden)
}

func (h *SystemHandler) HandleSocialFlowTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		authType := r.URL.Query().Get("auth_type")

		switch authType {
		case c.AuthTypeGoogle:
			h.service.HandleGoogleLoginTemp(w, r)
		case c.AuthTypeFacebook:
			// h.service.HandleFacebookLoginTemp(w, r)
		}
	}
}

func (h *SystemHandler) HandleSocialFlowCallbackTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		authType := r.URL.Query().Get("auth_type")

		switch authType {
		case c.AuthTypeGoogle:
			h.service.HandleGoogleLoginCallbackTemp(w, r)
		case c.AuthTypeFacebook:
			// h.service.HandleFacebookLoginTemp(w, r)
		}
	}
}
