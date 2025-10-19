package file_storage

import (
	"context"

	m "github.com/easy-comerce/backend/internal/file_storage/model"
)

type FileStorageRepository struct {
	client ObjectStorage
}

func NewFileStorageRepository(client ObjectStorage) *FileStorageRepository {
	return &FileStorageRepository{
		client: client,
	}
}

func (r *FileStorageRepository) Upload(ctx context.Context, req *m.StorageUploadRequest) (*m.StorageUploadResponse, error) {
	return r.client.Upload(ctx, req)
}

func (r *FileStorageRepository) UploadRaw(ctx context.Context, req *m.StorageUploadRawRequest) (*m.StorageUploadResponse, error) {
	return r.client.UploadRaw(ctx, req)
}

func (r *FileStorageRepository) Delete(ctx context.Context, req *m.StorageDeleteRequest) error {
	return r.client.Delete(ctx, req)
}

func (r *FileStorageRepository) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
