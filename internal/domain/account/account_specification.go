package account

import (
	"errors"

	"github.com/ucho456job/pocgo/pkg/ulid"
)

const (
	NameMinLength  = 1
	NameMaxLength  = 10
	PasswordLength = 4
	JPY            = "JPY"
)

var (
	ErrInvalidID             = errors.New("invalid account id")
	ErrInvalidAccountName    = errors.New("account name must be between 1 and 10 characters")
	ErrPasswordInvalidLength = errors.New("account password must be 4 characters")
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