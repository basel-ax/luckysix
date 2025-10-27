package service

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/alexvec/go-bip39"
	"github.com/shomali11/parallelizer"
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

	wordListSize := len(bip39.WordList)

	poolSize := runtime.NumCPU()
	group := parallelizer.NewGroup(parallelizer.WithPoolSize(poolSize))
	defer group.Close()

	for i := 0; i < wordListSize; i++ {
		i := i // Capture loop variable for the closure
		err := group.Add(func() error {
			batch := make([]entity.Luckytwo, 0, batchSize)
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

			if len(batch) > 0 { // Store remaining items
				if err := s.luckyTwoRepo.StoreBatch(ctx, batch); err != nil {
					err = fmt.Errorf("failed to store final batch for i=%d: %w", i, err)
					s.logger.Error(err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			s.logger.Error(fmt.Errorf("failed to add task to parallelizer group for i=%d: %w", i, err))
		}
	}

	if err := group.Wait(); err != nil {
		err = fmt.Errorf("error during LuckyTwo combination generation: %w", err)
		s.logger.Error(err)
		return err
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

	poolSize := runtime.NumCPU()
	group := parallelizer.NewGroup(parallelizer.WithPoolSize(poolSize))
	defer group.Close()

	// Distribute the work among workers
	workPerWorker := count / poolSize
	extraWork := count % poolSize

	for i := 0; i < poolSize; i++ {
		workerNum := i
		jobs := workPerWorker
		if workerNum < extraWork {
			jobs++
		}
		if jobs == 0 {
			continue
		}

		err := group.Add(func() error {
			// Seed the random number generator for each goroutine to ensure randomness.
			source := rand.NewSource(time.Now().UnixNano() + int64(workerNum))
			random := rand.New(source)
			batch := make([]entity.LuckyFive, 0, batchSize)

			for j := 0; j < jobs; j++ {
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
						s.logger.Error(fmt.Errorf("worker %d failed to store batch: %w", workerNum, err))
						return err // Propagate error to parallelizer
					}
					batch = make([]entity.LuckyFive, 0, batchSize) // Reset batch
				}
			}

			// Store any remaining items in the batch
			if len(batch) > 0 {
				if err := s.luckyFiveRepo.StoreBatch(ctx, batch); err != nil {
					s.logger.Error(fmt.Errorf("worker %d failed to store final batch: %w", workerNum, err))
					return err
				}
			}
			return nil
		})
		if err != nil {
			s.logger.Error(fmt.Errorf("failed to add task to parallelizer group for worker %d: %w", workerNum, err))
		}
	}

	if err := group.Wait(); err != nil {
		err = fmt.Errorf("error during LuckyFive combination generation: %w", err)
		s.logger.Error(err)
		return err
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

	poolSize := runtime.NumCPU()
	group := parallelizer.NewGroup(parallelizer.WithPoolSize(poolSize))
	defer group.Close()

	workPerWorker := count / poolSize
	extraWork := count % poolSize

	for i := 0; i < poolSize; i++ {
		workerNum := i
		jobs := workPerWorker
		if workerNum < extraWork {
			jobs++
		}
		if jobs == 0 {
			continue
		}

		err := group.Add(func() error {
			source := rand.NewSource(time.Now().UnixNano() + int64(workerNum))
			random := rand.New(source)
			batch := make([]entity.LuckySix, 0, batchSize)

			for j := 0; j < jobs; j++ {
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
						s.logger.Error(fmt.Errorf("worker %d failed to store LuckySix batch: %w", workerNum, err))
						return err
					}
					batch = make([]entity.LuckySix, 0, batchSize)
				}
			}

			if len(batch) > 0 {
				if err := s.luckySixRepo.StoreBatch(ctx, batch); err != nil {
					s.logger.Error(fmt.Errorf("worker %d failed to store final LuckySix batch: %w", workerNum, err))
					return err
				}
			}
			return nil
		})
		if err != nil {
			s.logger.Error(fmt.Errorf("failed to add LuckySix task to parallelizer group for worker %d: %w", workerNum, err))
		}
	}

	if err := group.Wait(); err != nil {
		err = fmt.Errorf("error during LuckySix combination generation: %w", err)
		s.logger.Error(err)
		return err
	}

	s.logger.Info("Finished generation of LuckySix combinations")
	return nil
}
