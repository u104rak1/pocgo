package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestUnitOfWork_RunInTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	bunDB := bun.NewDB(db, pgdialect.New())
	uow := repository.NewUnitOfWork(bunDB)

	t.Run("Positive: トランザクションが正常にコミットされる", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		err := uow.RunInTx(context.Background(), func(ctx context.Context) error {
			return nil
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Negative: エラーが発生した場合、ロールバックされる", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		err := uow.RunInTx(context.Background(), func(ctx context.Context) error {
			return assert.AnError
		})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Negative: トランザクションの開始に失敗する", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(assert.AnError)

		err := uow.RunInTx(context.Background(), func(ctx context.Context) error {
			return nil
		})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUnitOfWorkWithResult_RunInTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	bunDB := bun.NewDB(db, pgdialect.New())
	uow := repository.NewUnitOfWorkWithResult[string](bunDB)

	t.Run("Positive: トランザクションが正常にコミットされる", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		expected := "success"
		result, err := uow.RunInTx(context.Background(), func(ctx context.Context) (*string, error) {
			return &expected, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Negative: エラーが発生した場合、ロールバックされる", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		result, err := uow.RunInTx(context.Background(), func(ctx context.Context) (*string, error) {
			return nil, assert.AnError
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Negative: トランザクションの開始に失敗する", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(assert.AnError)

		result, err := uow.RunInTx(context.Background(), func(ctx context.Context) (*string, error) {
			return nil, nil
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
