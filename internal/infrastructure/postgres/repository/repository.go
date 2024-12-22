package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository[T any] struct {
	db *bun.DB
}

func NewRepository[T any](db *bun.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) ExecDB(ctx context.Context) bun.IDB {
	tx := getTx(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}
