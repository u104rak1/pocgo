package user

import (
	"errors"
	"fmt"

	emailUtil "github.com/u104rak1/pocgo/pkg/email"
)

const (
	NameMinLength = 3
	NameMaxLength = 20
)

var (
	ErrInvalidName        = fmt.Errorf("user name must be between %d and %d characters", NameMinLength, NameMaxLength)
	ErrInvalidEmail       = errors.New("the email format is invalid")
	ErrEmailAlreadyExists = errors.New("user email already exists")
	ErrNotFound           = errors.New("user not found")
)

func validName(name string) error {
	if len(name) < NameMinLength || len(name) > NameMaxLength {
		return ErrInvalidName
	}
	return nil
}

func validEmail(email string) error {
	if !emailUtil.IsValid(email) {
		return ErrInvalidEmail
	}
	return nil
}
