package object_storage

import (
	"context"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
)

type ObjectStorageRepository struct {
	client ObjectStorage
}

func NewObjectStorageRepository() *ObjectStorageRepository {
	cfg := config.Load()

	client, err := NewObjectStorage(constants.ObjStoreProviderR2, cfg.ObjStore)
	if err != nil {
		panic("Failed to initialize object storage: " + err.Error())
	}

	return &ObjectStorageRepository{
		client: client,
	}
}

func (r *ObjectStorageRepository) Upload(ctx context.Context, req StorageUploadRequest) (*StorageUploadResponse, error) {
	return r.client.Upload(ctx, req)
}

func (r *ObjectStorageRepository) Delete(ctx context.Context, req StorageDeleteRequest) error {
	return r.client.Delete(ctx, req)
}

func (r *ObjectStorageRepository) GetURL(ctx context.Context, key string) (string, error) {
	return r.client.GetURL(ctx, key)
}

func (r *ObjectStorageRepository) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
