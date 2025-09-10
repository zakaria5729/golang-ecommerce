package file_storage

import (
	"context"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
)

type FileStoreRepository struct {
	client ObjectStorage
}

func NewFileStoreRepository() *FileStoreRepository {
	cfg := config.Load()

	client, err := NewObjectStorage(constants.ObjStoreProviderR2, cfg.ObjStore)
	if err != nil {
		panic("Failed to initialize object storage: " + err.Error())
	}

	return &FileStoreRepository{
		client: client,
	}
}

func (r *FileStoreRepository) Upload(ctx context.Context, req StorageUploadRequest) (*StorageUploadResponse, error) {
	return r.client.Upload(ctx, req)
}

func (r *FileStoreRepository) Delete(ctx context.Context, req StorageDeleteRequest) error {
	return r.client.Delete(ctx, req)
}

func (r *FileStoreRepository) GetURL(ctx context.Context, key string) (string, error) {
	return r.client.GetURL(ctx, key)
}

func (r *FileStoreRepository) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
