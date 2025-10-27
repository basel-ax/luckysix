package entity

import "github.com/shopspring/decimal"

type WalletBalance struct {
	ID           uint
	Wallet       string
	BlockchainID uint
	LuckySixID   uint
	Balance      decimal.Decimal `json:"amount" sql:"type:decimal(20,8);"`
}
