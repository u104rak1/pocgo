package transaction

import (
	"errors"

	"github.com/ucho456job/pocgo/pkg/ulid"
)

// Operation types
const (
	Deposit  = "DEPOSIT"
	Withdraw = "WITHDRAW"
	Transfer = "TRANSFER"
)

const (
	ListTransactionsLimit = 100
)

var (
	ErrInvalidID       = errors.New("transaction id must be a valid ULID")
	ErrUnsupportedType = errors.New("unsupported transaction type")
)

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidID
	}
	return nil
}

func validOperationType(operationType string) error {
	var validOperationTypes = []string{
		Deposit,
		Withdraw,
		Transfer,
	}
	for _, validType := range validOperationTypes {
		if operationType == validType {
			return nil
		}
	}
	return ErrUnsupportedType
}
