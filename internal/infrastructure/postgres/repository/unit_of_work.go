package repository

import (
	"context"

	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	"github.com/uptrace/bun"
)

type IUnitOfWork interface {
	RunInTx(ctx context.Context, f func(ctx context.Context, tx unitofwork.ITransaction) error) error
}

type unitOfWork struct {
	db *bun.DB
}

func NewUnitOfWork(db *bun.DB) IUnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (u *unitOfWork) RunInTx(ctx context.Context, f func(ctx context.Context, tx unitofwork.ITransaction) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	bunTx := &transaction{tx: tx}

	err = f(ctx, bunTx)
	if err != nil {
		rollbackErr := bunTx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return bunTx.Commit()
}

type IUnitOfWorkWithResult[T any] interface {
	RunInTx(ctx context.Context, f func(ctx context.Context, tx unitofwork.ITransaction) (*T, error)) (*T, error)
}

type unitOfWorkWithResult[T any] struct {
	db *bun.DB
}

func NewUnitOfWorkWithResult[T any](db *bun.DB) IUnitOfWorkWithResult[T] {
	return &unitOfWorkWithResult[T]{
		db: db,
	}
}

func (u *unitOfWorkWithResult[T]) RunInTx(ctx context.Context, f func(ctx context.Context, tx unitofwork.ITransaction) (*T, error)) (*T, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	bunTx := &transaction{tx: tx}

	result, err := f(ctx, bunTx)
	if err != nil {
		rollbackErr := bunTx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	err = bunTx.Commit()
	return result, err
}

type transaction struct {
	tx bun.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}
