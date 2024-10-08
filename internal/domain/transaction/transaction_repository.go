package transaction_domain

type ITransactionRepository interface {
	Save(transaction *Transaction) error
	ListByAccountID(accountID string, limit, offset int) ([]*Transaction, error)
}
