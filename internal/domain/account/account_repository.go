package account

import "context"

type IAccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id string) (*Account, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
}
