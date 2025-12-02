package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexvec/go-bip39"
	"github.com/basel-ax/luckysix/entity"
	"gorm.io/gorm"
)

// GenerateWalletsFromLuckySix generates 12-word BIP39 wallet mnemonics from LuckySix combinations.
// It uses the 6 word indices from LuckySix as the first 6 words and generates the remaining 6 words
// to create a valid 12-word mnemonic. The generation is resumable from the last saved wallet.
func GenerateWalletsFromLuckySix(db *gorm.DB) error {
	// Get the BIP39 wordlist
	wordlist := bip39.GetWordList()
	if len(wordlist) != 2048 {
		return fmt.Errorf("expected 2048 words in BIP39 wordlist, got %d", len(wordlist))
	}

	// Find the last processed LuckySix ID from WalletBalance
	var lastWallet entity.WalletBalance
	err := db.Order("lucky_six_id desc").First(&lastWallet).Error

	startLuckySixID := uint(0)
	if err == nil {
		startLuckySixID = lastWallet.LuckySixID + 1
		log.Printf("Resuming from LuckySixID: %d", startLuckySixID)
	} else if err == gorm.ErrRecordNotFound {
		log.Println("No previous wallets found, starting from beginning")
		// Get the first LuckySix ID
		var firstSix entity.LuckySix
		if err := db.Order("id asc").First(&firstSix).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println("No LuckySix records found. Please generate LuckySix first.")
				return nil
			}
			return err
		}
		startLuckySixID = firstSix.ID
	} else {
		return err
	}

	// Get LuckySix records starting from startLuckySixID
	const batchSize = 1000
	offset := 0

	for {
		var luckySixes []entity.LuckySix
		err := db.Where("id >= ?", startLuckySixID).
			Order("id asc").
			Limit(batchSize).
			Offset(offset).
			Find(&luckySixes).Error

		if err != nil {
			return err
		}

		if len(luckySixes) == 0 {
			break
		}

		log.Printf("Processing batch of %d LuckySix records (offset: %d)", len(luckySixes), offset)

		for _, luckySix := range luckySixes {
			// Check if wallet already exists for this LuckySix
			var count int64
			db.Model(&entity.WalletBalance{}).Where("lucky_six_id = ?", luckySix.ID).Count(&count)
			if count > 0 {
				continue
			}

			// Convert word indices to words (BIP39 uses 0-indexed, but stored as 1-indexed)
			// Validate that all word indices are within valid range (1-2048)
			if luckySix.WordOne < 1 || luckySix.WordOne > 2048 ||
			   luckySix.WordTwo < 1 || luckySix.WordTwo > 2048 ||
			   luckySix.WordThree < 1 || luckySix.WordThree > 2048 ||
			   luckySix.WordFour < 1 || luckySix.WordFour > 2048 ||
			   luckySix.WordFive < 1 || luckySix.WordFive > 2048 ||
			   luckySix.WordSix < 1 || luckySix.WordSix > 2048 {
				log.Printf("Invalid word indices for LuckySix ID %d: %d,%d,%d,%d,%d,%d", luckySix.ID,
					luckySix.WordOne, luckySix.WordTwo, luckySix.WordThree,
					luckySix.WordFour, luckySix.WordFive, luckySix.WordSix)
				continue
			}

			words := []string{
				wordlist[luckySix.WordOne-1],
				wordlist[luckySix.WordTwo-1],
				wordlist[luckySix.WordThree-1],
				wordlist[luckySix.WordFour-1],
				wordlist[luckySix.WordFive-1],
				wordlist[luckySix.WordSix-1],
			}

			// Generate a valid 12-word mnemonic using the first 6 words
			// We need to generate the remaining 6 words that make it valid
			mnemonic, err := generateValidMnemonic(words, wordlist)
			if err != nil {
				log.Printf("Failed to generate mnemonic for LuckySix ID %d: %v", luckySix.ID, err)
				continue
			}

			// Save wallet to database
			wallet := entity.WalletBalance{
				LuckySixID: luckySix.ID,
				Mnemonic:   mnemonic,
			}

			if err := db.Create(&wallet).Error; err != nil {
				return err
			}
		}

		log.Printf("Completed batch: processed %d wallets", len(luckySixes))
		offset += batchSize
	}

	log.Println("All wallets generated and saved")
	return nil
}

// generateValidMnemonic takes the first 6 words and generates a valid 12-word BIP39 mnemonic.
// It uses a deterministic approach to find valid word combinations efficiently.
func generateValidMnemonic(firstSixWords []string, wordlist []string) (string, error) {
	// BIP39 uses entropy bits + checksum to generate mnemonics
	// For 12 words: 128 bits entropy + 4 bits checksum = 132 bits total
	// Each word represents 11 bits (2048 = 2^11)
	// First 6 words = 66 bits, we need to find remaining 6 words (66 bits)
	// The last word contains 7 bits of entropy + 4 bits checksum
	
	// Strategy: Try combinations with a smarter approach
	// We'll iterate through reasonable combinations and check validity
	maxAttempts := 100000 // Limit attempts to avoid infinite loops
	
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Use a deterministic but varied approach
		// Calculate indices based on the attempt number
		idx6 := attempt % len(wordlist)
		idx7 := (attempt / len(wordlist)) % len(wordlist)
		idx8 := (attempt / (len(wordlist) * len(wordlist))) % len(wordlist)
		idx9 := (attempt / (len(wordlist) * len(wordlist) * len(wordlist))) % len(wordlist)
		idx10 := (attempt / (len(wordlist) * len(wordlist) * len(wordlist) * len(wordlist))) % len(wordlist)
		
		// For the last word, we need to find one that makes the checksum valid
		// Try different values for the 11th word
		for idx11 := 0; idx11 < len(wordlist); idx11++ {
			// Build the 12-word mnemonic
			words := make([]string, 12)
			copy(words, firstSixWords)
			words[6] = wordlist[idx6]
			words[7] = wordlist[idx7]
			words[8] = wordlist[idx8]
			words[9] = wordlist[idx9]
			words[10] = wordlist[idx10]
			words[11] = wordlist[idx11]

			mnemonic := strings.Join(words, " ")

			// Check if it's a valid BIP39 mnemonic
			if bip39.IsMnemonicValid(mnemonic) {
				return mnemonic, nil
			}
		}
	}

	return "", fmt.Errorf("failed to generate valid mnemonic for words: %v after %d attempts", firstSixWords, maxAttempts)
}
