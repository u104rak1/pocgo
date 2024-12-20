package id

type userIDType struct{}

type UserID = ID[userIDType]

func NewUserID() UserID {
	return New[userIDType]()
}

func UserIDFromString(value string) (UserID, error) {
	return NewFromString[userIDType](value)
}

// NewUserIDForTest テスト用のUserIDを生成します
// 同じseedからは常に同じIDが生成されます
func NewUserIDForTest(seed string) UserID {
	return NewForTest[userIDType](seed)
}
