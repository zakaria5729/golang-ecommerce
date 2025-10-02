package file_storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
)

type CloudflareR2Client2 struct {
	client        *s3.Client
	uploader      *manager.Uploader
	presignClient *s3.PresignClient // Add presign client
	bucketName    string
	accountID     string
	publicDomain  string
}

func NewCloudflareR2Client2() (*CloudflareR2Client2, error) {
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
		return nil, fmt.Errorf("failed to create Cloudflare R2 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})
	uploader := manager.NewUploader(client)
	presignClient := s3.NewPresignClient(client) // Initialize presign client

	return &CloudflareR2Client2{
		client:        client,
		uploader:      uploader,
		presignClient: presignClient,
		bucketName:    cfg.BucketName,
		accountID:     cfg.AccountID,
		publicDomain:  cfg.PublicDomain,
	}, nil
}

// GeneratePresignedUploadURL generates a presigned URL for file upload
func (r *CloudflareR2Client2) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (*PresignedUploadResponse, error) {
	fullKey := r.buildKey("", key)

	logger.Logger.Info("Generating presigned upload URL",
		"key", fullKey,
		"contentType", contentType,
		"expiresIn", expiresIn.String())

	// Generate presigned PUT URL
	presignRequest, err := r.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(fullKey),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		logger.Logger.Error("Failed to generate presigned upload URL",
			"error", err,
			"key", fullKey)
		return nil, fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	logger.Logger.Info("Successfully generated presigned upload URL",
		"key", fullKey,
		"url", presignRequest.URL)

	return &PresignedUploadResponse{
		URL:       presignRequest.URL,
		Key:       fullKey,
		ExpiresAt: time.Now().Add(expiresIn),
		Method:    "PUT",
		Headers: map[string]string{
			"Content-Type": contentType,
		},
	}, nil
}

// GeneratePresignedUploadURLWithFolder generates a presigned URL for file upload with folder
func (r *CloudflareR2Client2) GeneratePresignedUploadURLWithFolder(ctx context.Context, folder, key, contentType string, expiresIn time.Duration) (*PresignedUploadResponse, error) {
	fullKey := r.buildKey(folder, key)

	logger.Logger.Info("Generating presigned upload URL with folder",
		"folder", folder,
		"key", key,
		"fullKey", fullKey,
		"contentType", contentType)

	presignRequest, err := r.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(fullKey),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return &PresignedUploadResponse{
		URL:       presignRequest.URL,
		Key:       fullKey,
		ExpiresAt: time.Now().Add(expiresIn),
		Method:    "PUT",
		Headers: map[string]string{
			"Content-Type": contentType,
		},
	}, nil
}

// GeneratePresignedDownloadURL generates a presigned URL for file download
func (r *CloudflareR2Client2) GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (*PresignedDownloadResponse, error) {
	logger.Logger.Info("Generating presigned download URL",
		"key", key,
		"expiresIn", expiresIn.String())

	presignRequest, err := r.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		logger.Logger.Error("Failed to generate presigned download URL",
			"error", err,
			"key", key)
		return nil, fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return &PresignedDownloadResponse{
		URL:       presignRequest.URL,
		Key:       key,
		ExpiresAt: time.Now().Add(expiresIn),
	}, nil
}

// GenerateMultipartUploadURL generates presigned URLs for multipart upload
func (r *CloudflareR2Client2) GenerateMultipartUploadURL(ctx context.Context, key string, contentType string, partNumber int32, uploadID string, expiresIn time.Duration) (*PresignedUploadResponse, error) {
	fullKey := r.buildKey("", key)

	presignRequest, err := r.presignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
		Bucket:     aws.String(r.bucketName),
		Key:        aws.String(fullKey),
		PartNumber: aws.Int32(partNumber),
		UploadId:   aws.String(uploadID),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate multipart upload URL: %w", err)
	}

	return &PresignedUploadResponse{
		URL:        presignRequest.URL,
		Key:        fullKey,
		ExpiresAt:  time.Now().Add(expiresIn),
		Method:     "PUT",
		PartNumber: partNumber,
		UploadID:   uploadID,
	}, nil
}

// Your existing methods remain unchanged...
func (r *CloudflareR2Client2) Upload(ctx context.Context, req StorageUploadRequest) (*StorageUploadResponse, error) {
	key := r.buildKey(req.Folder, req.Key)

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		Body:        req.Body,
		ContentType: aws.String(req.ContentType),
	}

	if req.Metadata != nil {
		uploadInput.Metadata = req.Metadata
	}

	result, err := r.uploader.Upload(ctx, uploadInput)
	if err != nil {
		logger.Logger.Error("Failed to upload file to Cloudflare R2", "error", err, "key", key)
		return nil, fmt.Errorf("failed to upload file to Cloudflare R2: %w", err)
	}

	etag := ""
	if result.ETag != nil {
		etag = strings.Trim(*result.ETag, "\"")
	}

	return &StorageUploadResponse{
		Key:      key,
		ETag:     etag,
		Location: result.Location,
	}, nil
}

func (r *CloudflareR2Client2) Delete(ctx context.Context, req StorageDeleteRequest) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(req.Key),
	})
	if err != nil {
		logger.Logger.Error("Failed to delete file from Cloudflare R2", "error", err, "key", req.Key)
		return fmt.Errorf("failed to delete file from Cloudflare R2: %w", err)
	}

	return nil
}

func (r *CloudflareR2Client2) GetURL(ctx context.Context, key string) (string, error) {
	publicURL := r.buildPublicURL(key)
	return publicURL, nil
}

func (r *CloudflareR2Client2) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
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

func (r *CloudflareR2Client2) buildKey(folder, key string) string {
	if folder == "" {
		return key
	}
	return folder + "/" + key
}

func (r *CloudflareR2Client2) buildPublicURL(key string) string {
	if r.publicDomain != "" {
		return fmt.Sprintf("https://%s/%s", r.publicDomain, key)
	}
	return fmt.Sprintf("https://%s.%s.r2.dev/%s", r.bucketName, r.accountID, key)
}

// Response types for presigned URLs
type PresignedUploadResponse struct {
	URL        string            `json:"url"`
	Key        string            `json:"key"`
	ExpiresAt  time.Time         `json:"expires_at"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers,omitempty"`
	PartNumber int32             `json:"part_number,omitempty"`
	UploadID   string            `json:"upload_id,omitempty"`
}

type PresignedDownloadResponse struct {
	URL       string    `json:"url"`
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}
