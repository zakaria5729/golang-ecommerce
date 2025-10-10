package file_storage

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
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/tokenutil"
)

type CloudflareR2Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
}

func NewCloudflareR2Client() (*CloudflareR2Client, error) {
	cfg := config.GetConfig().ObjStore
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.AccessKeySecret,
			"",
		)),
	)

	if err != nil {
		l.Logger.Error("Failed to create AWS config", "error", err)
		return nil, fmt.Errorf("failed to create Cloudflare R2 config: %w", err)
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = &endpoint
	})

	return &CloudflareR2Client{
		client:        client,
		bucketName:    cfg.BucketName,
		presignClient: s3.NewPresignClient(client),
	}, nil
}

func (r *CloudflareR2Client) Upload(ctx context.Context, req *StorageUploadRequest) (*StorageUploadResponse, error) {
	file, err := req.File.Open()
	if err != nil {
		l.Logger.Error("Failed to open file", "error", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		l.Logger.Error("Failed to read file content", "error", err)
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	fileUUID, _ := tokenutil.GenerateNewToken(true)
	if fileUUID == "" {
		l.Logger.Error("Failed to generate file UUID")
		return nil, fmt.Errorf("failed to generate file UUID")
	}

	fileName := fmt.Sprintf("%s%s", fileUUID, filepath.Ext(req.File.Filename))
	pathKey := filepath.Join(req.Folder, fileName)
	contentLength := int64(len(fileBytes))

	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &r.bucketName,
		Key:           &pathKey,
		Body:          bytes.NewReader(fileBytes),
		ContentType:   &req.ContentType,
		ContentLength: &contentLength,
	})

	if err == nil {
		return &StorageUploadResponse{
			PathKey: pathKey,
		}, nil
	}

	expiresIn := 2 * time.Minute
	presignedURL, presignErr := r.generatePresignedUploadURL(ctx, pathKey, req.ContentType, expiresIn)
	if presignErr != nil {
		l.Logger.Error("Failed to generate presigned URL", "error", presignErr, "pathKey", pathKey)
		return nil, fmt.Errorf("direct upload failed and fallback to presigned URL failed: %w", presignErr)
	}

	httpReq, httpErr := http.NewRequest(http.MethodPut, presignedURL, bytes.NewReader(fileBytes))
	if httpErr != nil {
		return nil, fmt.Errorf("failed to create HTTP request for presigned URL: %w", httpErr)
	}
	httpReq.ContentLength = contentLength
	httpReq.Header.Set("Content-Type", req.ContentType)

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

	return &StorageUploadResponse{
		PathKey: pathKey,
	}, nil
}

func (r *CloudflareR2Client) Delete(ctx context.Context, req *StorageDeleteRequest) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &r.bucketName,
		Key:    &req.PathKey,
	})
	if err != nil {
		logger.Logger.Error("Failed to delete file from Cloudflare R2", "error", err, "key", req.PathKey)
		return fmt.Errorf("failed to delete file from Cloudflare R2: %w", err)
	}

	return nil
}

func (r *CloudflareR2Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &r.bucketName,
		Key:    &key,
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}
		logger.Logger.Error("Failed to check if file exists in Cloudflare R2", "error", err, "key", key)
		return false, fmt.Errorf("failed to check if file exists in Cloudflare R2: %w", err)
	}

	return true, nil
}

func (r *CloudflareR2Client) generatePresignedUploadURL(ctx context.Context, pathKey string, contentType string, expiresIn time.Duration) (string, error) {
	if r.presignClient == nil {
		l.Logger.Error("presignClient is nil")
		return "", fmt.Errorf("presignClient is nil")
	}

	input := &s3.PutObjectInput{
		Bucket:      &r.bucketName,
		Key:         &pathKey,
		ContentType: &contentType,
	}

	request, err := r.presignClient.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
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
