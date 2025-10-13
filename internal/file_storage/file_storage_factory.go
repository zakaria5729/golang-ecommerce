package file_storage

import (
	"context"
	"fmt"

	c "github.com/easy-comerce/backend/pkg/constants"
)

type ObjectStorage interface {
	Upload(ctx context.Context, req *StorageUploadRequest) (*StorageUploadResponse, error)
	UploadRaw(ctx context.Context, req *StorageUploadRawRequest) (*StorageUploadResponse, error)
	Delete(ctx context.Context, req *StorageDeleteRequest) error
	Exists(ctx context.Context, key string) (bool, error)
}

func NewObjectStorage() (ObjectStorage, error) {
	var err error
	var storage ObjectStorage

	switch c.EnvActiveObjectStorage {
	case c.ObjStoreProviderR2:
		storage, err = NewCloudflareR2Client()
	default:
		err = fmt.Errorf("unsupported object storage provider: %s", c.EnvActiveObjectStorage)
	}

	return storage, err
}
