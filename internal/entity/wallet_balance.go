package entity

import "github.com/shopspring/decimal"

type WalletBalance struct {
	ID           uint
	Wallet       string
	BlockchainID uint
	LuckySixID   uint
	Balance      decimal.Decimal `json:"amount" sql:"type:decimal(20,8);"`
}
package entity

import "gorm.io/gorm"

// WalletBalance stores the generated mnemonic, address, and balance for a LuckySix combination.
type WalletBalance struct {
	gorm.Model
	LuckySixID uint   `gorm:"column:lucky_six_id;uniqueIndex"`
	Mnemonic   string `gorm:"column:mnemonic"`
	Address    string `gorm:"column:address"`
	Balance    string `gorm:"column:balance"` // Stored as string to handle large numbers
}
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
	Balance          string     `gorm:"column:balance"` // Stored as string to handle large numbers
	BalanceUpdatedAt *time.Time `gorm:"column:balance_updated_at"`
}
