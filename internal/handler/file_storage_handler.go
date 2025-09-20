package handler

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/file_storage"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
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

	var uploadReq file_storage.FileUploadAPIRequest
	if !utils.DecodeJSON(w, r, &uploadReq, "UploadFile") {
		return
	}

	result, err := h.usecase.UploadFile(r.Context(), *userID, uploadReq)
	if err != nil {
		logger.Logger.Error("Failed to upload file", "error", err, "user_id", *userID)
		response.SendErrorJSON(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	responseData := file_storage.FileUploadAPIResponse{
		PathKey: result.PathKey,
	}

	response.SendSuccessJSON(w, responseData)
}
