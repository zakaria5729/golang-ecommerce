package system

import (
	"fmt"
	"net/http"

	m "github.com/easy-comerce/backend/internal/system/model"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type SystemHandler interface {
	GetSystemHealthCheck(w http.ResponseWriter, r *http.Request)
	GetSystemDbStats(w http.ResponseWriter, r *http.Request)
	GetSystemLogFiles(w http.ResponseWriter, r *http.Request)
	DownloadSystemLogFile(w http.ResponseWriter, r *http.Request)
	DeleteSystemLogFile(w http.ResponseWriter, r *http.Request)
	HandleSocialFlowTemp(w http.ResponseWriter, r *http.Request)
	HandleSocialFlowCallbackTemp(w http.ResponseWriter, r *http.Request)
}

type systemHandler struct {
	service SystemService
}

func NewSystemHandler(service SystemService) SystemHandler {
	return &systemHandler{
		service: service,
	}
}

func (h *systemHandler) GetSystemHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req m.SystemHealthRequest
	if !utils.DecodeJSON(w, r, &req, "SystemHealthCheck") {
		return
	}

	if req.HealthToken != nil && *req.HealthToken == config.GetConfig().SecretConfig.AppHealthCheckToken {
		healthResponse := h.service.GetSystemHealthCheck(r.Context())
		response.SendApiResponse(w, healthResponse, nil)
		return
	}

	response.SendErrorJSON(w, "Invalid health token", http.StatusForbidden)
}

func (h *systemHandler) GetSystemDbStats(w http.ResponseWriter, r *http.Request) {
	dbStats, err := h.service.GetSystemDbStats(r.Context())
	response.SendApiResponse(w, dbStats, err)
}

func (h *systemHandler) GetSystemLogFiles(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get(c.FileName)
	logFiles, err := h.service.GetSystemLogFiles(fileName)
	response.SendApiResponse(w, logFiles, err)
}

func (h *systemHandler) DownloadSystemLogFile(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue(c.FileName)
	if filename == "" {
		response.SendErrorJSON(w, "Invalid filename provided", http.StatusBadRequest)
		return
	}

	content, err := h.service.DownloadSystemLogFile(filename)
	if err != nil {
		response.SendApiResponse(w, nil, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(content)
}

func (h *systemHandler) DeleteSystemLogFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue(c.FileName)
	err := h.service.DeleteSystemLogFile(fileName)
	var msg string
	if err == nil {
		msg = "Log file deleted successfully"
	}
	response.SendApiResponse(w, msg, err)
}

func (h *systemHandler) HandleSocialFlowTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		authType := r.URL.Query().Get("auth_type")

		switch authType {
		case c.AuthTypeGoogle:
			h.service.HandleGoogleLoginTemp(w, r)
		case c.AuthTypeFacebook:
			h.service.HandleFacebookLoginTemp(w, r)
		}
	}
}

func (h *systemHandler) HandleSocialFlowCallbackTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		authType := r.URL.Query().Get("auth_type")

		switch authType {
		case c.AuthTypeGoogle:
			h.service.HandleGoogleLoginCallbackTemp(w, r)
		case c.AuthTypeFacebook:
			h.service.HandleFacebookLoginCallbackTemp(w, r)
		}
	}
}
