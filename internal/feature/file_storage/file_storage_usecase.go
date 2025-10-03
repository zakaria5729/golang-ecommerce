package file_storage

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
)

type presignedURLCache struct {
	url       string
	expiresAt time.Time
}

type FileStorageUseCase struct {
	repo         *FileStorageRepository
	urlCache     map[string]presignedURLCache
	urlCacheLock sync.RWMutex
}

func NewFileStorageUseCase() *FileStorageUseCase {
	return &FileStorageUseCase{
		repo:     NewFileStorageRepository(),
		urlCache: make(map[string]presignedURLCache),
	}
}

func (uc *FileStorageUseCase) UploadFile(ctx context.Context, userID uint, req *StorageUploadRequest) (*StorageUploadResponse, error) {
	maxFileSize := c.SizeInMB * 3
	if req.File.Size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds, maximum allowed size of %d MB", maxFileSize/c.SizeInMB)
	}

	if req.Folder == "" || !uc.isValidFolder(req.Folder) {
		logger.Logger.Error("Invalid folder given", "method", "UploadFile", "folder", req.Folder)
		return nil, fmt.Errorf("invalid folder given")
	}

	allowedTypes := c.AllowedImageTypes + "," + c.AllowedDocTypes
	if !uc.isValidFileType(req.ContentType, allowedTypes) {
		return nil, fmt.Errorf("invalid file type: %s. Allowed types: %s", req.ContentType, allowedTypes)
	}

	uploadRes, err := uc.repo.Upload(&ctx, req)
	if err != nil {
		logger.Logger.Error("Failed to upload file", "error", err)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadRes, nil
}

func (uc *FileStorageUseCase) DeleteFile(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("file key is required")
	}

	err := uc.repo.Delete(&ctx, &StorageDeleteRequest{PathKey: key})
	if err != nil {
		logger.Logger.Error("Failed to delete file", "error", err, "key", key)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (uc *FileStorageUseCase) isValidFileType(contentType string, allowedTypes string) bool {
	types := strings.Split(allowedTypes, ",")
	for _, allowedType := range types {
		if strings.TrimSpace(allowedType) == contentType {
			return true
		}
	}
	return false
}

func (uc *FileStorageUseCase) isValidFolder(folder string) bool {
	validFolders := []string{
		c.FolderCategory,
		c.FolderProduct,
		c.FolderUser,
	}
	return slices.Contains(validFolders, folder)
}
