package authentication

import "errors"

var (
	ErrInvalidID                   = errors.New("invalid authentication id")
	ErrPasswordInvalidLength       = errors.New("password must be between 8 and 20 characters")
	ErrPasswordUnmatch             = errors.New("password unmatch")
	ErrAuthenticationAlreadyExists = errors.New("authentication already exists")
	ErrUnexpectedSigningMethod     = errors.New("unexpected signing method")
	ErrAuthenticationFailed        = errors.New("invalid token or missing userID")
)
