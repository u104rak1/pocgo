package unitofwork

import "context"

type IUnitOfWork interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
}

type IUnitOfWorkWithResult[T any] interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) (*T, error)) (*T, error)
}
