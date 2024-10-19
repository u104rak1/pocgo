package transaction

// mockgen -source=internal/domain/transaction/transaction_repository.go -destination=internal/domain/transaction/mock/mock_transaction_repository.go -package=mock

type ITransactionRepository interface {
	Save(transaction *Transaction) error
	ListByAccountID(accountID string, limit, offset *int) ([]*Transaction, error)
}
