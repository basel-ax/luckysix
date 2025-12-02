package main

import (
	"fmt"
	"log"
	"os"

	"github.com/basel-ax/luckysix/entity"
	"github.com/basel-ax/luckysix/service"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Build DSN from individual env vars
		host := getEnvWithDefault("DB_HOST", "localhost")
		port := getEnvWithDefault("DB_PORT", "5432")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := getEnvWithDefault("DB_SSLMODE", "disable")

		if user == "" || password == "" || dbname == "" {
			log.Fatal("DB_USER, DB_PASSWORD, and DB_NAME environment variables are required when DATABASE_URL is not set")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	}
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.AutoMigrate(&entity.Luckytwo{}, &entity.LuckyFive{}, &entity.LuckySix{}, &entity.WalletBalance{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load .env file if it exists
	godotenv.Load()

	initDB()

	var rootCmd = &cobra.Command{Use: "luckysix"}

	var luckytwoCmd = &cobra.Command{
		Use:   "luckytwo",
		Short: "Commands for LuckyTwo operations",
	}

	var generateLuckyTwoCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate LuckyTwo combinations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := generateLuckyTwo(); err != nil {
				log.Fatal(err)
			}
		},
	}

	luckytwoCmd.AddCommand(generateLuckyTwoCmd)

	var luckyfiveCmd = &cobra.Command{
		Use:   "luckyfive",
		Short: "Commands for LuckyFive operations",
	}

	var generateLuckyFiveCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate LuckyFive combinations",
		Run: func(cmd *cobra.Command, args []string) {
			all, _ := cmd.Flags().GetBool("all")
			if err := generateLuckyFive(all); err != nil {
				log.Fatal(err)
			}
		},
	}
	generateLuckyFiveCmd.Flags().BoolP("all", "a", false, "Generate all possible LuckyFive combinations instead of random samples")

	luckyfiveCmd.AddCommand(generateLuckyFiveCmd)

	var luckysixCmd = &cobra.Command{
		Use:   "luckysix",
		Short: "Commands for LuckySix operations",
	}

	var generateLuckySixCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate LuckySix combinations from LuckyFive and LuckyTwo",
		Run: func(cmd *cobra.Command, args []string) {
			if err := generateLuckySix(); err != nil {
				log.Fatal(err)
			}
		},
	}

	luckysixCmd.AddCommand(generateLuckySixCmd)

	var walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Commands for Wallet operations",
	}

	var generateWalletCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate 12-word wallet mnemonics from LuckySix combinations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := generateWallets(); err != nil {
				log.Fatal(err)
			}
		},
	}

	walletCmd.AddCommand(generateWalletCmd)
	rootCmd.AddCommand(luckytwoCmd)
	rootCmd.AddCommand(luckyfiveCmd)
	rootCmd.AddCommand(luckysixCmd)
	rootCmd.AddCommand(walletCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func generateLuckyTwo() error {
	return service.GenerateAndSaveLuckyTwo(db)
}

func generateLuckyFive(all bool) error {
	return service.GenerateAndSaveLuckyFive(db, all)
}

func generateLuckySix() error {
	return service.GenerateAndSaveLuckySix(db)
}

func generateWallets() error {
	return service.GenerateWalletsFromLuckySix(db)
}
