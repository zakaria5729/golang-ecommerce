package file_storage

import (
	"mime/multipart"
)

type StorageDeleteRequest struct {
	PathKey string `json:"path_key"`
}

type StorageUploadRequest struct {
	File        *multipart.FileHeader `json:"file"`
	Folder      string                `json:"folder"`
	ContentType string                `json:"content_type"`
}
