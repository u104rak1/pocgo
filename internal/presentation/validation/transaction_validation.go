package validation

import (
	"errors"
	"strconv"
	"strings"

	v "github.com/go-ozzo/ozzo-validation/v4"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
)

func ValidTransactionOperationType(operationType string) error {
	return v.Validate(operationType, v.Required, v.In(
		transactionDomain.Deposit, transactionDomain.Withdrawal, transactionDomain.Transfer))
}

// 取引操作タイプのカンマ区切り文字列を検証します。
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
	if limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	if limit > transactionDomain.ListTransactionsLimit {
		return errors.New("limit must be less than or equal to " + strconv.Itoa(transactionDomain.ListTransactionsLimit))
	}
	return nil
}
