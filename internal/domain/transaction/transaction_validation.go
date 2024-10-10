package transaction_domain

import "github.com/ucho456job/pocgo/pkg/ulid"

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidTransactionID
	}
	return nil
}

func validTransactionType(transactionType string) error {
	var validTransactionTypes = []string{
		TransactionDeposit,
		TransactionWithdraw,
		TransactionTransfer,
	}
	for _, validType := range validTransactionTypes {
		if transactionType == validType {
			return nil
		}
	}
	return ErrUnsupportTransactionType
}
