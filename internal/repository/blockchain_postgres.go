package repository

import (
	"context"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/db"
)

// BlockchainRepo -.
type BlockchainRepo struct {
	*db.Gorm
}

// NewBlockchainRepo -.

func NewBlockchainRepo(g *db.Gorm) *BlockchainRepo {
	return &BlockchainRepo{g}
}

// GetEnabledBlockchains -.
func (r *BlockchainRepo) GetEnabledBlockchains(ctx context.Context) ([]*entity.Blockchain, error) {

	var blockchains []*entity.Blockchain

	err := r.WithContext(ctx).Where("enabled = ?", true).Find(&blockchains).Error
	return blockchains, err
}
