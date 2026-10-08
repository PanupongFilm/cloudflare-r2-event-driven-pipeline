package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/config"
	"server/database"
	"server/internal/r2"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// ตั้งค่า log prefix
	log.SetPrefix("[server] ")
	log.SetFlags(log.Ldate | log.Ltime)

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Starting server...")
	log.Printf("Environment: %s", cfg.Server.Environment)

	// Connect to database (Required - server won't start without DB)
	db, err := database.Connect(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("❌Failed to connect database: %v", err)
	}
	log.Println("Database connected successfully")
	
	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

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
	log.Println("R2 client initialized successfully")

	// สร้าง Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "R2 Event Pipeline API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowOrigins,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	}))

	// Routes
	setupRoutes(app, db, r2Client)

	// Graceful shutdown
	go func() {
		if err := app.Listen(cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server is running on port %s", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if db != nil {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	log.Println("Server stopped gracefully")
}

func setupRoutes(app *fiber.App, db interface{}, r2Client interface{}) {
	// Health check
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "R2 Event-Driven Pipeline API",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API Routes
	api := app.Group("/api/v1")

	// TODO: Add your routes here
	_ = api
	_ = r2Client

	// 404 Handler
	app.Use(func(c fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error": "Route not found",
			"path":  c.Path(),
		})
	})
}

func customErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}
