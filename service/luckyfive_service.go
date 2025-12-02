package service

import (
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/alexvec/go-bip39"
	"github.com/basel-ax/luckysix/entity"
	"gorm.io/gorm"
)

func GenerateAndSaveLuckyFive(db *gorm.DB) error {
	// Get BIP39 word list
	wordList := bip39.GetWordList()
	wordCount := len(wordList)

	log.Printf("BIP39 word list has %d words", wordCount)

	// Seed random
	rand.Seed(time.Now().UnixNano())

	// Generate random combinations, say 100 per run
	generated := 0
	maxGenerate := 100
	for generated < maxGenerate {
		// Generate 5 distinct random indices
		indices := make(map[uint]bool)
		for len(indices) < 5 {
			idx := uint(rand.Intn(wordCount))
			indices[idx] = true
		}

		words := make([]uint, 0, 5)
		for idx := range indices {
			words = append(words, idx)
		}
		sort.Slice(words, func(i, j int) bool { return words[i] < words[j] })

		luckyFive := entity.LuckyFive{
			WordOne:   words[0],
			WordTwo:   words[1],
			WordThree: words[2],
			WordFour:  words[3],
			WordFive:  words[4],
		}

		// Check if exists
		var count int64
		db.Model(&entity.LuckyFive{}).Where("word_one = ? AND word_two = ? AND word_three = ? AND word_four = ? AND word_five = ?",
			words[0], words[1], words[2], words[3], words[4]).Count(&count)
		if count == 0 {
			if err := db.Create(&luckyFive).Error; err != nil {
				return err
			}
			log.Printf("Generated LuckyFive: %d %d %d %d %d", words[0], words[1], words[2], words[3], words[4])
			generated++
		} else {
			log.Println("Duplicate found, skipping")
		}
	}

	log.Println("LuckyFive generation completed")
	return nil
}
