package id

import "fmt"

type transactionIDType struct{}

type TransactionID = ID[transactionIDType]

func NewTransactionID() TransactionID {
	return New[transactionIDType]()
}

func TransactionIDFromString(value string) (TransactionID, error) {
	transactionID, err := NewFromString[transactionIDType](value)
	if err != nil {
		return TransactionID{}, fmt.Errorf("invalid transaction id: %w", err)
	}
	return transactionID, nil
}

// NewTransactionIDForTest テスト用のTransactionIDを生成します
// 同じseedからは常に同じIDが生成されます
func NewTransactionIDForTest(seed string) TransactionID {
	return NewForTest[transactionIDType](seed)
}
