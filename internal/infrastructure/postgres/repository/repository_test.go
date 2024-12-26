package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/config"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type testStruct struct {
	ID   int
	Name string
}

func TestExecDB(t *testing.T) {
	repo, mock, ctx, bunDB := PrepareTestRepository(t, func(db *bun.DB) *repository.Repository[testStruct] {
		return repository.NewRepository[testStruct](db)
	})

	t.Run("トランザクションなしの場合", func(t *testing.T) {
		execDB := repo.ExecDB(ctx)
		assert.Equal(t, bunDB, execDB)
	})

	t.Run("トランザクションありの場合", func(t *testing.T) {
		mock.ExpectBegin()

		tx, err := bunDB.BeginTx(ctx, nil)
		assert.NoError(t, err)

		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()

		txCtx := context.WithValue(ctx, config.CtxTransactionKey(), tx)
		execDB := repo.ExecDB(txCtx)
		assert.Equal(t, tx, execDB)
	})

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// テスト用のリポジトリを作成するためのヘルパー関数です。
func PrepareTestRepository[T any](t *testing.T, newRepo func(db *bun.DB) T) (T, sqlmock.Sqlmock, context.Context, *bun.DB) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	bunDB := bun.NewDB(db, pgdialect.New())

	repo := newRepo(bunDB)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err)
	})

	return repo, mock, ctx, bunDB
}
