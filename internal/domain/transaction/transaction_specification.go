package transaction

import (
	"errors"
)

// Operation types
const (
	Deposit    = "DEPOSIT"
	Withdrawal = "WITHDRAWAL"
	Transfer   = "TRANSFER"
)

const (
	ListTransactionsLimit = 100
)

var (
	ErrUnsupportedType = errors.New("unsupported transaction type")
)

func validOperationType(operationType string) error {
	var validOperationTypes = []string{
		Deposit,
		Withdrawal,
		Transfer,
	}
	for _, validType := range validOperationTypes {
		if operationType == validType {
			return nil
		}
	}
	return ErrUnsupportedType
}
