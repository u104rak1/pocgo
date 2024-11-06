package transaction

import "context"

type ITransactionRepository interface {
	Save(ctx context.Context, transaction *Transaction) error
	ListByAccountID(ctx context.Context, accountID string, limit, offset *int) ([]*Transaction, error)
}
