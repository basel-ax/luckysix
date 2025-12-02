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

// areAllDistinct checks if all elements in a slice are distinct
func areAllDistinct(words []uint) bool {
	seen := make(map[uint]bool)
	for _, word := range words {
		if seen[word] {
			return false
		}
		seen[word] = true
	}
	return true
}

func GenerateAndSaveLuckyFive(db *gorm.DB, all bool) error {
	// Get BIP39 word list
	wordList := bip39.GetWordList()
	wordCount := len(wordList)

	log.Printf("BIP39 word list has %d words", wordCount)

	// Seed random
	rand.Seed(time.Now().UnixNano())

	// Generate combinations based on the 'all' parameter
	if all {
		// Generate ALL possible combinations (2048^5)
		log.Println("Starting generation of ALL LuckyFive combinations - this will take a very long time!")
		generated := 0
		totalCombinations := uint64(wordCount) * uint64(wordCount) * uint64(wordCount) * uint64(wordCount) * uint64(wordCount)
		log.Printf("Total combinations to generate: %d", totalCombinations)

		// Use nested loops to generate all combinations
		for w1 := 0; w1 < wordCount; w1++ {
			for w2 := 0; w2 < wordCount; w2++ {
				for w3 := 0; w3 < wordCount; w3++ {
					for w4 := 0; w4 < wordCount; w4++ {
						for w5 := 0; w5 < wordCount; w5++ {
							// Ensure all words are distinct
							words := []uint{uint(w1), uint(w2), uint(w3), uint(w4), uint(w5)}
							if areAllDistinct(words) {
								luckyFive := entity.LuckyFive{
									WordOne:   uint(w1),
									WordTwo:   uint(w2),
									WordThree: uint(w3),
									WordFour:  uint(w4),
									WordFive:  uint(w5),
								}

								// Check if exists
								var count int64
								db.Model(&entity.LuckyFive{}).Where("word_one = ? AND word_two = ? AND word_three = ? AND word_four = ? AND word_five = ?",
									words[0], words[1], words[2], words[3], words[4]).Count(&count)
								if count == 0 {
									if err := db.Create(&luckyFive).Error; err != nil {
										return err
									}
									generated++
									if generated%1000 == 0 {
										log.Printf("Generated %d LuckyFive combinations", generated)
									}
								}
							}
						}
					}
				}
			}
		}
		log.Printf("Completed generation of all LuckyFive combinations. Total generated: %d", generated)
	} else {
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
	}

	log.Println("LuckyFive generation completed")
	return nil
}
