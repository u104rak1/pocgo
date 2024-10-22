package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var (
	ErrDB = errors.New("database error")
)

func PrepareTestRepository[T any](t *testing.T, newRepo func(db *bun.DB) T) (T, sqlmock.Sqlmock, context.Context, *bun.DB) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	bunDB := bun.NewDB(db, pgdialect.New())

	repo := newRepo(bunDB)

	t.Cleanup(func() {
		db.Close()
	})

	return repo, mock, ctx, bunDB
}
