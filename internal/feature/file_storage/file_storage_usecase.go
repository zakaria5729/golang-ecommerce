package file_storage

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/google/uuid"
)

type FileStorageUseCase struct {
	repo *FileStorageRepository
}

func NewFileStorageUseCase() *FileStorageUseCase {
	return &FileStorageUseCase{
		repo: NewFileStorageRepository(),
	}
}

func (uc *FileStorageUseCase) UploadFile(ctx context.Context, userID uint, req FileUploadAPIRequest) (*FileUploadResponse, error) {
	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	if req.Folder == "" || !uc.isValidFolder(req.Folder) {
		logger.Logger.Error("Invalid folder given", "method", "UploadFile", "folder", req.Folder)
		return nil, fmt.Errorf("invalid folder given")
	}

	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	if req.File.Size > constants.MaxFileSizeMB {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", constants.MaxFileSizeMB)
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
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), timeutil.NowUTC().Unix(), ext)
	uploadReq := StorageUploadRequest{
		Key:         filename,
		Body:        file,
		ContentType: contentType,
		Folder:      req.Folder,
		Metadata: map[string]string{
			"original_filename": req.File.Filename,
			"user_id":           fmt.Sprintf("%d", userID),
			"uploaded_at":       timeutil.NowUTC().Format(time.RFC3339),
		},
	}

	result, err := uc.repo.Upload(ctx, uploadReq)
	if err != nil {
		logger.Logger.Error("Failed to upload file", "error", err, "filename", req.File.Filename)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &FileUploadResponse{
		PathKey: result.Key,
	}, nil
}

func (uc *FileStorageUseCase) DeleteFile(ctx context.Context, key string) error {
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

func (uc *FileStorageUseCase) GetFileURL(ctx context.Context, key string) (string, error) {
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

// GeneratePresignedUploadURL generates a presigned URL for file upload
func (uc *FileStorageUseCase) GeneratePresignedUploadURL(ctx context.Context, userID uint, folder, fileName, contentType string, expiresIn time.Duration) (string, error) {
	if folder == "" || !uc.isValidFolder(folder) {
		return "", fmt.Errorf("invalid folder: %s", folder)
	}

	if fileName == "" {
		return "", fmt.Errorf("file name is required")
	}

	// Generate unique filename
	ext := filepath.Ext(fileName)
	uniqueFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), timeutil.NowUTC().Unix(), ext)
	key := folder + "/" + uniqueFileName

	// Generate presigned URL
	presignedURL, err := uc.repo.GeneratePresignedUploadURL(ctx, key, contentType, expiresIn)
	if err != nil {
		logger.Logger.Error("Failed to generate presigned URL", "error", err, "user_id", userID, "key", key)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL, nil
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
		constants.FolderCategory,
		constants.FolderProduct,
		constants.FolderUser,
	}
	return slices.Contains(validFolders, folder)
}
