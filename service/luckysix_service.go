package service

import (
	"log"
	"math/rand"
	"time"

	"github.com/basel-ax/luckysix/entity"
	"gorm.io/gorm"
)

// GenerateAndSaveLuckySix generates LuckySix combinations from LuckyFive and LuckyTwo.
// Each LuckySix is created by combining one LuckyFive (5 words) with one word from a LuckyTwo,
// ensuring no word repetition. The generation is resumable from the last saved combination.
func GenerateAndSaveLuckySix(db *gorm.DB) error {
	// Find the last generated LuckySix to resume from
	var lastLuckySix entity.LuckySix
	err := db.Order("lucky_five_id desc, lucky_two_id desc").First(&lastLuckySix).Error

	startLuckyFiveID := uint(0)
	startLuckyTwoID := uint(0)

	if err == nil {
		startLuckyFiveID = lastLuckySix.LuckyFiveID
		startLuckyTwoID = lastLuckySix.LuckyTwoID + 1
		log.Printf("Resuming from LuckyFiveID: %d, LuckyTwoID: %d", startLuckyFiveID, startLuckyTwoID)
	} else if err == gorm.ErrRecordNotFound {
		log.Println("No previous LuckySix found, starting from beginning")
		// Get the first LuckyFive ID
		var firstFive entity.LuckyFive
		if err := db.Order("id asc").First(&firstFive).Error; err != nil {
			return err
		}
		startLuckyFiveID = firstFive.ID
	} else {
		return err
	}

	// Load all LuckyTwo records ONCE to avoid repeated queries
	log.Println("Loading all LuckyTwo records...")
	var allLuckyTwos []entity.Luckytwo
	if err := db.Model(&entity.Luckytwo{}).Order("id asc").Select("id, word_one, word_two").Find(&allLuckyTwos).Error; err != nil {
		return err
	}

	if len(allLuckyTwos) == 0 {
		log.Println("No LuckyTwo records found. Please generate LuckyTwo first.")
		return nil
	}
	log.Printf("Loaded %d LuckyTwo records", len(allLuckyTwos))

	// Process LuckyFive in batches to manage memory
	const luckyFiveBatchSize = 100
	luckyFiveOffset := 0

	for {
		var luckyFives []entity.LuckyFive
		err := db.Where("id >= ?", startLuckyFiveID).
			Order("id asc").
			Limit(luckyFiveBatchSize).
			Offset(luckyFiveOffset).
			Find(&luckyFives).Error

		if err != nil {
			return err
		}

		if len(luckyFives) == 0 {
			break
		}

		log.Printf("Processing batch of %d LuckyFive records (offset: %d)", len(luckyFives), luckyFiveOffset)

		// Batch insert buffer
		const insertBatchSize = 500
		luckySixBatch := make([]entity.LuckySix, 0, insertBatchSize)

		// Process each LuckyFive
		for _, luckyFive := range luckyFives {
			// Create a set of words already in LuckyFive
			existingWords := map[uint]bool{
				luckyFive.WordOne:   true,
				luckyFive.WordTwo:   true,
				luckyFive.WordThree: true,
				luckyFive.WordFour:  true,
				luckyFive.WordFive:  true,
			}

			// Determine starting LuckyTwo ID for this LuckyFive
			luckyTwoStartIdx := 0
			if luckyFive.ID == startLuckyFiveID && startLuckyTwoID > 0 {
				// Find the starting index in our preloaded slice
				for i, lt := range allLuckyTwos {
					if lt.ID >= startLuckyTwoID {
						luckyTwoStartIdx = i
						break
					}
				}
			}

			// Try to combine with each LuckyTwo
			for i := luckyTwoStartIdx; i < len(allLuckyTwos); i++ {
				luckyTwo := allLuckyTwos[i]

				// Check if either word from LuckyTwo is not in LuckyFive
				word1Available := !existingWords[luckyTwo.WordOne]
				word2Available := !existingWords[luckyTwo.WordTwo]

				// Try WordOne from LuckyTwo
				if word1Available {
					luckySix := entity.LuckySix{
						LuckyFiveID: luckyFive.ID,
						LuckyTwoID:  luckyTwo.ID,
						WordOne:     luckyFive.WordOne,
						WordTwo:     luckyFive.WordTwo,
						WordThree:   luckyFive.WordThree,
						WordFour:    luckyFive.WordFour,
						WordFive:    luckyFive.WordFive,
						WordSix:     luckyTwo.WordOne,
					}

					luckySixBatch = append(luckySixBatch, luckySix)

					// Batch insert when buffer is full
					if len(luckySixBatch) >= insertBatchSize {
						if err := insertLuckySixBatch(db, luckySixBatch); err != nil {
							return err
						}
						luckySixBatch = luckySixBatch[:0] // Clear the batch
					}
				}

				// Try WordTwo from LuckyTwo (if different from WordOne and also available)
				if word2Available && luckyTwo.WordOne != luckyTwo.WordTwo {
					luckySix := entity.LuckySix{
						LuckyFiveID: luckyFive.ID,
						LuckyTwoID:  luckyTwo.ID,
						WordOne:     luckyFive.WordOne,
						WordTwo:     luckyFive.WordTwo,
						WordThree:   luckyFive.WordThree,
						WordFour:    luckyFive.WordFour,
						WordFive:    luckyFive.WordFive,
						WordSix:     luckyTwo.WordTwo,
					}

					luckySixBatch = append(luckySixBatch, luckySix)

					// Batch insert when buffer is full
					if len(luckySixBatch) >= insertBatchSize {
						if err := insertLuckySixBatch(db, luckySixBatch); err != nil {
							return err
						}
						luckySixBatch = luckySixBatch[:0] // Clear the batch
					}
				}
			}

			// Reset startLuckyTwoID after first LuckyFive
			if luckyFive.ID == startLuckyFiveID {
				startLuckyTwoID = 0
			}
		}

		// Insert any remaining records in the batch
		if len(luckySixBatch) > 0 {
			if err := insertLuckySixBatch(db, luckySixBatch); err != nil {
				return err
			}
		}

		log.Printf("Completed batch: processed %d LuckyFive records", len(luckyFives))
		luckyFiveOffset += luckyFiveBatchSize
	}

	log.Println("All LuckySix combinations generated and saved")
	return nil
}

// insertLuckySixBatch inserts a batch of LuckySix records, skipping duplicates
func insertLuckySixBatch(db *gorm.DB, batch []entity.LuckySix) error {
	if len(batch) == 0 {
		return nil
	}

	// Use individual inserts with ON CONFLICT DO NOTHING to handle duplicates
	// This is more reliable than CreateInBatches for handling unique constraint violations
	for _, ls := range batch {
		// Check if this combination already exists
		var count int64
		err := db.Model(&entity.LuckySix{}).
			Where("word_one = ? AND word_two = ? AND word_three = ? AND word_four = ? AND word_five = ? AND word_six = ?",
				ls.WordOne, ls.WordTwo, ls.WordThree, ls.WordFour, ls.WordFive, ls.WordSix).
			Count(&count).Error
		if err != nil {
			return err
		}

		// Only insert if it doesn't exist
		if count == 0 {
			if err := db.Create(&ls).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// GenerateAndSaveRandomLuckySix generates LuckySix combinations using random LuckyFive entries
// This function selects random LuckyFive entries and combines them with LuckyTwo entries
// to create LuckySix combinations, avoiding duplicates by checking lucky_five_id
func GenerateAndSaveRandomLuckySix(db *gorm.DB) error {
	// Load all LuckyTwo records
	log.Println("Loading all LuckyTwo records...")
	var allLuckyTwos []entity.Luckytwo
	if err := db.Model(&entity.Luckytwo{}).Order("id asc").Select("id, word_one, word_two").Find(&allLuckyTwos).Error; err != nil {
		return err
	}

	if len(allLuckyTwos) == 0 {
		log.Println("No LuckyTwo records found. Please generate LuckyTwo first.")
		return nil
	}
	log.Printf("Loaded %d LuckyTwo records", len(allLuckyTwos))

	// Seed random
	rand.Seed(time.Now().UnixNano())

	// Get total count of LuckyFive records
	var luckyFiveCount int64
	if err := db.Model(&entity.LuckyFive{}).Count(&luckyFiveCount).Error; err != nil {
		return err
	}

	if luckyFiveCount == 0 {
		log.Println("No LuckyFive records found. Please generate LuckyFive first.")
		return nil
	}
	log.Printf("Found %d LuckyFive records to work with", luckyFiveCount)

	// Process in batches
	const batchSize = 100
	luckySixBatch := make([]entity.LuckySix, 0, batchSize)

	// Generate random LuckySix combinations
	generated := 0
	maxGenerate := 10000 // Generate 10000 random combinations per run

	for generated < maxGenerate {
		// Select a random LuckyFive
		randomOffset := rand.Intn(int(luckyFiveCount))
		var luckyFive entity.LuckyFive
		if err := db.Offset(randomOffset).First(&luckyFive).Error; err != nil {
			return err
		}

		// Check if this LuckyFive has already been used for LuckySix generation
		var existingCount int64
		if err := db.Model(&entity.LuckySix{}).
			Where("lucky_five_id = ?", luckyFive.ID).
			Count(&existingCount).Error; err != nil {
			return err
		}

		// Skip if this LuckyFive has already been used
		if existingCount > 0 {
			log.Printf("LuckyFive ID %d already used, skipping", luckyFive.ID)
			continue
		}

		// Create a set of words already in LuckyFive
		existingWords := map[uint]bool{
			luckyFive.WordOne:   true,
			luckyFive.WordTwo:   true,
			luckyFive.WordThree: true,
			luckyFive.WordFour:  true,
			luckyFive.WordFive:  true,
		}

		// Try to combine with random LuckyTwo entries
		// Shuffle the LuckyTwo list to get random combinations
		shuffledTwos := make([]entity.Luckytwo, len(allLuckyTwos))
		copy(shuffledTwos, allLuckyTwos)

		// Fisher-Yates shuffle
		for i := len(shuffledTwos) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			shuffledTwos[i], shuffledTwos[j] = shuffledTwos[j], shuffledTwos[i]
		}

		// Try first few shuffled LuckyTwo entries
		twoCount := 0
		for _, luckyTwo := range shuffledTwos {
			// Check if either word from LuckyTwo is not in LuckyFive
			word1Available := !existingWords[luckyTwo.WordOne]
			word2Available := !existingWords[luckyTwo.WordTwo]

			// Try WordOne from LuckyTwo
			if word1Available {
				luckySix := entity.LuckySix{
					LuckyFiveID: luckyFive.ID,
					LuckyTwoID:  luckyTwo.ID,
					WordOne:     luckyFive.WordOne,
					WordTwo:     luckyFive.WordTwo,
					WordThree:   luckyFive.WordThree,
					WordFour:    luckyFive.WordFour,
					WordFive:    luckyFive.WordFive,
					WordSix:     luckyTwo.WordOne,
				}

				// Check if this specific combination already exists
				var count int64
				err := db.Model(&entity.LuckySix{}).
					Where("word_one = ? AND word_two = ? AND word_three = ? AND word_four = ? AND word_five = ? AND word_six = ?",
						luckySix.WordOne, luckySix.WordTwo, luckySix.WordThree, luckySix.WordFour, luckySix.WordFive, luckySix.WordSix).
					Count(&count).Error
				if err != nil {
					return err
				}

				// Only insert if it doesn't exist
				if count == 0 {
					luckySixBatch = append(luckySixBatch, luckySix)
					generated++
					log.Printf("Generated LuckySix from LuckyFive ID %d: %d %d %d %d %d %d",
						luckyFive.ID, luckySix.WordOne, luckySix.WordTwo, luckySix.WordThree,
						luckySix.WordFour, luckySix.WordFive, luckySix.WordSix)

					// Batch insert when buffer is full
					if len(luckySixBatch) >= batchSize {
						if err := insertLuckySixBatch(db, luckySixBatch); err != nil {
							return err
						}
						luckySixBatch = luckySixBatch[:0] // Clear the batch
					}
				}

				// Move to next LuckyFive after finding one valid combination
				break
			}

			// Try WordTwo from LuckyTwo (if different from WordOne and also available)
			if word2Available && luckyTwo.WordOne != luckyTwo.WordTwo {
				luckySix := entity.LuckySix{
					LuckyFiveID: luckyFive.ID,
					LuckyTwoID:  luckyTwo.ID,
					WordOne:     luckyFive.WordOne,
					WordTwo:     luckyFive.WordTwo,
					WordThree:   luckyFive.WordThree,
					WordFour:    luckyFive.WordFour,
					WordFive:    luckyFive.WordFive,
					WordSix:     luckyTwo.WordTwo,
				}

				// Check if this specific combination already exists
				var count int64
				err := db.Model(&entity.LuckySix{}).
					Where("word_one = ? AND word_two = ? AND word_three = ? AND word_four = ? AND word_five = ? AND word_six = ?",
						luckySix.WordOne, luckySix.WordTwo, luckySix.WordThree, luckySix.WordFour, luckySix.WordFive, luckySix.WordSix).
					Count(&count).Error
				if err != nil {
					return err
				}

				// Only insert if it doesn't exist
				if count == 0 {
					luckySixBatch = append(luckySixBatch, luckySix)
					generated++
					log.Printf("Generated LuckySix from LuckyFive ID %d: %d %d %d %d %d %d",
						luckyFive.ID, luckySix.WordOne, luckySix.WordTwo, luckySix.WordThree,
						luckySix.WordFour, luckySix.WordFive, luckySix.WordSix)

					// Batch insert when buffer is full
					if len(luckySixBatch) >= batchSize {
						if err := insertLuckySixBatch(db, luckySixBatch); err != nil {
							return err
						}
						luckySixBatch = luckySixBatch[:0] // Clear the batch
					}
				}

				// Move to next LuckyFive after finding one valid combination
				break
			}

			twoCount++
			// Limit the number of LuckyTwo attempts per LuckyFive
			if twoCount >= 10 {
				break
			}
		}
	}

	// Insert any remaining records in the batch
	if len(luckySixBatch) > 0 {
		if err := insertLuckySixBatch(db, luckySixBatch); err != nil {
			return err
		}
	}

	log.Printf("Random LuckySix generation completed. Total generated: %d", generated)
	return nil
}
