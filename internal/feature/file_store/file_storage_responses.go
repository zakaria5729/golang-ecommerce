package file_storage

type StorageUploadResponse struct {
	Key      string
	URL      string
	ETag     string
	Location string
}

type FileUploadResponse struct {
	PathKey string `json:"path_key"`
}

type FileUploadAPIResponse struct {
	PathKey string `json:"path_key" example:"category/uuid_timestamp.pdf"`
}
