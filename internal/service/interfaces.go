// Package service implements services for working with business logic.
package service

import (
	"context"

	"github.com/basel-ax/luckysix/internal/entity"
)

type (
	// LuckySix -.
	LuckySix interface {
		GenerateAndProcess(ctx context.Context)
	}
	// Wallet -.
	Wallet interface {
		GenerateMnemonic(ctx context.Context, luckySix *entity.LuckySix) (string, error)
		GenerateWallet(ctx context.Context, mnemonic string) (string, string, error)
		CheckBalance(ctx context.Context, address string, blockchains []*entity.Blockchain) (string, error)
	}

	// Notification -.
	Notification interface {
		Notify(ctx context.Context, message string) error
	}
)
