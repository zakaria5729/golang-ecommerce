package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/file_storage"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
)

type FileStorageHandler struct {
	service *file_storage.FileStorageService
}

func NewFileStorageHandler(service *file_storage.FileStorageService) *FileStorageHandler {
	return &FileStorageHandler{
		service: service,
	}
}

func (h *FileStorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	file, fileHeader, err := r.FormFile(c.File)
	if err != nil {
		logger.Logger.Error("Failed to get file from form", "error", err)
		response.SendErrorJSON(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if file == nil {
		response.SendErrorJSON(w, "No file provided", http.StatusBadRequest)
		return
	}

	uploadReq := file_storage.StorageUploadRequest{
		File:        fileHeader,
		Folder:      r.FormValue(c.Folder),
		ContentType: fileHeader.Header.Get(c.ContentType),
	}

	result, err := h.service.UploadFile(r.Context(), *userID, &uploadReq)
	response.SendResponse(w, result, err, http.StatusInternalServerError)
}
