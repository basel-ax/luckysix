package repository

import (
	"context"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/db"
)

// WalletBalanceRepo -.
type WalletBalanceRepo struct {
	*db.Gorm
}

// NewWalletBalanceRepo -.
func NewWalletBalanceRepo(g *db.Gorm) *WalletBalanceRepo {
	return &WalletBalanceRepo{g}
}

// Create -.
func (r *WalletBalanceRepo) Create(ctx context.Context, wb *entity.WalletBalance) error {
	return r.WithContext(ctx).Create(wb).Error
}

// Update -.
func (r *WalletBalanceRepo) Update(ctx context.Context, wb *entity.WalletBalance) error {
	return r.WithContext(ctx).Save(wb).Error
}

// GetByLuckySixID -.
func (r *WalletBalanceRepo) GetByLuckySixID(ctx context.Context, luckySixID uint) (*entity.WalletBalance, error) {
	var wb entity.WalletBalance
	err := r.WithContext(ctx).Where("lucky_six_id = ?", luckySixID).First(&wb).Error
	return &wb, err
}
