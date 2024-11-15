package validation

import (
	"strings"

	v "github.com/go-ozzo/ozzo-validation/v4"
	transactionDomain "github.com/ucho456job/pocgo/internal/domain/transaction"
)

func ValidTransactionOperationType(operationType string) error {
	return v.Validate(operationType, v.Required, v.In(
		transactionDomain.Deposit, transactionDomain.Withdraw, transactionDomain.Transfer))
}

// ValidTransactionOperationTypes validates a comma-separated string of transaction operation types.
func ValidTransactionOperationTypes(operationTypes string) error {
	types := strings.Split(operationTypes, ",")
	var errors []string

	for _, t := range types {
		t = strings.TrimSpace(t)
		if err := ValidTransactionOperationType(t); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return v.Errors{
			"operation_types": v.NewError("validation", strings.Join(errors, "; ")),
		}
	}

	return nil
}

func ValidListTransactionsLimit(limit int) error {
	return v.Validate(limit, v.Min(1), v.Max(transactionDomain.ListTransactionsLimit))
}
