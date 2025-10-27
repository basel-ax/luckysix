package repo

import (
	"context"
	"fmt"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/postgres"
)

// LuckySixRepo -.
type LuckySixRepo struct {
	*postgres.Postgres
}

// NewLuckySixRepo -.
func NewLuckySixRepo(pg *postgres.Postgres) *LuckySixRepo {
	return &LuckySixRepo{pg}
}

// StoreBatch stores a batch of LuckySix entities, ignoring duplicates.
func (r *LuckySixRepo) StoreBatch(ctx context.Context, luckySixes []entity.LuckySix) error {
	if len(luckySixes) == 0 {
		return nil
	}

	builder := r.Builder.Insert("lucky_sixes").Columns("lucky_five_id", "lucky_two_id", "created_at", "updated_at")

	for _, ls := range luckySixes {
		builder = builder.Values(ls.LuckyFiveID, ls.LuckyTwoID, ls.CreatedAt, ls.UpdatedAt)
	}

	// This suffix handles the "without repeat" requirement at the database level.
	builder = builder.Suffix("ON CONFLICT (lucky_five_id, lucky_two_id) DO NOTHING")

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("LuckySixRepo - StoreBatch - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("LuckySixRepo - StoreBatch - r.Pool.Exec: %w", err)
	}

	return nil
}
