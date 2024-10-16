package authentication

import (
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidID
	}
	return nil
}

func validPassword(password string) error {
	if len(password) < PasswordMinLength || len(password) > PasswordMaxLength {
		return ErrPasswordInvalidLength
	}
	return nil
}
