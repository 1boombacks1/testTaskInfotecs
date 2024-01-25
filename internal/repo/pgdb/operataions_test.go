package pgdb

import (
	"context"
	"errors"
	"testing"

	uuidgen "github.com/1boombacks1/testTaskInfotecs/pkg/UUIDGen"
	"github.com/1boombacks1/testTaskInfotecs/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestOperations_Create_Success(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	pg := &postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}
	repo := NewOperationsRepo(pg)

	idGen := uuidgen.DefaultGenerator{}
	from, _ := idGen.New()
	to, _ := idGen.New()
	var amount float32 = 30.0
	var expectedResult = 52

	pool.ExpectQuery("INSERT INTO operations \\(from_wallet_id,to_wallet_id,amount\\) VALUES \\(\\$1,\\$2,\\$3\\) RETURNING id").
		WithArgs(from, to, amount).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedResult))

	id, err := repo.Create(context.Background(), from, to, amount)

	assert.Equal(t, expectedResult, id)
	assert.NoError(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestOperations_Create_DBError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	pg := &postgres.Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}
	repo := NewOperationsRepo(pg)

	idGen := uuidgen.DefaultGenerator{}
	from, _ := idGen.New()
	to, _ := idGen.New()
	var amount float32 = 30.0

	pool.ExpectQuery("INSERT INTO operations \\(from_wallet_id,to_wallet_id,amount\\) VALUES \\(\\$1,\\$2,\\$3\\) RETURNING id").
		WithArgs(from, to, amount).
		WillReturnError(errors.New("Test Error"))

	_, err = repo.Create(context.Background(), from, to, amount)

	assert.Error(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}
