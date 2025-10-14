package system

import (
	"fmt"
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

func (h *SystemHandler) GetSystemLogFiles(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get(c.FileName)
	logFiles, err := h.service.GetSystemLogFiles(fileName)
	response.SendResponse(w, logFiles, err, http.StatusInternalServerError)
}

func (h *SystemHandler) DownloadSystemLogFile(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue(c.FileName)
	if filename == "" {
		response.SendErrorJSON(w, "filename parameter is required", http.StatusBadRequest)
		return
	}

	content, err := h.service.DownloadSystemLogFile(filename)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(content)
}

func (h *SystemHandler) DeleteSystemLogFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue(c.FileName)
	err := h.service.DeleteSystemLogFile(fileName)
	var msg string
	if err == nil {
		msg = "Log file deleted successfully"
	}
	response.SendResponse(w, msg, err, http.StatusInternalServerError)
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
