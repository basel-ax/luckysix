package entity

import "gorm.io/gorm"

// Luckyfive represents a combination of 5 LuckyTwo pair IDs.
type Luckyfive struct {
	gorm.Model
	PairOne   uint `gorm:"column:pair_one"`
	PairTwo   uint `gorm:"column:pair_two"`
	PairThree uint `gorm:"column:pair_three"`
	PairFour  uint `gorm:"column:pair_four"`
	PairFive  uint `gorm:"column:pair_five"`
}
