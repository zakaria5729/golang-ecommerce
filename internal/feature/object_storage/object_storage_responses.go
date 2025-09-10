package object_storage

type StorageUploadResponse struct {
	Key      string
	URL      string
	ETag     string
	Location string
}

type FileUploadResponse struct {
	Key      string `json:"key"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type FileUploadAPIResponse struct {
	Key      string `json:"key" example:"category/uuid_timestamp.pdf"`
	URL      string `json:"url" example:"https://bucket.r2.cloudflarestorage.com/category/uuid_timestamp.pdf"`
	Filename string `json:"filename" example:"document.pdf"`
	Size     int64  `json:"size" example:"1024000"`
}
