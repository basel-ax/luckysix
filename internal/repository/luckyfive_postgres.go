package repository

import (
	"context"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/db"
)

// LuckyFiveRepo -.
type LuckyFiveRepo struct {
	*db.Gorm
}

// NewLuckyFiveRepo -.
func NewLuckyFiveRepo(g *db.Gorm) *LuckyFiveRepo {
	return &LuckyFiveRepo{g}
}

// GetAll -.
func (r *LuckyFiveRepo) GetAll(ctx context.Context) ([]*entity.LuckyFive, error) {
	var luckyfives []*entity.LuckyFive
	err := r.WithContext(ctx).Find(&luckyfives).Error
	return luckyfives, err
}

// GetByID -.
func (r *LuckyFiveRepo) GetByID(ctx context.Context, id uint) (*entity.LuckyFive, error) {
	var luckyfive entity.LuckyFive
	err := r.WithContext(ctx).First(&luckyfive, id).Error
	return &luckyfive, err
}
