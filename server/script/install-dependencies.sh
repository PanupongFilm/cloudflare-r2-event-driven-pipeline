#!/bin/bash

echo "Installing Go dependencies..."

# Environment variables
go get github.com/joho/godotenv

# Fiber v3
go get github.com/gofiber/fiber/v3
go get github.com/gofiber/fiber/v3/middleware/cors
go get github.com/gofiber/fiber/v3/middleware/logger
go get github.com/gofiber/fiber/v3/middleware/recover

# Validation
go get github.com/go-playground/validator/v10

# AWS SDK v2 for R2 (Cloudflare R2 Storage)
go get github.com/aws/aws-sdk-go-v2/aws
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/credentials
go get github.com/aws/aws-sdk-go-v2/service/s3

# GORM
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Clean up and download
go mod tidy
go mod download

echo "✅ All dependencies installed successfully!"
