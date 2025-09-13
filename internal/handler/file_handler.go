package handler

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	file_storage "github.com/easy-comerce/backend/internal/feature/file_storage"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
)

type FileHandler struct {
	fileUseCase *file_storage.FileStoreUseCase
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileUseCase: file_storage.NewFileStoreUseCase(),
	}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	user, exists := r.Context().Value(constants.UserContextKey).(*user.User)
	if !exists {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(constants.MaxRequestSizeMB)
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

	if fileHeader.Size > constants.MaxFileSizeMB {
		response.SendErrorJSON(w, fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", constants.MaxFileSizeMB), http.StatusBadRequest)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	allowedTypes := constants.AllowedImageTypes + "," + constants.AllowedDocTypes
	if !h.isValidFileType(contentType, allowedTypes) {
		response.SendErrorJSON(w, fmt.Sprintf("Invalid file type: %s. Allowed types: %s", contentType, allowedTypes), http.StatusBadRequest)
		return
	}

	uploadReq := file_storage.FileUploadAPIRequest{
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

	responseData := file_storage.FileUploadAPIResponse{
		PathKey: result.PathKey,
	}

	response.SendSuccessJSON(w, responseData)
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
