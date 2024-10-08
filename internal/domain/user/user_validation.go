package user_domain

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

func ValidName(name string) error {
	const nameMinLength = 1
	const nameMaxLength = 20
	if len(name) < nameMinLength || len(name) > nameMaxLength {
		return ErrInvalidUserName
	}
	return nil
}

func ValidEmail(email string) error {
	if !emailUtil.IsValid(email) {
		return ErrInvalidEmail
	}
	return nil
}
