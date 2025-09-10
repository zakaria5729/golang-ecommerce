package object_storage

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
)

type CloudflareR2Client struct {
	client       *s3.Client
	uploader     *manager.Uploader
	bucketName   string
	accountID    string
	publicDomain string
}

func NewCloudflareR2Client(cfg config.ObjectStoreConfig) (*CloudflareR2Client, error) {
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

	return &CloudflareR2Client{
		client:       client,
		uploader:     uploader,
		bucketName:   cfg.BucketName,
		accountID:    cfg.AccountID,
		publicDomain: cfg.PublicDomain,
	}, nil
}

func (r *CloudflareR2Client) Upload(ctx context.Context, req StorageUploadRequest) (*StorageUploadResponse, error) {
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

	logger.Logger.Info("File uploaded to Cloudflare R2 successfully", "key", key, "location", result.Location)

	etag := ""
	if result.ETag != nil {
		etag = strings.Trim(*result.ETag, "\"")
	}

	// Construct public URL for Cloudflare R2
	publicURL := r.buildPublicURL(key)

	return &StorageUploadResponse{
		Key:      key,
		URL:      publicURL,
		ETag:     etag,
		Location: result.Location,
	}, nil
}

func (r *CloudflareR2Client) Delete(ctx context.Context, req StorageDeleteRequest) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(req.Key),
	})
	if err != nil {
		logger.Logger.Error("Failed to delete file from Cloudflare R2", "error", err, "key", req.Key)
		return fmt.Errorf("failed to delete file from Cloudflare R2: %w", err)
	}

	logger.Logger.Info("File deleted from Cloudflare R2 successfully", "key", req.Key)
	return nil
}

func (r *CloudflareR2Client) GetURL(ctx context.Context, key string) (string, error) {
	// Return public URL for Cloudflare R2
	publicURL := r.buildPublicURL(key)
	return publicURL, nil
}

func (r *CloudflareR2Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		logger.Logger.Error("Failed to check if file exists in Cloudflare R2", "error", err, "key", key)
		return false, fmt.Errorf("failed to check if file exists in Cloudflare R2: %w", err)
	}

	return true, nil
}

func (r *CloudflareR2Client) buildKey(folder, key string) string {
	if folder == "" {
		return key
	}
	return filepath.Join(folder, key)
}

func (r *CloudflareR2Client) buildPublicURL(key string) string {
	if r.publicDomain != "" {
		return fmt.Sprintf("https://%s/%s", r.publicDomain, key)
	}
	return fmt.Sprintf("https://%s.%s.r2.dev/%s", r.bucketName, r.accountID, key)
}
