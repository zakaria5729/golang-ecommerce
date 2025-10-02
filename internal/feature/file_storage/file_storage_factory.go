package file_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/easy-comerce/backend/pkg/constants"
)

type ObjectStorage interface {
	Upload(ctx context.Context, req StorageUploadRequest) (*StorageUploadResponse, error)
	Delete(ctx context.Context, req StorageDeleteRequest) error
	GetURL(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiresIn time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

func NewObjectStorage(provider string) (ObjectStorage, error) {
	switch provider {
	case string(constants.ObjStoreProviderR2):
		return NewCloudflareR2Client()
	default:
		return nil, fmt.Errorf("unsupported object storage provider: %s. Currently only R2 is supported", provider)
	}
}
