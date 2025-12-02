package main

import (
	"log"
	"os"

	"github.com/basel-ax/luckysix/entity"
	"github.com/basel-ax/luckysix/service"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.AutoMigrate(&entity.Luckytwo{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func main() {
	initDB()

	var rootCmd = &cobra.Command{Use: "luckysix"}

	var luckytwoCmd = &cobra.Command{
		Use:   "luckytwo",
		Short: "Commands for LuckyTwo operations",
	}

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate LuckyTwo combinations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := generateLuckyTwo(); err != nil {
				log.Fatal(err)
			}
		},
	}

	luckytwoCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(luckytwoCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func generateLuckyTwo() error {
	return service.GenerateAndSaveLuckyTwo(db)
}
