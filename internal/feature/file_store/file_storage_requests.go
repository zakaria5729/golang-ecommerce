package file_storage

import (
	"io"
	"mime/multipart"
)

type StorageUploadRequest struct {
	Key         string
	Body        io.Reader
	ContentType string
	Folder      string
	Metadata    map[string]string
}

type StorageDeleteRequest struct {
	Key string
}

type FileUploadAPIRequest struct {
	File   *multipart.FileHeader
	Folder string
	UserID string
}
