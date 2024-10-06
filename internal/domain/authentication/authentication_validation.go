package password_domain

import (
	"errors"

	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	ErrInvalidID             = errors.New("invalid authentication id")
	ErrPasswordInvalidLength = errors.New("password must be between 8 and 20 characters")
	ErrPasswordUnmatch       = errors.New("password unmatch")
)

func IsValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidID
	}
	return nil
}

func isValidPassword(password string) error {
	const passwordMinLength = 8
	const passwordMaxLength = 20
	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return ErrPasswordInvalidLength
	}
	return nil
}
