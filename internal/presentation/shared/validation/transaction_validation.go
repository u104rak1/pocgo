package validation

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	transactionDomain "github.com/ucho456job/pocgo/internal/domain/transaction"
)

func ValidTransactionOperationType(operationType string) error {
	return v.Validate(operationType, v.Required, v.In(
		transactionDomain.Deposit, transactionDomain.Withdraw, transactionDomain.Transfer))
}
