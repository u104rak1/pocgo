package user

import (
	"errors"

	emailUtil "github.com/ucho456job/pocgo/pkg/email"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

const (
	NameMinLength = 3
	NameMaxLength = 20
)

var (
	ErrInvalidUserID          = errors.New("user id must be a valid ULID")
	ErrInvalidUserName        = errors.New("user name must be between 3 and 20 characters")
	ErrInvalidEmail           = errors.New("the email format is invalid")
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
)

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidUserID
	}
	return nil
}

func validName(name string) error {
	if len(name) < NameMinLength || len(name) > NameMaxLength {
		return ErrInvalidUserName
	}
	return nil
}

func validEmail(email string) error {
	if !emailUtil.IsValid(email) {
		return ErrInvalidEmail
	}
	return nil
}
