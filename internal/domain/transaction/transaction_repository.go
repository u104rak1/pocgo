package transaction

import (
	"context"
	"time"

	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type ListTransactionsParams struct {
	AccountID      idVO.AccountID
	From           *time.Time
	To             *time.Time
	OperationTypes []string
	Sort           *string
	Limit          *int
	Page           *int
}

type ITransactionRepository interface {
	Save(ctx context.Context, transaction *Transaction) error
	ListWithTotalByAccountID(ctx context.Context, params ListTransactionsParams) (transactions []Transaction, total int, err error)
}
