package unitofwork

import "context"

type IUnitOfWork interface {
	RunInTx(ctx context.Context, f func(ctx context.Context, tx ITransaction) error) error
}

type IUnitOfWorkWithResult[T any] interface {
	RunInTx(ctx context.Context, f func(ctx context.Context, tx ITransaction) (*T, error)) (*T, error)
}

type ITransaction interface {
	Commit() error
	Rollback() error
}
