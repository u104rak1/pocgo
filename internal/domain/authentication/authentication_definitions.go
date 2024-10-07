package authentication_domain

import "errors"

var (
	ErrInvalidID             = errors.New("invalid authentication id")
	ErrPasswordInvalidLength = errors.New("password must be between 8 and 20 characters")
	ErrPasswordUnmatch       = errors.New("password unmatch")
)
