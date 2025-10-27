package entity

type LuckySix struct {
	ID            uint
	LuckyFiveID   uint
	LuckyTwoID    uint
	WalletAddress string
	Status        string `gorm:"default:New"`
}
