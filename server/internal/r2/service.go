package r2

import (
	"context"
	"fmt"
	"server/database"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)


type Service interface {
	GetAllUploadLog(ctx context.Context, limit, offset int) ([]database.UploadLog, int64, error)
	GeneratePresignedUploadURL(ctx context.Context, objectKey string, expiresIn time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, objectKey string, expiresIn time.Duration) (string, error)
}

type service struct {
	repository      Repository
	r2Client        *s3.Client
	r2PresignClient *s3.PresignClient
	bucketName      string
}

func NewService(repository Repository, r2Client *s3.Client, r2PresignClient *s3.PresignClient, bucketName string) Service {
	return &service{
		repository:      repository,
		r2Client:        r2Client,
		r2PresignClient: r2PresignClient,
		bucketName:      bucketName,
	}
}

func (s *service) GetAllUploadLog(ctx context.Context, limit, offset int) ([]database.UploadLog, int64, error) {

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	logs, totalCount, err := s.repository.GetAllUploadLog(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}

func (s *service) GeneratePresignedUploadURL(ctx context.Context, objectKey string, expiresIn time.Duration) (string, error) {
	
	if objectKey == "" {
		return "", fmt.Errorf("object key cannot be empty")
	}

	if expiresIn <= 0 {
		expiresIn = 5 * time.Minute 
	}

	contentType := "application/x-tar" // .tar only

	putObjectInput := &s3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &objectKey,
		ContentType: &contentType,
	}

	// Generate presigned URL
	presignResult, err := s.r2PresignClient.PresignPutObject(ctx, putObjectInput, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	

	return presignResult.URL, nil
}

func (s *service) GeneratePresignedDownloadURL(ctx context.Context, objectKey string, expiresIn time.Duration) (string, error) {
	
	if objectKey == "" {
		return "", fmt.Errorf("object key cannot be empty")
	}

	if expiresIn <= 0 {
		expiresIn = 15 * time.Minute
	}

	// สร้าง GetObject request
	getObjectInput := &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &objectKey,
	}

	// Generate presigned URL
	presignResult, err := s.r2PresignClient.PresignGetObject(ctx, getObjectInput, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return presignResult.URL, nil
}
