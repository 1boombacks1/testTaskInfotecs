package repo

import (
	"context"

	"github.com/1boombacks1/testTaskInfotecs/internal/models"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo/pgdb"
	"github.com/1boombacks1/testTaskInfotecs/pkg/postgres"
	"github.com/gofrs/uuid/v5"
)

type Wallet interface {
	Create(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (models.Wallet, error)
	Transfer(ctx context.Context, from uuid.UUID, to uuid.UUID, amount float32) error
}
type Operations interface {
	Create(ctx context.Context, from, to uuid.UUID, amount float32) (int, error)
	GetAllByWalletID(ctx context.Context, id uuid.UUID) ([]models.Operations, error)
}

type Repos struct {
	Wallet
	Operations
}

func NewRepos(pg *postgres.Postgres) *Repos {
	return &Repos{
		Wallet:     pgdb.NewWalletRepo(pg),
		Operations: pgdb.NewOperationsRepo(pg),
	}
}
