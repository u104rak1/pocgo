package transaction

import (
	"context"
	"time"
)

type ListTransactionsParams struct {
	AccountID      string
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
