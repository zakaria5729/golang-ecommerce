package file_storage

import (
	"context"

	m "github.com/easy-comerce/backend/internal/file_storage/model"
)

type FileStorageRepository interface {
	Upload(ctx context.Context, req *m.StorageUploadRequest) (*m.StorageUploadResponse, error)
	UploadRaw(ctx context.Context, req *m.StorageUploadRawRequest) (*m.StorageUploadResponse, error)
	Delete(ctx context.Context, req *m.StorageDeleteRequest) error
}

type fileStorageRepository struct {
	client ObjectStorage
}

func NewFileStorageRepository(client ObjectStorage) FileStorageRepository {
	return &fileStorageRepository{
		client: client,
	}
}

func (r *fileStorageRepository) Upload(ctx context.Context, req *m.StorageUploadRequest) (*m.StorageUploadResponse, error) {
	return r.client.Upload(ctx, req)
}

func (r *fileStorageRepository) UploadRaw(ctx context.Context, req *m.StorageUploadRawRequest) (*m.StorageUploadResponse, error) {
	return r.client.UploadRaw(ctx, req)
}

func (r *fileStorageRepository) Delete(ctx context.Context, req *m.StorageDeleteRequest) error {
	return r.client.Delete(ctx, req)
}

