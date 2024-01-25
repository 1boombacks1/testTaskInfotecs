package pgdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/1boombacks1/testTaskInfotecs/internal/models"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo/repoerrs"
	"github.com/1boombacks1/testTaskInfotecs/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type WalletRepo struct {
	*postgres.Postgres
}

func NewWalletRepo(pg *postgres.Postgres) *WalletRepo {
	return &WalletRepo{
		Postgres: pg,
	}
}

func (r *WalletRepo) Create(ctx context.Context, id uuid.UUID) error {
	sql, args, _ := r.Builder.
		Insert("wallets").
		Columns("id").
		Values(id).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("WalletRepo.CreateWallet - r.Pool.Exec: %v", err)
	}

	return nil
}

func (r *WalletRepo) GetByID(ctx context.Context, id uuid.UUID) (models.Wallet, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("wallets").
		Where("id = ?", id).
		ToSql()

	var wallet models.Wallet
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Wallet{}, repoerrs.ErrNotFound
		}
		return models.Wallet{}, fmt.Errorf("WalletRepo.GetWalletByID - r.Pool.QueryRow: %v", err)
	}

	return wallet, nil
}

func (r *WalletRepo) Transfer(ctx context.Context, from uuid.UUID, to uuid.UUID, amount float32) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("WalletRepo.Transfer - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.
		Select("balance").
		From("wallets").
		Where("id = ?", from).
		ToSql()
	var fromBalance float32
	if err := tx.QueryRow(ctx, sql, args...).Scan(&fromBalance); err != nil {
		return fmt.Errorf("WalletRepo.Transfer - tx.QueryRow: %v", err)
	}

	if fromBalance < amount {
		return repoerrs.ErrNotEnoughBalance
	}

	sql, args, _ = r.Builder.
		Update("wallets").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Set("updated_at", time.Now()).
		Where("id = ?", from).
		ToSql()
	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("WalletRepo.Transfer - (Update 'From') tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Update("wallets").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Set("updated_at", time.Now()).
		Where("id = ?", to).
		ToSql()
	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("WalletRepo.Transfer - (Update 'To') tx.Exec: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("WalletRepo.Transfer - tx.Commit: %v", err)
	}

	return nil
}
