package entity

import "gorm.io/gorm"

type Luckytwo struct {
	gorm.Model
	WordOne uint `gorm:"uniqueIndex:idx_luckytwo"`
	WordTwo uint `gorm:"uniqueIndex:idx_luckytwo"`
}
