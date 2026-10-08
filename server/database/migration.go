package database

import (
	"log"

	"gorm.io/gorm"
)

// RunMigrations
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto migrate models
	err := db.AutoMigrate(
		&User{},
		&UploadLog{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// DropAllTables (Dev Only)
func DropAllTables(db *gorm.DB) error {
	log.Println("Dropping all tables...")

	err := db.Migrator().DropTable(
		&UploadLog{}, 
		&User{},     
	)

	if err != nil {
		return err
	}

	log.Println("All tables dropped successfully")
	return nil
}

// ResetDatabase (Dev Only)
func ResetDatabase(db *gorm.DB) error {
	if err := DropAllTables(db); err != nil {
		return err
	}
	return RunMigrations(db)
}
