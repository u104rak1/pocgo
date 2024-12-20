package id

import "fmt"

type userIDType struct{}

type UserID = ID[userIDType]

func NewUserID() UserID {
	return New[userIDType]()
}

func UserIDFromString(value string) (UserID, error) {
	userID, err := NewFromString[userIDType](value)
	if err != nil {
		return UserID{}, fmt.Errorf("invalid user id: %w", err)
	}
	return userID, nil
}

// NewUserIDForTest テスト用のUserIDを生成します
// 同じseedからは常に同じIDが生成されます
func NewUserIDForTest(seed string) UserID {
	return NewForTest[userIDType](seed)
}
