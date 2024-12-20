package account

import (
	"context"

	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IAccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id idVO.AccountID) (*Account, error)
	CountByUserID(ctx context.Context, userID idVO.UserID) (int, error)
}
