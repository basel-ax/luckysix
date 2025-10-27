package repo

import (
	"context"
	"fmt"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/postgres"
)

// LuckyTwoRepo -.
type LuckyTwoRepo struct {
	*postgres.Postgres
}

// NewLuckyTwoRepo -.
func NewLuckyTwoRepo(pg *postgres.Postgres) *LuckyTwoRepo {
	return &LuckyTwoRepo{pg}
}

// StoreBatch stores a batch of Luckytwo entities.
func (r *LuckyTwoRepo) StoreBatch(ctx context.Context, luckyTwos []entity.Luckytwo) error {
	if len(luckyTwos) == 0 {
		return nil
	}

	builder := r.Builder.Insert("luckytwos").Columns("word_one", "word_two")

	for _, lt := range luckyTwos {
		builder = builder.Values(lt.WordOne, lt.WordTwo)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("LuckyTwoRepo - StoreBatch - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("LuckyTwoRepo - StoreBatch - r.Pool.Exec: %w", err)
	}

	return nil
}
