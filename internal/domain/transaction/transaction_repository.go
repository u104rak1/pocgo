package transaction

import (
	"context"
	"time"

	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
)

type ListTransactionsParams struct {
	AccountID      accountDomain.AccountID
	From           *time.Time
	To             *time.Time
	OperationTypes []string
	Sort           *string
	Limit          *int
	Page           *int
}

type ITransactionRepository interface {
	Save(ctx context.Context, transaction *Transaction) error
	ListWithTotalByAccountID(ctx context.Context, params ListTransactionsParams) (transactions []*Transaction, total int, err error)
}
