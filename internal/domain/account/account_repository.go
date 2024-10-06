package account_domain

import "context"

type IAccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id string) (*Account, error)
	Delete(ctx context.Context, id string) error
}
