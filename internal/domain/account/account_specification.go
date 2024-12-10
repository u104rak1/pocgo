package account

import (
	"errors"
	"fmt"

	"github.com/u104rak1/pocgo/pkg/ulid"
)

const (
	NameMinLength   = 3
	NameMaxLength   = 20
	PasswordLength  = 4
	MaxAccountLimit = 3
)

var (
	ErrInvalidID             = errors.New("account id must be a valid ULID")
	ErrInvalidName           = fmt.Errorf("account name must be between %d and %d characters", NameMinLength, NameMaxLength)
	ErrPasswordInvalidLength = fmt.Errorf("account password must be %d characters", PasswordLength)
	ErrNotFound              = errors.New("account not found")
	ErrReceiverNotFound      = errors.New("receiver account not found")
	ErrUnmatchedPassword     = errors.New("passwords do not match")
	ErrLimitReached          = fmt.Errorf("account limit reached, maximum %d accounts", MaxAccountLimit)
	ErrUnauthorized          = errors.New("unauthorized access to account")
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
