package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// テスト用のリポジトリを作成するためのヘルパー関数です。
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
