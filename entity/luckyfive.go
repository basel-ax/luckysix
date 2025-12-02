package entity

import "gorm.io/gorm"

// LuckyFive represents a combination of 5 BIP39 word indices.
type LuckyFive struct {
	gorm.Model
	WordOne   uint `gorm:"uniqueIndex:idx_luckyfive"`
	WordTwo   uint `gorm:"uniqueIndex:idx_luckyfive"`
	WordThree uint `gorm:"uniqueIndex:idx_luckyfive"`
	WordFour  uint `gorm:"uniqueIndex:idx_luckyfive"`
	WordFive  uint `gorm:"uniqueIndex:idx_luckyfive"`
}
