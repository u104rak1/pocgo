package authentication

import (
	"errors"
	"fmt"
)

const (
	PasswordMinLength = 8
	PasswordMaxLength = 20
)

var (
	ErrPasswordInvalidLength   = fmt.Errorf("password must be between %d and %d characters", PasswordMinLength, PasswordMaxLength)
	ErrAlreadyExists           = errors.New("authentication already exists")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrAuthenticationFailed    = errors.New("email or password is incorrect")
	ErrNotFound                = errors.New("authentication not found")
	ErrUnmatchedPassword       = errors.New("passwords do not match")
)

func validPassword(password string) error {
	if len(password) < PasswordMinLength || len(password) > PasswordMaxLength {
		return ErrPasswordInvalidLength
	}
	return nil
}
