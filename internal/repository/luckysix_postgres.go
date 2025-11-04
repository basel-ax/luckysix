package repository

import (
	"context"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/db"
)

// LuckySixRepo -.
type LuckySixRepo struct {
	*db.Gorm
}

// NewLuckySixRepo -.
func NewLuckySixRepo(g *db.Gorm) *LuckySixRepo {
	return &LuckySixRepo{g}
}

// GetAll -.
func (r *LuckySixRepo) GetAll(ctx context.Context) ([]*entity.LuckySix, error) {
	var luckysixes []*entity.LuckySix
	err := r.WithContext(ctx).Find(&luckysixes).Error
	return luckysixes, err
}

// GetByID -.
func (r *LuckySixRepo) GetByID(ctx context.Context, id uint) (*entity.LuckySix, error) {
	var luckysix entity.LuckySix
	err := r.WithContext(ctx).First(&luckysix, id).Error
	return &luckysix, err
}
