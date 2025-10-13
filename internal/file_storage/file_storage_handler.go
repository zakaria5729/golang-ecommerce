package file_storage

import (
	"net/http"

	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	"github.com/easy-comerce/backend/pkg/response"
)

type FileStorageHandler struct {
	service *FileStorageService
}

func NewFileStorageHandler(service *FileStorageService) *FileStorageHandler {
	return &FileStorageHandler{
		service: service,
	}
}

func (h *FileStorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	file, fileHeader, err := r.FormFile(c.File)
	if err != nil {
		response.SendErrorJSON(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if file == nil {
		response.SendErrorJSON(w, "No file provided", http.StatusBadRequest)
		return
	}

	uploadReq := StorageUploadRequest{
		File:        fileHeader,
		Folder:      r.FormValue(c.Folder),
		ContentType: fileHeader.Header.Get(c.ContentType),
	}

	result, err := h.service.UploadFile(r.Context(), *userID, &uploadReq)
	response.SendResponse(w, result, err, http.StatusInternalServerError)
}
