package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/basel-ax/luckysix/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDatabaseConnection(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	host := getEnvWithDefault("DB_HOST", "localhost")
	port := getEnvWithDefault("DB_PORT", "5432")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := getEnvWithDefault("DB_SSLMODE", "disable")

	if user == "" || password == "" || dbname == "" {
		t.Fatal("DB_USER, DB_PASSWORD, and DB_NAME environment variables are required")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	tables := []interface{}{
		&entity.Luckytwo{},
		&entity.LuckyFive{},
		&entity.LuckySix{},
		&entity.WalletBalance{},
	}

	migrator := db.Migrator()
	for _, table := range tables {
		if !migrator.HasTable(table) {
			if err := migrator.CreateTable(table); err != nil {
				t.Fatalf("Failed to create table: %v", err)
			}
		}
	}

	t.Log("Database connection and migration test passed")
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
