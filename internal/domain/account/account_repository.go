package account

import "context"

type IAccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id string) (*Account, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*Account, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	Delete(ctx context.Context, id string) error
}
