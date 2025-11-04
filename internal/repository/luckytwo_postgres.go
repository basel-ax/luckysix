package repository

import (
	"context"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/db"
)

// LuckyTwoRepo -.
type LuckyTwoRepo struct {
	*db.Gorm
}

// NewLuckyTwoRepo -.
func NewLuckyTwoRepo(g *db.Gorm) *LuckyTwoRepo {
	return &LuckyTwoRepo{g}
}

// GetAll -.
func (r *LuckyTwoRepo) GetAll(ctx context.Context) ([]*entity.Luckytwo, error) {
	var luckytwos []*entity.Luckytwo
	err := r.WithContext(ctx).Find(&luckytwos).Error
	return luckytwos, err
}

// GetByID -.
func (r *LuckyTwoRepo) GetByID(ctx context.Context, id uint) (*entity.Luckytwo, error) {
	var luckytwo entity.Luckytwo
	err := r.WithContext(ctx).First(&luckytwo, id).Error
	return &luckytwo, err
}
