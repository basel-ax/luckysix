package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/basel-ax/luckysix/entity"
	"github.com/basel-ax/luckysix/service"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGenerateAndSaveLuckyTwo(t *testing.T) {
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

	// Ensure tables exist
	migrator := db.Migrator()
	tables := []interface{}{
		&entity.Luckytwo{},
	}
	for _, table := range tables {
		if !migrator.HasTable(table) {
			if err := migrator.CreateTable(table); err != nil {
				t.Fatalf("Failed to create table: %v", err)
			}
		}
	}

	// Clear any existing data for clean test state
	db.Where("1=1").Delete(&entity.Luckytwo{})

	// Test the function (this will generate all combinations - be careful in real env)
	// For safety in test env, we'll skip the actual generation and just verify setup works
	_, err = service.GenerateAndSaveLuckyTwo(db)
	if err != nil {
		t.Fatalf("GenerateAndSaveLuckyTwo returned error: %v", err)
	}

	t.Log("GenerateAndSaveLuckyTwo test completed")
}
