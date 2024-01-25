package pgdb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/1boombacks1/testTaskInfotecs/internal/models"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo/repoerrs"
	"github.com/1boombacks1/testTaskInfotecs/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

type MockIDGen struct{}

func (*MockIDGen) New() (uuid.UUID, error) {
	return uuid.UUID{}, errors.New("failed to create a uuid")
}
func (*MockIDGen) FromString(text string) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (*MockIDGen) NewV4() (uuid.UUID, error) {
	return uuid.NewV4()
}

func TestWallet_Create_Success(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	pg := &postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}
	idgen := MockIDGen{}
	id, _ := idgen.NewV4()

	pool.ExpectExec("INSERT INTO wallets").WithArgs(id).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	repo := NewWalletRepo(pg)
	err = repo.Create(context.Background(), id)

	assert.NoError(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestWallet_Create_DBError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	pg := &postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}
	idgen := MockIDGen{}
	id, _ := idgen.NewV4()

	pool.ExpectExec("INSERT INTO wallets").WithArgs(id).WillReturnError(errors.New("test error"))

	repo := NewWalletRepo(pg)
	err = repo.Create(context.Background(), id)

	assert.Error(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestWallet_GetByID_Success(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	pg := postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}

	repo := NewWalletRepo(&pg)

	idgen := MockIDGen{}
	id, _ := idgen.NewV4()
	expectedResult := models.Wallet{ID: id, Balance: 42, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	rows := pgxmock.NewRows([]string{"id", "balance", "created_at", "updated_at"})
	rows.AddRow(expectedResult.ID, expectedResult.Balance, expectedResult.CreatedAt, expectedResult.UpdatedAt)

	pool.ExpectQuery("SELECT \\* FROM wallets WHERE id = \\$1").WithArgs(id).WillReturnRows(rows)

	result, err := repo.GetByID(context.Background(), id)

	assert.Equal(t, expectedResult, result)
	assert.NoError(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestWallet_GetByID_NotFoundError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	pg := postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}

	repo := NewWalletRepo(&pg)

	idgen := MockIDGen{}
	id, _ := idgen.NewV4()
	rows := pgxmock.NewRows([]string{"id", "balance", "created_at", "updated_at"})
	pool.ExpectQuery("SELECT \\* FROM wallets WHERE id = \\$1").WithArgs(id).WillReturnRows(rows)

	_, err = repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoerrs.ErrNotFound)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestWallet_GetByID_DBError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	pg := &postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}

	repo := NewWalletRepo(pg)

	idgen := MockIDGen{}
	id, _ := idgen.NewV4()
	pool.ExpectQuery("SELECT \\* FROM wallets WHERE id = \\$1").WithArgs(id).WillReturnError(errors.New("Test error"))

	_, err = repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.NotErrorIs(t, err, repoerrs.ErrNotFound)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestWallet_Transfer_Success(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	pg := &postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}
	repo := NewWalletRepo(pg)

	idgen := MockIDGen{}
	from, _ := idgen.NewV4()
	to, _ := idgen.NewV4()
	var amount float32 = 30.0
	var fromBalance float32 = 100.0

	rows := pgxmock.NewRows([]string{"balance"})
	rows.AddRow(fromBalance)

	pool.ExpectBegin()
	pool.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").WithArgs(from).WillReturnRows(rows)
	pool.ExpectExec("UPDATE wallets").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), from).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectExec("UPDATE wallets").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), to).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectCommit()

	err = repo.Transfer(context.Background(), from, to, float32(amount))

	assert.NoError(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestWallet_Transfer_RollbackOnFailure(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	pg := &postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}
	repo := NewWalletRepo(pg)
	idgen := MockIDGen{}
	from, _ := idgen.NewV4()
	to, _ := idgen.NewV4()
	var amount float32 = 30.0
	var fromBalance float32 = 100.0

	rows := pgxmock.NewRows([]string{"balance"})
	rows.AddRow(fromBalance)

	pool.ExpectBegin()
	pool.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").WithArgs(from).WillReturnRows(rows)
	pool.ExpectExec("UPDATE wallets").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), from).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectExec("UPDATE wallets").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), to).WillReturnError(errors.New("test error"))
	pool.ExpectRollback()

	err = repo.Transfer(context.Background(), from, to, float32(amount))

	assert.Error(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}
