package entity

import "gorm.io/gorm"

// LuckyFive represents a combination of 5 LuckyTwo pair IDs.
type LuckyFive struct {
	gorm.Model
	PairOne   uint `gorm:"column:pair_one"`
	PairTwo   uint `gorm:"column:pair_two"`
	PairThree uint `gorm:"column:pair_three"`
	PairFour  uint `gorm:"column:pair_four"`
	PairFive  uint `gorm:"column:pair_five"`
}
