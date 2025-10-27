package repo

import (
	"context"
	"fmt"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/postgres"
)

// LuckyFiveRepo -.
type LuckyFiveRepo struct {
	*postgres.Postgres
}

// NewLuckyFiveRepo -.
func NewLuckyFiveRepo(pg *postgres.Postgres) *LuckyFiveRepo {
	return &LuckyFiveRepo{pg}
}

// StoreBatch stores a batch of LuckyFive entities.
func (r *LuckyFiveRepo) StoreBatch(ctx context.Context, luckyFives []entity.LuckyFive) error {
	if len(luckyFives) == 0 {
		return nil
	}

	columns := []string{"pair_one", "pair_two", "pair_three", "pair_four", "pair_five", "created_at", "updated_at"}
	builder := r.Builder.Insert("luckyfives").Columns(columns...)

	for _, lf := range luckyFives {
		builder = builder.Values(lf.PairOne, lf.PairTwo, lf.PairThree, lf.PairFour, lf.PairFive, lf.CreatedAt, lf.UpdatedAt)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("LuckyFiveRepo - StoreBatch - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("LuckyFiveRepo - StoreBatch - r.Pool.Exec: %w", err)
	}

	return nil
}

// Count returns the total number of records in the luckyfives table.
func (r *LuckyFiveRepo) Count(ctx context.Context) (int64, error) {
	sql, _, err := r.Builder.
		Select("COUNT(*)").
		From("luckyfives").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("LuckyFiveRepo - Count - r.Builder: %w", err)
	}

	var count int64
	err = r.Pool.QueryRow(ctx, sql).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("LuckyFiveRepo - Count - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}
