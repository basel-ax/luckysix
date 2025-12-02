package entity

import "gorm.io/gorm"

// LuckySix represents a combination of 6 BIP39 word indices,
// created from one LuckyFive and one additional word from a LuckyTwo.
type LuckySix struct {
	gorm.Model
	LuckyFiveID uint `gorm:"index:idx_luckysix_source"`
	LuckyTwoID  uint `gorm:"index:idx_luckysix_source"`
	WordOne     uint `gorm:"uniqueIndex:idx_luckysix"`
	WordTwo     uint `gorm:"uniqueIndex:idx_luckysix"`
	WordThree   uint `gorm:"uniqueIndex:idx_luckysix"`
	WordFour    uint `gorm:"uniqueIndex:idx_luckysix"`
	WordFive    uint `gorm:"uniqueIndex:idx_luckysix"`
	WordSix     uint `gorm:"uniqueIndex:idx_luckysix"`
}
