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

func ValidPassword(password string) error {
	const passwordMinLength = 8
	const passwordMaxLength = 20
	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return ErrPasswordInvalidLength
	}
	return nil
}
