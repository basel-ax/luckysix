package service

import (
	"log"

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
	if err := db.Order("id asc").Find(&allLuckyTwos).Error; err != nil {
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
						if err := db.CreateInBatches(luckySixBatch, insertBatchSize).Error; err != nil {
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
						if err := db.CreateInBatches(luckySixBatch, insertBatchSize).Error; err != nil {
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
			if err := db.CreateInBatches(luckySixBatch, insertBatchSize).Error; err != nil {
				return err
			}
		}

		log.Printf("Completed batch: processed %d LuckyFive records", len(luckyFives))
		luckyFiveOffset += luckyFiveBatchSize
	}

	log.Println("All LuckySix combinations generated and saved")
	return nil
}
