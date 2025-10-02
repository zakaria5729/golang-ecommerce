package handler

import (
	"net/http"
	"time"

	"github.com/easy-comerce/backend/internal/feature/file_storage"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
)

type FileStorageHandler struct {
	usecase *file_storage.FileStorageUseCase
}

func NewFileStorageHandler() *FileStorageHandler {
	return &FileStorageHandler{
		usecase: file_storage.NewFileStorageUseCase(),
	}
}

func (h *FileStorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
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

	uploadReq := file_storage.FileUploadAPIRequest{
		File:   fileHeader,
		Folder: r.FormValue(c.Folder),
	}

	result, err := h.usecase.UploadFile(r.Context(), *userID, uploadReq)
	if err != nil {
		logger.Logger.Error("Failed to upload file", "error", err, "user_id", *userID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := file_storage.FileUploadAPIResponse{
		PathKey: result.PathKey,
	}

	response.SendSuccessJSON(w, responseData)
}

// GeneratePresignedUploadURL generates a presigned URL for file upload
func (h *FileStorageHandler) GeneratePresignedUploadURL(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	folder := r.URL.Query().Get("folder")
	fileName := r.URL.Query().Get("file_name")
	contentType := r.URL.Query().Get("content_type")
	expiresIn := r.URL.Query().Get("expires_in") // in minutes, default 60

	if fileName == "" {
		response.SendErrorJSON(w, "file_name is required", http.StatusBadRequest)
		return
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Parse expires_in (default to 60 minutes)
	expirationMinutes := 60
	if expiresIn != "" {
		if parsed, err := time.ParseDuration(expiresIn + "m"); err == nil {
			expirationMinutes = int(parsed.Minutes())
		}
	}

	expirationDuration := time.Duration(expirationMinutes) * time.Minute

	// Generate presigned URL
	presignedURL, err := h.usecase.GeneratePresignedUploadURL(r.Context(), *userID, folder, fileName, contentType, expirationDuration)
	if err != nil {
		logger.Logger.Error("Failed to generate presigned URL", "error", err, "user_id", *userID)
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"presigned_url": presignedURL,
		"expires_in":    expirationMinutes,
		"expires_at":    time.Now().Add(expirationDuration).Format(time.RFC3339),
		"key":           folder + "/" + fileName,
	}

	response.SendSuccessJSON(w, responseData)
}
