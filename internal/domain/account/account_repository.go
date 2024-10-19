package account

// mockgen -source=internal/domain/account/account_repository.go -destination=internal/domain/account/mock/mock_account_repository.go -package=mock

import "context"

type IAccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id string) (*Account, error)
	ListByUserID(ctx context.Context, userID string) ([]*Account, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	Delete(ctx context.Context, id string) error
}
