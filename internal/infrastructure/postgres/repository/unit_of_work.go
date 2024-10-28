package repository

import (
	"context"

	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	"github.com/ucho456job/pocgo/internal/config"
	"github.com/uptrace/bun"
)

type unitOfWork struct {
	db *bun.DB
}

func NewUnitOfWork(db *bun.DB) unitofwork.IUnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (u *unitOfWork) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctxWithTx := setTx(ctx, tx)

	err = f(ctxWithTx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

type unitOfWorkWithResult[T any] struct {
	db *bun.DB
}

func NewUnitOfWorkWithResult[T any](db *bun.DB) unitofwork.IUnitOfWorkWithResult[T] {
	return &unitOfWorkWithResult[T]{
		db: db,
	}
}

func (u *unitOfWorkWithResult[T]) RunInTx(ctx context.Context, f func(ctx context.Context) (*T, error)) (*T, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	ctxWithTx := setTx(ctx, tx)

	result, err := f(ctxWithTx)
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

func setTx(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, config.CtxTransactionKey(), tx)
}

func getTx(ctx context.Context) bun.IDB {
	tx, ok := ctx.Value(config.CtxTransactionKey()).(bun.IDB)
	if !ok {
		return nil
	}
	return tx
}
