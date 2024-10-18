package transaction

import "github.com/ucho456job/pocgo/pkg/ulid"

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
