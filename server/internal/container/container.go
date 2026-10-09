package container

import (
	"log"
	"server/config"
	"server/internal/r2"
	"server/internal/webhook"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB

	// R2 Module
	R2Client        *s3.Client
	R2PresignClient *s3.PresignClient
	R2Repository    r2.Repository
	R2Service       r2.Service
	R2Handler       *r2.Handler

	// Webhook Module
	WebhookService webhook.Service
	WebhookHandler *webhook.Handler
}


func NewContainer(db *gorm.DB, cfg *config.Config) *Container {

	// Initialize R2 Client
	r2Cfg := r2.R2Config{
		AccountID:       cfg.R2.AccountID,
		AccessKeyID:     cfg.R2.AccessKeyID,
		SecretAccessKey: cfg.R2.SecretAccessKey,
		BucketName:      cfg.R2.BucketName,
		Endpoint:        cfg.R2.Endpoint,
	}

	// Validate R2 config
	if err := r2.ValidateConfig(r2Cfg); err != nil {
		log.Fatalf("❌R2 configuration validation failed: %v", err)
	}

	// Create R2 client
	r2Client, err := r2.NewR2Client(r2Cfg)
	if err != nil {
		log.Fatalf("❌Failed to create R2 client: %v", err)
	}

	// Create Presign client
	r2PresignClient := r2.NewPresignClient(r2Client)

	log.Println("R2 client initialized successfully")

	// Initialize R2 module layers
	r2Repository := r2.NewRepository(db)
	r2Service := r2.NewService(r2Repository, r2Client, r2PresignClient, cfg.R2.BucketName)
	r2Handler := r2.NewHandler(r2Service)

	// Initialize Webhook module
	webhookService := webhook.NewService(r2Repository)
	webhookHandler := webhook.NewHandler(webhookService, cfg.Webhook.WebhookSecret)

	return &Container{
		DB: db,

		// R2 Module
		R2Client:        r2Client,
		R2PresignClient: r2PresignClient,
		R2Repository:    r2Repository,
		R2Service:       r2Service,
		R2Handler:       r2Handler,

		// Webhook Module
		WebhookService: webhookService,
		WebhookHandler: webhookHandler,

		// User Module
		// UserRepository: userRepository,
		// UserService:    userService,
		// UserHandler:    userHandler,
	}
}

// Close - ปิด connections และ cleanup resources
func (c *Container) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
