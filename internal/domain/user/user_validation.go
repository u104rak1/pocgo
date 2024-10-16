package user

import (
	emailUtil "github.com/ucho456job/pocgo/pkg/email"
	"github.com/ucho456job/pocgo/pkg/ulid"
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
