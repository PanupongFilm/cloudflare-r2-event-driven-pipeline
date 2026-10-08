package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	R2       R2Config
	Database DatabaseConfig
}

type ServerConfig struct {
	Port         string
	Environment  string
	AllowOrigins []string
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Endpoint        string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig from environment variable
func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", ":3000"),
			Environment: getEnv("ENVIRONMENT", "development"),
			AllowOrigins: []string{
				getEnv("ALLOWED_ORIGIN_1", "http://localhost:3001"),
				getEnv("ALLOWED_ORIGIN_2", "http://localhost:3000"),
			},
		},
		R2: R2Config{
			AccountID:       mustGetEnv("R2_ACCOUNT_ID"),
			AccessKeyID:     mustGetEnv("R2_ACCESS_KEY_ID"),
			SecretAccessKey: mustGetEnv("R2_SECRET_ACCESS_KEY"),
			BucketName:      getEnv("R2_BUCKET_NAME", "my-bucket"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "r2_pipeline"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	// Create R2 Endpoint
	cfg.R2.Endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2.AccountID)

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}


func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}


func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}


func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}


func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}
