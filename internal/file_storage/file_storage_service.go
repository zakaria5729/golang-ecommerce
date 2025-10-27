package file_storage

import (
	"context"
	"fmt"
	"slices"
	"strings"

	m "github.com/easy-comerce/backend/internal/file_storage/model"
	c "github.com/easy-comerce/backend/pkg/constants"
)

type FileStorageService interface {
	UploadFile(ctx context.Context, userID uint, req *m.StorageUploadRequest) (*m.StorageUploadResponse, error)
	DeleteFile(ctx context.Context, key string) error
}

type fileStorageService struct {
	repo FileStorageRepository
}

func NewFileStorageService(repo FileStorageRepository) FileStorageService {
	return &fileStorageService{
		repo: repo,
	}
}

func (s *fileStorageService) UploadFile(ctx context.Context, userID uint, req *m.StorageUploadRequest) (*m.StorageUploadResponse, error) {
	maxFileSize := c.SizeInMB * 3
	if req.File.Size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds, maximum allowed size of %d MB", maxFileSize/c.SizeInMB)
	}

	if req.Folder == "" || !isValidFolder(req.Folder) {
		return nil, fmt.Errorf("invalid folder given")
	}

	allowedTypes := c.AllowedImageTypes + "," + c.AllowedDocTypes
	if !isValidFileType(req.ContentType, allowedTypes) {
		return nil, fmt.Errorf("invalid file type: %s. Allowed types: %s", req.ContentType, allowedTypes)
	}

	uploadRes, err := s.repo.Upload(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadRes, nil
}

func (s *fileStorageService) DeleteFile(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("file key is required")
	}

	err := s.repo.Delete(ctx, &m.StorageDeleteRequest{PathKey: key})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func isValidFileType(contentType string, allowedTypes string) bool {
	types := strings.SplitSeq(allowedTypes, ",")
	for allowedType := range types {
		if strings.TrimSpace(allowedType) == contentType {
			return true
		}
	}
	return false
}

func isValidFolder(folder string) bool {
	validFolders := []string{
		c.FolderCategory,
		c.FolderProduct,
		c.FolderUser,
	}
	return slices.Contains(validFolders, folder)
}
