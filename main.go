package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/basel-ax/luckysix/entity"
	"github.com/basel-ax/luckysix/service"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var prodMode bool

// logWriter implements io.Writer for suppressing logs in prod mode
type logWriter struct{}

func (logWriter) Write(p []byte) (n int, err error) {
	if prodMode {
		return len(p), nil
	}
	return os.Stderr.Write(p)
}

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
	migrator := db.Migrator()
	tables := []interface{}{
		&entity.Luckytwo{},
		&entity.LuckyFive{},
		&entity.LuckySix{},
		&entity.WalletBalance{},
	}
	for _, t := range tables {
		if !migrator.HasTable(t) {
			if err := migrator.CreateTable(t); err != nil {
				// Handle the case where table was created by another process after our check
				if strings.Contains(err.Error(), "already exists") {
					log.Printf("Table already exists (created by another process): %T", t)
					continue
				}
				log.Fatalf("Failed to create table: %v", err)
			}
			log.Printf("Created table: %T", t)
		} else {
			log.Printf("Table already exists: %T", t)
		}
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

	rootCmd.PersistentFlags().BoolVarP(&prodMode, "prod", "p", false, "Quiet mode: suppress console output and disable Telegram")

	var luckytwoCmd = &cobra.Command{
		Use:   "luckytwo",
		Short: "Commands for LuckyTwo operations",
	}

	var generateLuckyTwoCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate LuckyTwo combinations",
		Run: func(cmd *cobra.Command, args []string) {
			// Suppress logs in prod mode (quiet for cron)
			if prodMode {
				log.SetOutput(logWriter{})
			}

			// Check for cron lock
			runFunc := func() error {
				rows, err := service.GenerateAndSaveLuckyTwo(db)
				if err != nil {
					return err
				}
				// Send Telegram notification in prod mode
				if !prodMode {
					return service.SendGenerationNotification("LuckyTwo", rows, nil)
				}
				return nil
			}

			acquired, err := service.CheckAndRunCommand("luckytwo-generate", runFunc)
			if err != nil {
				if !prodMode {
					service.SendGenerationNotification("LuckyTwo", 0, err)
				}
				log.Fatal(err)
			}
			if !acquired {
				log.Println("luckytwo generate is already running, skipping...")
				if !prodMode {
					service.SendGenerationNotification("LuckyTwo", -1, nil)
				}
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
			// Suppress logs in prod mode (quiet for cron)
			if prodMode {
				log.SetOutput(logWriter{})
			}

			all, _ := cmd.Flags().GetBool("all")

			// Check for cron lock
			runFunc := func() error {
				rows, err := service.GenerateAndSaveLuckyFive(db, all)
				if err != nil {
					return err
				}
				// Send Telegram notification in prod mode
				if !prodMode {
					return service.SendGenerationNotification("LuckyFive", rows, nil)
				}
				return nil
			}

			acquired, err := service.CheckAndRunCommand("luckyfive-generate", runFunc)
			if err != nil {
				if !prodMode {
					service.SendGenerationNotification("LuckyFive", 0, err)
				}
				log.Fatal(err)
			}
			if !acquired {
				log.Println("luckyfive generate is already running, skipping...")
				if !prodMode {
					service.SendGenerationNotification("LuckyFive", -1, nil)
				}
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
			// Suppress logs in prod mode (quiet for cron)
			if prodMode {
				log.SetOutput(logWriter{})
			}

			// Check for cron lock
			runFunc := func() error {
				rows, err := service.GenerateAndSaveLuckySix(db)
				if err != nil {
					return err
				}
				// Send Telegram notification in prod mode
				if !prodMode {
					return service.SendGenerationNotification("LuckySix", rows, nil)
				}
				return nil
			}

			acquired, err := service.CheckAndRunCommand("luckysix-generate", runFunc)
			if err != nil {
				if !prodMode {
					service.SendGenerationNotification("LuckySix", 0, err)
				}
				log.Fatal(err)
			}
			if !acquired {
				log.Println("luckysix generate is already running, skipping...")
				if !prodMode {
					service.SendGenerationNotification("LuckySix", -1, nil)
				}
			}
		},
	}

	var generateRandomLuckySixCmd = &cobra.Command{
		Use:   "generate-random",
		Short: "Generate LuckySix combinations using random LuckyFive entries",
		Run: func(cmd *cobra.Command, args []string) {
			// Suppress logs in prod mode (quiet for cron)
			if prodMode {
				log.SetOutput(logWriter{})
			}

			if err := service.GenerateAndSaveRandomLuckySix(db); err != nil {
				log.Fatal(err)
			}
		},
	}

	luckysixCmd.AddCommand(generateLuckySixCmd)
	luckysixCmd.AddCommand(generateRandomLuckySixCmd)

	var walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Commands for Wallet operations",
	}

	var generateWalletCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate 12-word wallet mnemonics from LuckySix combinations",
		Run: func(cmd *cobra.Command, args []string) {
			// Suppress logs in prod mode (quiet for cron)
			if prodMode {
				log.SetOutput(logWriter{})
			}

			count, _ := cmd.Flags().GetInt("count")
			if count == 0 {
				count = 100000 // Default count
			}

			// Check for cron lock
			runFunc := func() error {
				rows, err := service.GenerateWalletsFromLuckySix(db, count)
				if err != nil {
					return err
				}
				// Send Telegram notification in prod mode
				if !prodMode {
					return service.SendGenerationNotification("Wallet", rows, nil)
				}
				return nil
			}

			acquired, err := service.CheckAndRunCommand("wallet-generate", runFunc)
			if err != nil {
				if !prodMode {
					service.SendGenerationNotification("Wallet", 0, err)
				}
				log.Fatal(err)
			}
			if !acquired {
				log.Println("wallet generate is already running, skipping...")
				if !prodMode {
					service.SendGenerationNotification("Wallet", -1, nil)
				}
			}
		},
	}
	generateWalletCmd.Flags().IntP("count", "c", 100000, "Number of wallets to generate")

walletCmd.AddCommand(generateWalletCmd)

	var statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Send project statistics to Telegram",
		Run: func(cmd *cobra.Command, args []string) {
			// Suppress logs in prod mode (quiet for cron)
			if prodMode {
				log.SetOutput(logWriter{})
			}

			// Get stats from database
			var luckytwosCount, luckyfivesCount, luckysixesCount, walletsCount, walletsWithAddress, walletsWithBalance int64

			// Count luckytwos
			db.Model(&entity.Luckytwo{}).Count(&luckytwosCount)
			// Count luckyfives
			db.Model(&entity.LuckyFive{}).Count(&luckyfivesCount)
			// Count luckysixes
			db.Model(&entity.LuckySix{}).Count(&luckysixesCount)
			// Count total wallets
			db.Model(&entity.WalletBalance{}).Count(&walletsCount)
			// Count wallets with non-empty address
			db.Model(&entity.WalletBalance{}).Where("address IS NOT NULL AND address != ''").Count(&walletsWithAddress)
			// Count wallets with balance > 0
			db.Model(&entity.WalletBalance{}).Where("balance IS NOT NULL AND balance != '0' AND balance != ''").Count(&walletsWithBalance)

			// Format the message
			message := fmt.Sprintf(`📊 LuckySix Project Statistics

🔢 Combinations Generated:
• LuckyTwos: %d
• LuckyFives: %d
• LuckySixes: %d

💰 Wallets:
• Total Wallets: %d
• With Address: %d
• With Balance > 0: %d`, luckytwosCount, luckyfivesCount, luckysixesCount, walletsCount, walletsWithAddress, walletsWithBalance)

			// Print to stdout
			fmt.Println(message)

			// Send to Telegram in debug mode (when --prod is NOT set)
			if !prodMode {
				if err := service.TelegramNotification(message); err != nil {
					log.Printf("Failed to send Telegram notification: %v", err)
				}
			}
		},
	}

	rootCmd.AddCommand(luckytwoCmd)
	rootCmd.AddCommand(luckyfiveCmd)
	rootCmd.AddCommand(luckysixCmd)
	rootCmd.AddCommand(walletCmd)
	rootCmd.AddCommand(statsCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Suppress unused variable warning
var _ = io.Writer(nil)
