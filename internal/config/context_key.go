package config

import "errors"

type ctxUserIDKey struct{}

func CtxUserIDKey() interface{} {
	return ctxUserIDKey{}
}

var ErrUserIDMissing = errors.New("user id is missing")

type ctxTransactionKey struct{}

func CtxTransactionKey() interface{} {
	return ctxTransactionKey{}
}
