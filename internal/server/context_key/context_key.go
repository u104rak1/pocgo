package contextkey

type userIDKey struct{}

func UserIDKey() interface{} {
	return userIDKey{}
}
