package entity

import "gorm.io/gorm"

// WalletBalance represents a 12-word BIP39 wallet mnemonic generated from a LuckySix combination.
type WalletBalance struct {
	gorm.Model
	LuckySixID uint   `gorm:"index:idx_wallet_luckysix;uniqueIndex:idx_wallet_mnemonic"`
	Mnemonic   string `gorm:"type:text;uniqueIndex:idx_wallet_mnemonic"`
}
