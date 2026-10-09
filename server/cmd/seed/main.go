package main

import (
	"flag"
	"log"
	"server/config"
	"server/database"
)

func main() {
	// ตั้งค่า log prefix
	log.SetPrefix("[seed] ")
	log.SetFlags(log.Ldate | log.Ltime)

	// Define flags
	reset := flag.Bool("reset", false, "Reset database before seeding")
	flag.Parse()

	log.Println("Starting database seeding...")

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("❌ Failed to connect database: %v", err)
	}
	defer database.Close()
	log.Println("✓ Database connected successfully")

	// Reset database if flag is set
	if *reset {
		log.Println("Resetting database...")
		if err := database.ResetDatabase(db); err != nil {
			log.Fatalf("❌ Failed to reset database: %v", err)
		}
		log.Println("✓ Database reset completed")
	} else {
		// Run migrations
		if err := database.RunMigrations(db); err != nil {
			log.Fatalf("❌ Failed to run migrations: %v", err)
		}
		log.Println("✓ Database migrations completed")
	}

	// Run seeds
	if err := database.RunSeeds(db); err != nil {
		log.Fatalf("❌ Failed to run seeds: %v", err)
	}

	log.Println("🎉 Seeding completed successfully!")
}
