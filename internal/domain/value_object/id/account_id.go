package id

import "fmt"

type accountIDType struct{}

type AccountID = ID[accountIDType]

func NewAccountID() AccountID {
	return New[accountIDType]()
}

func AccountIDFromString(value string) (AccountID, error) {
	accountID, err := NewFromString[accountIDType](value)
	if err != nil {
		return AccountID{}, fmt.Errorf("invalid account id: %w", err)
	}
	return accountID, nil
}

// NewAccountIDForTest テスト用のAccountIDを生成します
// 同じseedからは常に同じIDが生成されます
func NewAccountIDForTest(seed string) AccountID {
	return NewForTest[accountIDType](seed)
}
