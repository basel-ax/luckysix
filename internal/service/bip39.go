package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/alexvec/go-bip39"
	"gorm.io/gorm"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/logger"
)

const (
	batchSize = 1000
)

// Bip39Service -.
type Bip39Service struct {
	luckyTwoRepo  LuckyTwoRepo
	luckyFiveRepo LuckyFiveRepo
	luckySixRepo  LuckySixRepo
	logger        logger.Interface
}

// NewBip39Service -.
func NewBip39Service(luckyTwoRepo LuckyTwoRepo, luckyFiveRepo LuckyFiveRepo, luckySixRepo LuckySixRepo, l logger.Interface) *Bip39Service {
	return &Bip39Service{
		luckyTwoRepo:  luckyTwoRepo,
		luckyFiveRepo: luckyFiveRepo,
		luckySixRepo:  luckySixRepo,
		logger:        l,
	}
}

// GenerateAndStoreLuckyTwoCombinations generates all 2-word combinations from BIP39 word list indices and stores them.
func (s *Bip39Service) GenerateAndStoreLuckyTwoCombinations(ctx context.Context) error {
	s.logger.Info("Starting generation of LuckyTwo combinations")

	wordListSize := len(bip39.GetWordList())

	batch := make([]entity.Luckytwo, 0, batchSize)

	for i := 0; i < wordListSize; i++ {
		for j := 0; j < wordListSize; j++ {
			batch = append(batch, entity.Luckytwo{
				WordOne: uint(i),
				WordTwo: uint(j),
			})

			if len(batch) >= batchSize {
				if err := s.luckyTwoRepo.StoreBatch(ctx, batch); err != nil {
					err = fmt.Errorf("failed to store batch for i=%d: %w", i, err)
					s.logger.Error(err)
					return err
				}
				batch = make([]entity.Luckytwo, 0, batchSize) // Reset batch
			}
		}
	}

	if len(batch) > 0 { // Store remaining items
		if err := s.luckyTwoRepo.StoreBatch(ctx, batch); err != nil {
			err = fmt.Errorf("failed to store final batch: %w", err)
			s.logger.Error(err)
			return err
		}
	}

	s.logger.Info("Finished generation of LuckyTwo combinations")
	return nil
}

// GenerateAndStoreLuckyFiveCombinations generates and stores a specified number of random LuckyFive combinations.
func (s *Bip39Service) GenerateAndStoreLuckyFiveCombinations(ctx context.Context, count int) error {
	s.logger.Info("Starting generation of LuckyFive combinations", "count", count)

	totalLuckyTwo, err := s.luckyTwoRepo.Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to get count of LuckyTwo pairs: %w", err)
	}

	if totalLuckyTwo < 5 {
		return fmt.Errorf("not enough LuckyTwo pairs to generate LuckyFive, found %d", totalLuckyTwo)
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	batch := make([]entity.LuckyFive, 0, batchSize)

	for i := 0; i < count; i++ {
		pickedIDs := make(map[uint]struct{}, 5)
		pairs := make([]uint, 0, 5)

		for len(pairs) < 5 {
			// Generate a random ID between 1 and totalLuckyTwo (inclusive)
			randomID := uint(random.Int63n(totalLuckyTwo) + 1)
			if _, exists := pickedIDs[randomID]; !exists {
				pickedIDs[randomID] = struct{}{}
				pairs = append(pairs, randomID)
			}
		}

		luckyFive := entity.LuckyFive{
			PairOne:   pairs[0],
			PairTwo:   pairs[1],
			PairThree: pairs[2],
			PairFour:  pairs[3],
			PairFive:  pairs[4],
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		batch = append(batch, luckyFive)

		if len(batch) >= batchSize {
			if err := s.luckyFiveRepo.StoreBatch(ctx, batch); err != nil {
				s.logger.Error(fmt.Errorf("failed to store batch: %w", err))
				return err
			}
			batch = make([]entity.LuckyFive, 0, batchSize) // Reset batch
		}
	}

	// Store any remaining items in the batch
	if len(batch) > 0 {
		if err := s.luckyFiveRepo.StoreBatch(ctx, batch); err != nil {
			s.logger.Error(fmt.Errorf("failed to store final batch: %w", err))
			return err
		}
	}

	s.logger.Info("Finished generation of LuckyFive combinations")
	return nil
}

// GenerateAndStoreLuckySixCombinations generates and stores a specified number of random LuckySix combinations.
func (s *Bip39Service) GenerateAndStoreLuckySixCombinations(ctx context.Context, count int) error {
	s.logger.Info("Starting generation of LuckySix combinations", "count", count)

	totalLuckyTwo, err := s.luckyTwoRepo.Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to get count of LuckyTwo pairs: %w", err)
	}
	totalLuckyFive, err := s.luckyFiveRepo.Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to get count of LuckyFive combinations: %w", err)
	}

	if totalLuckyTwo == 0 || totalLuckyFive == 0 {
		return fmt.Errorf("not enough LuckyTwo (%d) or LuckyFive (%d) to generate LuckySix", totalLuckyTwo, totalLuckyFive)
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	batch := make([]entity.LuckySix, 0, batchSize)

	for i := 0; i < count; i++ {
		randomTwoID := uint(random.Int63n(totalLuckyTwo) + 1)
		randomFiveID := uint(random.Int63n(totalLuckyFive) + 1)

		luckySix := entity.LuckySix{
			LuckyFiveID: randomFiveID,
			LuckyTwoID:  randomTwoID,
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		batch = append(batch, luckySix)

		if len(batch) >= batchSize {
			if err := s.luckySixRepo.StoreBatch(ctx, batch); err != nil {
				s.logger.Error(fmt.Errorf("failed to store LuckySix batch: %w", err))
				return err
			}
			batch = make([]entity.LuckySix, 0, batchSize)
		}
	}

	if len(batch) > 0 {
		if err := s.luckySixRepo.StoreBatch(ctx, batch); err != nil {
			s.logger.Error(fmt.Errorf("failed to store final LuckySix batch: %w", err))
			return err
		}
	}

	s.logger.Info("Finished generation of LuckySix combinations")
	return nil
}
