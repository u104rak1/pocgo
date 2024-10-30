package account

import (
	"errors"
	"fmt"

	"github.com/ucho456job/pocgo/pkg/ulid"
)

const (
	NameMinLength  = 3
	NameMaxLength  = 20
	PasswordLength = 4
)

var (
	ErrInvalidID             = errors.New("account id must be a valid ULID")
	ErrInvalidName           = fmt.Errorf("account name must be between %d and %d characters", NameMinLength, NameMaxLength)
	ErrPasswordInvalidLength = fmt.Errorf("account password must be %d characters", PasswordLength)
	ErrNotFound              = errors.New("account not found")
	ErrUnmatchedPassword     = errors.New("passwords do not match")
	ErrLimitReached          = errors.New("account limit reached")
)

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidID
	}
	return nil
}

func validName(name string) error {
	if len(name) < NameMinLength || len(name) > NameMaxLength {
		return ErrInvalidName
	}
	return nil
}

func validPassword(password string) error {
	if len(password) != PasswordLength {
		return ErrPasswordInvalidLength
	}
	return nil
}
