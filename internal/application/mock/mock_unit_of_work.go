package mock

import "context"

type MockIUnitOfWork struct {
	Func func(ctx context.Context) error
}

func (m *MockIUnitOfWork) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

type MockIUnitOfWorkWithResult[T any] struct {
	Func func(ctx context.Context) (*T, error)
}

func (m *MockIUnitOfWorkWithResult[T]) RunInTx(ctx context.Context, f func(ctx context.Context) (*T, error)) (*T, error) {
	return f(ctx)
}
