// Package service implements application business logic. Each logic group in own file.
package service

import (
	"context"

	entity "github.com/basel-ax/luckysix/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=service_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}

	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// Tasks -.
	Tasks interface {
		CheckRabbitTask() string
	}

	// PoolRepo -.
	PoolRepo interface {
	}

	// Bip39 -.
	Bip39 interface {
		GenerateAndStoreLuckyTwoCombinations(context.Context) error
		GenerateAndStoreLuckyFiveCombinations(ctx context.Context, count int) error
		GenerateAndStoreLuckySixCombinations(ctx context.Context, count int) error
	}

	// LuckyTwoRepo -.
	LuckyTwoRepo interface {
		StoreBatch(context.Context, []entity.Luckytwo) error
		Count(ctx context.Context) (int64, error)
		GetByIDs(ctx context.Context, ids []uint) (map[uint]entity.Luckytwo, error)
	}

	// LuckyFiveRepo -.
	LuckyFiveRepo interface {
		StoreBatch(ctx context.Context, luckyFives []entity.LuckyFive) error
		Count(ctx context.Context) (int64, error)
		GetByID(ctx context.Context, id uint) (entity.LuckyFive, error)
	}

	// LuckySixRepo -.
	LuckySixRepo interface {
		StoreBatch(ctx context.Context, luckySixes []entity.LuckySix) error
		GetBatchStartingFromID(ctx context.Context, startID uint, limit int) ([]entity.LuckySix, error)
	}

	// Blockchain -.
	Blockchain interface {
		ProcessLuckySixCombinations(ctx context.Context, count int) error
	}

	// WalletBalanceRepo -.
	WalletBalanceRepo interface {
		StoreBatch(ctx context.Context, balances []entity.WalletBalance) error
		GetLastProcessedLuckySixID(ctx context.Context) (uint, error)
	}
)
