package id

type transactionIDType struct{}

type TransactionID = ID[transactionIDType]

func NewTransactionID() TransactionID {
	return New[transactionIDType]()
}

func TransactionIDFromString(value string) (TransactionID, error) {
	return NewFromString[transactionIDType](value)
}

// NewTransactionIDForTest テスト用のTransactionIDを生成します
// 同じseedからは常に同じIDが生成されます
func NewTransactionIDForTest(seed string) TransactionID {
	return NewForTest[transactionIDType](seed)
}
