package r2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Config holds configuration for connecting to Cloudflare R2
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Endpoint        string
}

// NewR2Client creates a new S3 client configured for Cloudflare R2
// Reference: https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/
func NewR2Client(cfg R2Config) (*s3.Client, error) {
	// Load default AWS SDK configuration with R2 credentials
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"", // Session token (not required for R2)
			),
		),
		// Region is required by SDK but not used by R2
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with R2 endpoint
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// Set R2 endpoint using account ID
		o.BaseEndpoint = aws.String(
			fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
		)
	})

	return client, nil
}

// NewPresignClient creates a presign client for generating presigned URLs
// Presigned URLs can be used to temporarily share public read/write access to objects
func NewPresignClient(client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(client)
}

// ValidateConfig checks if all required R2 configuration fields are set
func ValidateConfig(cfg R2Config) error {
	if cfg.AccountID == "" {
		return fmt.Errorf("R2 Account ID is required")
	}
	if cfg.AccessKeyID == "" {
		return fmt.Errorf("R2 Access Key ID is required")
	}
	if cfg.SecretAccessKey == "" {
		return fmt.Errorf("R2 Secret Access Key is required")
	}
	if cfg.BucketName == "" {
		return fmt.Errorf("R2 Bucket Name is required")
	}
	return nil
}
