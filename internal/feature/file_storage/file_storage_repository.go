package file_storage

import (
	"context"
)

type FileStorageRepository struct {
	client ObjectStorage
}

func NewFileStorageRepository(client ObjectStorage) *FileStorageRepository {
	return &FileStorageRepository{
		client: client,
	}
}

func (r *FileStorageRepository) Upload(ctx context.Context, req *StorageUploadRequest) (*StorageUploadResponse, error) {
	return r.client.Upload(ctx, req)
}

func (r *FileStorageRepository) Delete(ctx context.Context, req *StorageDeleteRequest) error {
	return r.client.Delete(ctx, req)
}

func (r *FileStorageRepository) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
