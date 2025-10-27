package entity

import (
	"time"

	"gorm.io/gorm"
)

// WalletBalance stores the generated mnemonic, address, and balance for a LuckySix combination.
type WalletBalance struct {
	gorm.Model
	LuckySixID       uint       `gorm:"column:lucky_six_id;uniqueIndex"`
	Mnemonic         string     `gorm:"column:mnemonic"`
	Address          string     `gorm:"column:address"`
	CosmosAddress    string     `gorm:"column:cosmos_address"`
	Balance          string     `gorm:"column:balance"` // Stored as string to handle large numbers
	BalanceUpdatedAt *time.Time `gorm:"column:balance_updated_at"`
	IsNotified       bool       `gorm:"column:is_notified;default:false"`
}
