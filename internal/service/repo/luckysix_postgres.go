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

// GetBatchStartingFromID fetches a batch of LuckySix records with IDs greater than the startID.
func (r *LuckySixRepo) GetBatchStartingFromID(ctx context.Context, startID uint, limit int) ([]entity.LuckySix, error) {
	sql, args, err := r.Builder.
		Select("id", "lucky_five_id", "lucky_two_id").
		From("lucky_sixes").
		Where("id > ?", startID).
		OrderBy("id ASC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LuckySixRepo - GetBatchStartingFromID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("LuckySixRepo - GetBatchStartingFromID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var results []entity.LuckySix
	for rows.Next() {
		var ls entity.LuckySix
		if err := rows.Scan(&ls.ID, &ls.LuckyFiveID, &ls.LuckyTwoID); err != nil {
			return nil, fmt.Errorf("LuckySixRepo - GetBatchStartingFromID - rows.Scan: %w", err)
		}
		results = append(results, ls)
	}

	return results, nil
}
