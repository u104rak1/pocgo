package user_domain

import (
	"errors"

	emailUtil "github.com/ucho456job/pocgo/pkg/email"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	ErrInvalidUserID          = errors.New("invalid user id")
	ErrInvalidUserName        = errors.New("user name must be between 1 and 20 characters")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
)

// idのみ特別にexportしている。要検討。
func IsValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidUserID
	}
	return nil
}

func isValidName(name string) error {
	const nameMinLength = 1
	const nameMaxLength = 20
	if len(name) < nameMinLength || len(name) > nameMaxLength {
		return ErrInvalidUserName
	}
	return nil
}

func isValidEmail(email string) error {
	if !emailUtil.IsValid(email) {
		return ErrInvalidEmail
	}
	return nil
}
