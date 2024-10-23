package transaction

import (
	"errors"

	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	ErrInvalidTransactionID     = errors.New("transaction id must be a valid ULID")
	ErrUnsupportTransactionType = errors.New("unsupported transaction type")
)

const (
	Deposit  = "DEPOSIT"
	Withdraw = "WITHDRAW"
	Transfer = "TRANSFER"
)

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidTransactionID
	}
	return nil
}

func validTransactionType(transactionType string) error {
	var validTransactionTypes = []string{
		Deposit,
		Withdraw,
		Transfer,
	}
	for _, validType := range validTransactionTypes {
		if transactionType == validType {
			return nil
		}
	}
	return ErrUnsupportTransactionType
}
