package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/1boombacks1/testTaskInfotecs/internal/models"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo/repoerrs"
	"github.com/1boombacks1/testTaskInfotecs/pkg/postgres"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

type OperationsRepo struct {
	*postgres.Postgres
}

func NewOperationsRepo(pg *postgres.Postgres) *OperationsRepo {
	return &OperationsRepo{pg}
}

func (r *OperationsRepo) Create(ctx context.Context, from, to uuid.UUID, amount float32) (int, error) {
	sql, args, _ := r.Builder.
		Insert("operations").
		Columns("from_wallet_id", "to_wallet_id", "amount").
		Values(from, to, amount).
		Suffix("RETURNING id").
		ToSql()
	var id int
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, fmt.Errorf("OperationRepo.Create - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *OperationsRepo) GetAllByWalletID(ctx context.Context, id uuid.UUID) ([]models.Operations, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("operations").
		Where("from_wallet_id = $1 or to_wallet_id = $1", id).
		ToSql()
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("OperationsRepo.GetAllByWalletID - r.PoolQuery: %v", err)
	}

	operations, err := pgx.CollectRows[models.Operations](rows, pgx.RowToStructByName[models.Operations])
	if err != nil {
		return nil, fmt.Errorf("OperationsRepo.GetAllByWalletID - pgx.CollectRows: %v", err)
	}

	return operations, nil
}
