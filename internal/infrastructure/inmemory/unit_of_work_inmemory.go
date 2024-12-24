package inmemory

import (
	"context"

	unitofwork "github.com/u104rak1/pocgo/internal/application/unit_of_work"
)

type unitOfWorkInMemory struct{}

func NewUnitOfWorkInMemory() unitofwork.IUnitOfWork {
	return &unitOfWorkInMemory{}
}

func (u *unitOfWorkInMemory) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

type unitOfWorkInMemoryWithResult[T any] struct{}

func NewUnitOfWorkInMemoryWithResult[T any]() unitofwork.IUnitOfWorkWithResult[T] {
	return &unitOfWorkInMemoryWithResult[T]{}
}

func (u *unitOfWorkInMemoryWithResult[T]) RunInTx(ctx context.Context, f func(ctx context.Context) (*T, error)) (*T, error) {
	return f(ctx)
}
