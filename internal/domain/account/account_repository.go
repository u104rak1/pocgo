package account

import (
	"context"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
)

type IAccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id AccountID) (*Account, error)
	CountByUserID(ctx context.Context, userID userDomain.UserID) (int, error)
}
