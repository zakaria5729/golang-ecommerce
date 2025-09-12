package file_storage

import (
	"io"
	"mime/multipart"
)

type StorageUploadRequest struct {
	Body        io.Reader
	Metadata    map[string]string
	Key         string
	ContentType string
	Folder      string
}

type StorageDeleteRequest struct {
	Key string
}

type FileUploadAPIRequest struct {
	File   *multipart.FileHeader `json:"file"`
	Folder string                `json:"folder"`
	UserID string                `json:"user_id"`
}
