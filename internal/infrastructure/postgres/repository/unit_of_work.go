package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type IUnitOfWork interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
}

type unitOfWork struct {
	db *bun.DB
}

func NewUnitOfWork(db *bun.DB) IUnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (u *unitOfWork) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = f(ctx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

type IUnitOfWorkWithResult[T any] interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) (*T, error)) (*T, error)
}

type unitOfWorkWithResult[T any] struct {
	db *bun.DB
}

func NewUnitOfWorkWithResult[T any](db *bun.DB) IUnitOfWorkWithResult[T] {
	return &unitOfWorkWithResult[T]{
		db: db,
	}
}

func (u *unitOfWorkWithResult[T]) RunInTx(ctx context.Context, f func(ctx context.Context) (*T, error)) (*T, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	result, err := f(ctx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	err = tx.Commit()
	return result, err
}
