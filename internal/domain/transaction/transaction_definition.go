package transaction_domain

import "errors"

var (
	ErrInvalidTransactionID     = errors.New("invalid transaction id")
	ErrUnsupportTransactionType = errors.New("unsupported transaction type")
)

const (
	TransactionDeposit  = "DEPOSIT"
	TransactionWithdraw = "WITHDRAW"
	TransactionTransfer = "TRANSFER"
)
