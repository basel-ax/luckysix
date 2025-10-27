package app

import (
	"log"

	"github.com/basel-ax/luckysix/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	databaseURL := "postgres://luckysix:luckysix@127.0.0.1:5432/luckysix"
	databaseURL += "?sslmode=disable"

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("gorm.Open error: %s", err)
	}

	db.AutoMigrate(&entity.Task{})
	log.Print("AutoMigrate: entity.Task")

	db.AutoMigrate(&entity.Luckytwo{})
	log.Print("AutoMigrate: entity.Luckytwo")

	db.AutoMigrate(&entity.LuckyFive{})
	log.Print("AutoMigrate: entity.LuckyFive")

	db.AutoMigrate(&entity.LuckySix{})
	log.Print("AutoMigrate: entity.LuckySix")

	db.AutoMigrate(&entity.Blockchain{})
	log.Print("AutoMigrate: entity.Blockchain")

	db.AutoMigrate(&entity.WalletBalance{})
	log.Print("AutoMigrate: entity.WalletBalance")
}
