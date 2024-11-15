package transaction

type TransactionDTO struct {
	ID                string
	AccountID         string
	ReceiverAccountID *string
	OperationType     string
	Amount            float64
	Currency          string
	TransactionAt     string
}
