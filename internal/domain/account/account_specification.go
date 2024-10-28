package account

import (
	"errors"

	"github.com/ucho456job/pocgo/pkg/ulid"
)

const (
	NameMinLength  = 3
	NameMaxLength  = 20
	PasswordLength = 4
)

var (
	ErrInvalidID             = errors.New("account id must be a valid ULID")
	ErrInvalidAccountName    = errors.New("account name must be between 1 and 10 characters")
	ErrPasswordInvalidLength = errors.New("account password must be 4 characters")
	ErrAccountNotFound       = errors.New("account not found")
	ErrUnmatchedPassword     = errors.New("passwords do not match")
)

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidID
	}
	return nil
}

func validName(name string) error {
	if len(name) < NameMinLength || len(name) > NameMaxLength {
		return ErrInvalidAccountName
	}
	return nil
}

func validPassword(password string) error {
	if len(password) != PasswordLength {
		return ErrPasswordInvalidLength
	}
	return nil
}
