package object_storage

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/google/uuid"
)

type ObjectStorageUseCase struct {
	repo *ObjectStorageRepository
}

func NewObjectStorageUseCase() *ObjectStorageUseCase {
	return &ObjectStorageUseCase{
		repo: NewObjectStorageRepository(),
	}
}

func (uc *ObjectStorageUseCase) UploadFile(ctx context.Context, req FileUploadAPIRequest) (*FileUploadResponse, error) {
	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	if req.File.Size > constants.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", constants.MaxFileSize)
	}

	contentType := req.File.Header.Get("Content-Type")
	allowedTypes := constants.AllowedImageTypes + "," + constants.AllowedDocTypes
	if !uc.isValidFileType(contentType, allowedTypes) {
		return nil, fmt.Errorf("invalid file type: %s. Allowed types: %s", contentType, allowedTypes)
	}

	file, err := req.File.Open()
	if err != nil {
		logger.Logger.Error("Failed to open uploaded file", "error", err)
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	ext := filepath.Ext(req.File.Filename)
	sanitizedOriginalName := strings.ReplaceAll(req.File.Filename, " ", "_")
	sanitizedOriginalName = strings.ReplaceAll(sanitizedOriginalName, "(", "")
	sanitizedOriginalName = strings.ReplaceAll(sanitizedOriginalName, ")", "")
	sanitizedOriginalName = strings.ReplaceAll(sanitizedOriginalName, "~", "")
	sanitizedOriginalName = strings.ReplaceAll(sanitizedOriginalName, "-", "_")

	filename := fmt.Sprintf("%s_%d_%s%s", uuid.New().String(), time.Now().Unix(), sanitizedOriginalName, ext)
	uploadReq := StorageUploadRequest{
		Key:         filename,
		Body:        file,
		ContentType: contentType,
		Folder:      req.Folder,
		Metadata: map[string]string{
			"original_filename": req.File.Filename,
			"user_id":           req.UserID,
			"uploaded_at":       time.Now().Format(time.RFC3339),
		},
	}

	result, err := uc.repo.Upload(ctx, uploadReq)
	if err != nil {
		logger.Logger.Error("Failed to upload file", "error", err, "filename", req.File.Filename)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &FileUploadResponse{
		Key:      result.Key,
		URL:      result.URL,
		Filename: req.File.Filename,
		Size:     req.File.Size,
	}, nil
}

func (uc *ObjectStorageUseCase) DeleteFile(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("file key is required")
	}

	err := uc.repo.Delete(ctx, StorageDeleteRequest{Key: key})
	if err != nil {
		logger.Logger.Error("Failed to delete file", "error", err, "key", key)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (uc *ObjectStorageUseCase) GetFileURL(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("file key is required")
	}

	url, err := uc.repo.GetURL(ctx, key)
	if err != nil {
		logger.Logger.Error("Failed to get file URL", "error", err, "key", key)
		return "", fmt.Errorf("failed to get file URL: %w", err)
	}

	return url, nil
}

func (uc *ObjectStorageUseCase) isValidFileType(contentType string, allowedTypes string) bool {
	types := strings.SplitSeq(allowedTypes, ",")
	for allowedType := range types {
		if strings.TrimSpace(allowedType) == contentType {
			return true
		}
	}
	return false
}
