// Package repository implements repository for working with db
package repository

import (
	"context"

	"github.com/basel-ax/luckysix/internal/entity"
)

type (
	// Blockchain -.
	Blockchain interface {
		GetEnabledBlockchains(context.Context) ([]*entity.Blockchain, error)
	}
	// LuckyTwo -.
	LuckyTwo interface {
		GetAll(context.Context) ([]*entity.Luckytwo, error)
		GetByID(ctx context.Context, id uint) (*entity.Luckytwo, error)
	}

	// LuckyFive -.
	LuckyFive interface {
		GetAll(context.Context) ([]*entity.LuckyFive, error)
		GetByID(ctx context.Context, id uint) (*entity.LuckyFive, error)
	}

	// LuckySix -.
	LuckySix interface {
		GetAll(context.Context) ([]*entity.LuckySix, error)
		GetByID(ctx context.Context, id uint) (*entity.LuckySix, error)
	}
	// WalletBalance -.
	WalletBalance interface {
		Create(ctx context.Context, wb *entity.WalletBalance) error
		Update(ctx context.Context, wb *entity.WalletBalance) error
		GetByLuckySixID(ctx context.Context, luckySixID uint) (*entity.WalletBalance, error)
	}
)
