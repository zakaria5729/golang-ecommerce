package file_storage

import (
	"context"
	"fmt"

	m "github.com/easy-comerce/backend/internal/file_storage/model"
	"github.com/easy-comerce/backend/internal/file_storage/provider"
	c "github.com/easy-comerce/backend/pkg/constants"
)

type ObjectStorage interface {
	Upload(ctx context.Context, req *m.StorageUploadRequest) (*m.StorageUploadResponse, error)
	UploadRaw(ctx context.Context, req *m.StorageUploadRawRequest) (*m.StorageUploadResponse, error)
	Delete(ctx context.Context, req *m.StorageDeleteRequest) error
	Exists(ctx context.Context, key string) (bool, error)
}

func NewObjectStorage() (ObjectStorage, error) {
	var err error
	var storage ObjectStorage

	switch c.EnvActiveObjectStorage {
	case c.ObjStoreProviderR2:
		storage, err = provider.NewCloudflareR2Provider()
	default:
		err = fmt.Errorf("unsupported object storage provider: %s", c.EnvActiveObjectStorage)
	}

	return storage, err
}
