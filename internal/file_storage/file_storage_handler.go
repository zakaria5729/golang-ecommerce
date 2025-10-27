package file_storage

import (
	"net/http"

	m "github.com/easy-comerce/backend/internal/file_storage/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	"github.com/easy-comerce/backend/pkg/response"
)

type FileStorageHandler interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}

type fileStorageHandler struct {
	service FileStorageService
}

func NewFileStorageHandler(service FileStorageService) FileStorageHandler {
	return &fileStorageHandler{
		service: service,
	}
}

func (h *fileStorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
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

	uploadReq := m.StorageUploadRequest{
		File:        fileHeader,
		Folder:      r.FormValue(c.Folder),
		ContentType: fileHeader.Header.Get(c.ContentType),
	}

	result, err := h.service.UploadFile(r.Context(), *userID, &uploadReq)
	response.SendResponse(w, result, err, http.StatusInternalServerError)
}
