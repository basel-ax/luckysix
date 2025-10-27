package service

import (
	"context"
	"fmt"
	"runtime"

	"github.com/alexvec/go-bip39"
	"github.com/shomali11/parallelizer"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/logger"
)

const (
	batchSize = 1000
)

// Bip39Service -.
type Bip39Service struct {
	repo   LuckyTwoRepo
	logger logger.Interface
}

// NewBip39Service -.
func NewBip39Service(repo LuckyTwoRepo, l logger.Interface) *Bip39Service {
	return &Bip39Service{repo: repo, logger: l}
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
					if err := s.repo.StoreBatch(ctx, batch); err != nil {
						err = fmt.Errorf("failed to store batch for i=%d: %w", i, err)
						s.logger.Error(err)
						return err
					}
					batch = make([]entity.Luckytwo, 0, batchSize) // Reset batch
				}
			}

			if len(batch) > 0 { // Store remaining items
				if err := s.repo.StoreBatch(ctx, batch); err != nil {
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
