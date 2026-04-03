package service

import (
	"log"

	"github.com/alexvec/go-bip39"
	"github.com/basel-ax/luckysix/entity"
	"gorm.io/gorm"
)

func GenerateAndSaveLuckyTwo(db *gorm.DB) (int64, error) {
	// Get BIP39 word list
	wordList := bip39.GetWordList()
	wordCount := len(wordList)

	log.Printf("BIP39 word list has %d words", wordCount)

	// Find the last generated LuckyTwo to resume from
	var last entity.Luckytwo
	err := db.Order("word_one desc, word_two desc").First(&last).Error
	startI := uint(0)
	startJ := uint(0)
	if err == nil {
		startI = last.WordOne
		startJ = last.WordTwo + 1
		if startJ >= uint(wordCount) {
			startI++
			startJ = 0
		}
		log.Printf("Resuming from WordOne: %d, WordTwo: %d", startI, startJ)
	} else if err == gorm.ErrRecordNotFound {
		log.Println("No previous LuckyTwo found, starting from beginning")
	} else {
		return 0, err
	}

	// Count records before generation
	var countBefore int64
	db.Model(&entity.Luckytwo{}).Count(&countBefore)
	log.Printf("Records in database before generation: %d", countBefore)

	// Generate combinations starting from startI, startJ
	for i := startI; i < uint(wordCount); i++ {
		jStart := uint(0)
		if i == startI {
			jStart = startJ
		}
		for j := jStart; j < uint(wordCount); j++ {
			luckyTwo := entity.Luckytwo{
				WordOne: i,
				WordTwo: j,
			}
			if err := db.Create(&luckyTwo).Error; err != nil {
				return 0, err
			}
		}
		log.Printf("Completed combinations for word %d", i)
	}

	// Count records after generation
	var countAfter int64
	db.Model(&entity.Luckytwo{}).Count(&countAfter)
	log.Printf("Records in database after generation: %d", countAfter)

	generated := countAfter - countBefore
	log.Printf("New combinations generated: %d", generated)
	log.Println("All LuckyTwo combinations generated and saved")
	return generated, nil
}
