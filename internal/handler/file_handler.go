package handler

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	object_storage "github.com/easy-comerce/backend/internal/feature/object_storage"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
)

type FileHandler struct {
	fileUseCase *object_storage.ObjectStorageUseCase
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileUseCase: object_storage.NewObjectStorageUseCase(),
	}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {

	user, exists := r.Context().Value(constants.UserContextKey).(*user.User)
	if !exists {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(constants.MaxRequestSize)
	if err != nil {
		response.SendErrorJSON(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	folder := r.FormValue("folder")
	if folder == "" {
		response.SendErrorJSON(w, "Folder is required", http.StatusBadRequest)
		return
	}

	if !h.isValidFolder(folder) {
		response.SendErrorJSON(w, "Invalid folder. Allowed folders: category, product, user, review, address, wishlist", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.SendErrorJSON(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if fileHeader.Size > constants.MaxFileSize {
		response.SendErrorJSON(w, fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", constants.MaxFileSize), http.StatusBadRequest)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	allowedTypes := constants.AllowedImageTypes + "," + constants.AllowedDocTypes
	if !h.isValidFileType(contentType, allowedTypes) {
		response.SendErrorJSON(w, fmt.Sprintf("Invalid file type: %s. Allowed types: %s", contentType, allowedTypes), http.StatusBadRequest)
		return
	}

	uploadReq := object_storage.FileUploadAPIRequest{
		File:   fileHeader,
		Folder: folder,
		UserID: fmt.Sprintf("%d", user.ID),
	}

	result, err := h.fileUseCase.UploadFile(r.Context(), uploadReq)
	if err != nil {
		logger.Logger.Error("Failed to upload file", "error", err, "user_id", user.ID)
		response.SendErrorJSON(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	responseData := object_storage.FileUploadAPIResponse{
		Key:      result.Key,
		URL:      result.URL,
		Filename: result.Filename,
		Size:     result.Size,
	}

	response.SendSuccessJSON(w, responseData)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {

	user, exists := r.Context().Value(constants.UserContextKey).(*user.User)
	if !exists {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	key := r.PathValue("key")
	if key == "" {
		response.SendErrorJSON(w, "File key is required", http.StatusBadRequest)
		return
	}

	err := h.fileUseCase.DeleteFile(r.Context(), key)
	if err != nil {
		logger.Logger.Error("Failed to delete file", "error", err, "user_id", user.ID, "key", key)
		response.SendErrorJSON(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{
		"key": key,
	})
}

func (h *FileHandler) GetFileURL(w http.ResponseWriter, r *http.Request) {

	key := r.PathValue("key")
	if key == "" {
		response.SendErrorJSON(w, "File key is required", http.StatusBadRequest)
		return
	}

	url, err := h.fileUseCase.GetFileURL(r.Context(), key)
	if err != nil {
		logger.Logger.Error("Failed to get file URL", "error", err, "key", key)
		response.SendErrorJSON(w, "Failed to get file URL", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{
		"url": url,
	})
}

func (h *FileHandler) isValidFolder(folder string) bool {
	validFolders := []string{
		constants.FolderCategory,
		constants.FolderProduct,
		constants.FolderUser,
		constants.FolderReview,
		constants.FolderAddress,
		constants.FolderWishlist,
	}

	return slices.Contains(validFolders, folder)
}

func (h *FileHandler) isValidFileType(contentType string, allowedTypes string) bool {
	types := strings.Split(allowedTypes, ",")
	for _, allowedType := range types {
		if strings.TrimSpace(allowedType) == contentType {
			return true
		}
	}
	return false
}
