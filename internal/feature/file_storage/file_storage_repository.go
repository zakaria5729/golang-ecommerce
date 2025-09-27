package file_storage

import (
	"context"

	"github.com/easy-comerce/backend/pkg/constants"
)

type FileStorageRepository struct {
	client ObjectStorage
}

func NewFileStorageRepository() *FileStorageRepository {
	client, err := NewObjectStorage(constants.ObjStoreProviderR2)
	if err != nil {
		panic("Failed to initialize object storage: " + err.Error())
	}

	return &FileStorageRepository{
		client: client,
	}
}

func (r *FileStorageRepository) Upload(ctx context.Context, req StorageUploadRequest) (*StorageUploadResponse, error) {
	return r.client.Upload(ctx, req)
}

func (r *FileStorageRepository) Delete(ctx context.Context, req StorageDeleteRequest) error {
	return r.client.Delete(ctx, req)
}

func (r *FileStorageRepository) GetURL(ctx context.Context, key string) (string, error) {
	return r.client.GetURL(ctx, key)
}

func (r *FileStorageRepository) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
