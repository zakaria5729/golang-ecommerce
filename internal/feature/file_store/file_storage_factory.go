package file_storage

import (
	"context"
	"fmt"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
)

type ObjectStorage interface {
	Upload(ctx context.Context, req StorageUploadRequest) (*StorageUploadResponse, error)
	Delete(ctx context.Context, req StorageDeleteRequest) error
	GetURL(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

func NewObjectStorage(provider string, cfg config.ObjectStoreConfig) (ObjectStorage, error) {
	switch provider {
	case string(constants.ObjStoreProviderR2):
		return NewCloudflareR2Client(cfg)
	default:
		return nil, fmt.Errorf("unsupported object storage provider: %s. Currently only R2 is supported", provider)
	}
}
