package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedUsers creates initial user data
func SeedUsers(db *gorm.DB) error {
	log.Println("Seeding users...")

	// Check if users already exist
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		log.Println("Users already exist, skipping seed")
		return nil
	}

	// Hash passwords
	password1, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	password2, _ := bcrypt.GenerateFromPassword([]byte("password456"), bcrypt.DefaultCost)

	// Create seed users
	users := []User{
		{
			UserName: "john_doe",
			Email:    "john@example.com",
			Password: string(password1),
		},
		{
			UserName: "jane_smith",
			Email:    "jane@example.com",
			Password: string(password2),
		},
	}

	// Insert users
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.UserName, err)
			return err
		}
		log.Printf("✓ Created user: %s", user.UserName)
	}

	log.Println("User seeding completed successfully")
	return nil
}

// RunSeeds runs all seed functions
func RunSeeds(db *gorm.DB) error {
	log.Println("Running database seeds...")

	if err := SeedUsers(db); err != nil {
		return err
	}

	log.Println("All seeds completed successfully")
	return nil
}
