package file_storage

import (
	"context"
	"fmt"

	c "github.com/easy-comerce/backend/pkg/constants"
)

type ObjectStorage interface {
	Upload(ctx context.Context, req *StorageUploadRequest) (*StorageUploadResponse, error)
	Delete(ctx context.Context, req *StorageDeleteRequest) error
	Exists(ctx context.Context, key string) (bool, error)
}

func NewObjectStorage(provider string) (ObjectStorage, error) {
	var err error
	var storage ObjectStorage

	switch provider {
	case c.ObjStoreProviderR2:
		storage, err = NewCloudflareR2Client()
	default:
		err = fmt.Errorf("unsupported object storage provider: %s", provider)
	}

	return storage, err
}
