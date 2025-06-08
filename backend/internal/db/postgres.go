package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/models"
)

var DB *gorm.DB

func Init() {
	// Load .env file in dev (optional)
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate all models here
	err = db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.RecipeIngredient{},
		&models.Tag{},
		&models.RecipeTag{},
		&models.MenuEntry{},
		&models.RecipeStep{},
		&models.MenuEntryExtra{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Assign to global var
	DB = db

	fmt.Println("database connection initialized")

	models.Seed(db)

	fmt.Println("db seed succesfully completed")
}
