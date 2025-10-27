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

	columns := []string{"lucky_six_id", "mnemonic", "address", "cosmos_address", "balance", "is_notified", "created_at", "updated_at"}
	builder := r.Builder.Insert("wallet_balances").Columns(columns...)

	for _, wb := range balances {
		builder = builder.Values(wb.LuckySixID, wb.Mnemonic, wb.Address, wb.CosmosAddress, wb.Balance, wb.IsNotified, wb.CreatedAt, wb.UpdatedAt)
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

// GetWalletsForBalanceCheck fetches wallets that have a zero balance or have never been checked.
func (r *WalletBalanceRepo) GetWalletsForBalanceCheck(ctx context.Context, limit int) ([]entity.WalletBalance, error) {
	sql, args, err := r.Builder.
		Select("id", "address", "cosmos_address").
		From("wallet_balances").
		Where("balance = '0' OR balance IS NULL").
		OrderBy("id ASC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletBalanceRepo - GetWalletsForBalanceCheck - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("WalletBalanceRepo - GetWalletsForBalanceCheck - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var results []entity.WalletBalance
	for rows.Next() {
		var wb entity.WalletBalance
		if err := rows.Scan(&wb.ID, &wb.Address, &wb.CosmosAddress); err != nil {
			return nil, fmt.Errorf("WalletBalanceRepo - GetWalletsForBalanceCheck - rows.Scan: %w", err)
		}
		results = append(results, wb)
	}

	return results, nil
}

// UpdateBalances performs a batch update of wallet balances.
func (r *WalletBalanceRepo) UpdateBalances(ctx context.Context, balances []entity.WalletBalance) error {
	if len(balances) == 0 {
		return nil
	}

	// Using a transaction for batch update
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("WalletBalanceRepo - UpdateBalances - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback is a no-op if the tx has been committed.

	for _, wb := range balances {
		sql, args, err := r.Builder.
			Update("wallet_balances").
			Set("balance", wb.Balance).
			Set("balance_updated_at", wb.BalanceUpdatedAt).
			Where("id = ?", wb.ID).
			ToSql()
		if err != nil {
			return fmt.Errorf("WalletBalanceRepo - UpdateBalances - r.Builder: %w", err)
		}
		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("WalletBalanceRepo - UpdateBalances - tx.Exec: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetPositiveBalanceWallets fetches wallets that have a balance > 0 and have not been notified yet.
func (r *WalletBalanceRepo) GetPositiveBalanceWallets(ctx context.Context, limit int) ([]entity.WalletBalance, error) {
	sql, args, err := r.Builder.
		Select("id", "address").
		From("wallet_balances").
		Where("balance > '0' AND is_notified = false").
		OrderBy("id ASC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletBalanceRepo - GetPositiveBalanceWallets - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("WalletBalanceRepo - GetPositiveBalanceWallets - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var results []entity.WalletBalance
	for rows.Next() {
		var wb entity.WalletBalance
		if err := rows.Scan(&wb.ID, &wb.Address); err != nil {
			return nil, fmt.Errorf("WalletBalanceRepo - GetPositiveBalanceWallets - rows.Scan: %w", err)
		}
		results = append(results, wb)
	}

	return results, nil
}

// MarkAsNotified updates the is_notified flag for a list of wallet IDs.
func (r *WalletBalanceRepo) MarkAsNotified(ctx context.Context, walletIDs []uint) error {
	if len(walletIDs) == 0 {
		return nil
	}

	sql, args, err := r.Builder.
		Update("wallet_balances").
		Set("is_notified", true).
		Where("id = ANY(?)", walletIDs).
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletBalanceRepo - MarkAsNotified - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletBalanceRepo - MarkAsNotified - r.Pool.Exec: %w", err)
	}

	return nil
}
