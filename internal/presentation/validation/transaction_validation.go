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
	if operationTypes == "" {
		return errors.New("operation types cannot be blank")
	}

	types := strings.Split(operationTypes, ",")
	var errorMsgs []string

	for _, t := range types {
		t = strings.TrimSpace(t)
		if err := ValidTransactionOperationType(t); err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		}
	}

	if len(errorMsgs) > 0 {
		return errors.New("contains an invalid operation type")
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
