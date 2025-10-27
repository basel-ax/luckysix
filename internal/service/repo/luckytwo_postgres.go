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

// Count returns the total number of records in the luckytwos table.
func (r *LuckyTwoRepo) Count(ctx context.Context) (int64, error) {
	sql, _, err := r.Builder.
		Select("COUNT(*)").
		From("luckytwos").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("LuckyTwoRepo - Count - r.Builder: %w", err)
	}

	var count int64
	err = r.Pool.QueryRow(ctx, sql).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("LuckyTwoRepo - Count - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}

// GetByIDs fetches multiple Luckytwo entities by their primary keys.
func (r *LuckyTwoRepo) GetByIDs(ctx context.Context, ids []uint) (map[uint]entity.Luckytwo, error) {
	if len(ids) == 0 {
		return make(map[uint]entity.Luckytwo), nil
	}

	sql, args, err := r.Builder.
		Select("id", "word_one", "word_two").
		From("luckytwos").
		Where("id = ANY(?)", ids).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LuckyTwoRepo - GetByIDs - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("LuckyTwoRepo - GetByIDs - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	results := make(map[uint]entity.Luckytwo)
	for rows.Next() {
		var lt entity.Luckytwo
		if err := rows.Scan(&lt.ID, &lt.WordOne, &lt.WordTwo); err != nil {
			return nil, fmt.Errorf("LuckyTwoRepo - GetByIDs - rows.Scan: %w", err)
		}
		results[lt.ID] = lt
	}

	return results, nil
}
