package config

import "errors"

// CtxUserIDKey is a key for user id in context
type ctxUserIDKey struct{}

func CtxUserIDKey() interface{} {
	return ctxUserIDKey{}
}

var ErrUserIDMissing = errors.New("user id is missing")

type ctxTransactionKey struct{}

func CtxTransactionKey() interface{} {
	return ctxTransactionKey{}
}
