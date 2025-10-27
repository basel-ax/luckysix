package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

// WalletBalanceRepo -.
type WalletBalanceRepo struct {
	*postgres.Postgres
}

// NewWalletBalanceRepo -.
func NewWalletBalanceRepo(pg *postgres.Postgres) *WalletBalanceRepo {
	return &WalletBalanceRepo{pg}
}

// StoreBatch stores a batch of WalletBalance entities.
func (r *WalletBalanceRepo) StoreBatch(ctx context.Context, balances []entity.WalletBalance) error {
	if len(balances) == 0 {
		return nil
	}

	columns := []string{"lucky_six_id", "mnemonic", "address", "balance", "created_at", "updated_at"}
	builder := r.Builder.Insert("wallet_balances").Columns(columns...)

	for _, wb := range balances {
		builder = builder.Values(wb.LuckySixID, wb.Mnemonic, wb.Address, wb.Balance, wb.CreatedAt, wb.UpdatedAt)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("WalletBalanceRepo - StoreBatch - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletBalanceRepo - StoreBatch - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetLastProcessedLuckySixID finds the highest LuckySixID processed so far.
func (r *WalletBalanceRepo) GetLastProcessedLuckySixID(ctx context.Context) (uint, error) {
	sql, _, err := r.Builder.
		Select("MAX(lucky_six_id)").
		From("wallet_balances").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("WalletBalanceRepo - GetLastProcessedLuckySixID - r.Builder: %w", err)
	}

	var lastID *uint
	err = r.Pool.QueryRow(ctx, sql).Scan(&lastID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || lastID == nil {
			return 0, nil // No records yet, start from the beginning
		}
		return 0, fmt.Errorf("WalletBalanceRepo - GetLastProcessedLuckySixID - r.Pool.QueryRow: %w", err)
	}

	if lastID == nil {
		return 0, nil
	}

	return *lastID, nil
}
