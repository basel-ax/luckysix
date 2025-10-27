package entity

import "gorm.io/gorm"

// LuckySix represents a combination of a LuckyFive and a LuckyTwo.
// It's recommended to add a unique constraint on (lucky_five_id, lucky_two_id) in your database migration.
type LuckySix struct {
	gorm.Model
	LuckyFiveID uint `gorm:"column:lucky_five_id"`
	LuckyTwoID  uint `gorm:"column:lucky_two_id"`
}
