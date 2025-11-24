package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	m "github.com/easy-comerce/backend/internal/file_storage/model"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/tokenutil"
)

type CloudflareR2Provider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
}

func NewCloudflareR2Provider() (*CloudflareR2Provider, error) {
	cfg := config.GetConfig().ObjStoreConfig
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.AccessKeySecret,
			"",
		)),
	)

	if err != nil {
		l.Error("❌ Failed to create AWS config", "error", err)
		return nil, fmt.Errorf("failed to create Cloudflare R2 config: %w", err)
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = &endpoint
	})

	return &CloudflareR2Provider{
		client:        client,
		bucketName:    cfg.BucketName,
		presignClient: s3.NewPresignClient(client),
	}, nil
}

func (r *CloudflareR2Provider) Upload(ctx context.Context, req *m.StorageUploadRequest) (*m.StorageUploadResponse, error) {
	file, err := req.File.Open()
	if err != nil {
		l.Error("❌ Failed to open file", "error", err, "method", "Upload")
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		l.Error("❌ Failed to read file content", "error", err, "method", "Upload")
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return uploadFile(ctx, r, req.Folder, req.File.Filename, req.ContentType, fileBytes)
}

func (r *CloudflareR2Provider) UploadRaw(ctx context.Context, req *m.StorageUploadRawRequest) (*m.StorageUploadResponse, error) {
	return uploadFile(ctx, r, req.Folder, req.FileName, req.ContentType, req.FileData)
}

func (r *CloudflareR2Provider) Delete(ctx context.Context, req *m.StorageDeleteRequest) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &r.bucketName,
		Key:    &req.PathKey,
	})
	if err != nil {
		l.Error("❌ Failed to delete file from Cloudflare R2", "error", err, "key", req.PathKey, "method", "Delete")
		return fmt.Errorf("failed to delete file from Cloudflare R2: %w", err)
	}

	return nil
}

func (r *CloudflareR2Provider) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &r.bucketName,
		Key:    &key,
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}
		l.Error("❌ Failed to check if file exists in Cloudflare R2", "error", err, "key", key, "method", "Exists")
		return false, fmt.Errorf("failed to check if file exists in Cloudflare R2: %w", err)
	}

	return true, nil
}

func generatePresignedUploadURL(ctx context.Context, r2 *CloudflareR2Provider, pathKey string, contentType string, expiresIn time.Duration) (string, error) {
	if r2.presignClient == nil {
		l.Error("❌ presignClient is nil", "method", "generatePresignedUploadURL")
		return "", fmt.Errorf("presignClient is nil")
	}

	input := &s3.PutObjectInput{
		Bucket:      &r2.bucketName,
		Key:         &pathKey,
		ContentType: &contentType,
	}

	request, err := r2.presignClient.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	if request == nil || request.URL == "" {
		return "", fmt.Errorf("generated empty presigned URL")
	}

	return request.URL, nil
}

func uploadFile(ctx context.Context, r2 *CloudflareR2Provider, folderName string, fileName string, contentType string, fileBytes []byte) (*m.StorageUploadResponse, error) {
	fileUUID, _ := tokenutil.GenerateNewToken(true)
	if fileUUID == "" {
		l.Error("❌ Failed to generate file UUID", "method", "Upload")
		return nil, fmt.Errorf("failed to generate file UUID")
	}

	fileName = fmt.Sprintf("%s%s", fileUUID, filepath.Ext(fileName))
	pathKey := filepath.Join(folderName, fileName)
	contentLength := int64(len(fileBytes))

	_, err := r2.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &r2.bucketName,
		Key:           &pathKey,
		Body:          bytes.NewReader(fileBytes),
		ContentType:   &contentType,
		ContentLength: &contentLength,
	})

	if err == nil {
		return &m.StorageUploadResponse{
			PathKey: pathKey,
		}, nil
	}

	expiresIn := 2 * time.Minute
	presignedURL, presignErr := generatePresignedUploadURL(ctx, r2, pathKey, contentType, expiresIn)
	if presignErr != nil {
		l.Error("❌ Failed to generate presigned URL", "error", presignErr, "pathKey", pathKey, "method", "Upload")
		return nil, fmt.Errorf("direct upload failed and fallback to presigned URL failed: %w", presignErr)
	}

	httpReq, httpErr := http.NewRequest(http.MethodPut, presignedURL, bytes.NewReader(fileBytes))
	if httpErr != nil {
		return nil, fmt.Errorf("failed to create HTTP request for presigned URL: %w", httpErr)
	}
	httpReq.ContentLength = contentLength
	httpReq.Header.Set(c.ContentType, contentType)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	res, httpErr := httpClient.Do(httpReq)
	if httpErr != nil {
		return nil, fmt.Errorf("failed to upload via presigned URL: %w", httpErr)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("presigned URL upload failed with status %d: %s", res.StatusCode, string(body))
	}

	return &m.StorageUploadResponse{
		PathKey: pathKey,
	}, nil
}
