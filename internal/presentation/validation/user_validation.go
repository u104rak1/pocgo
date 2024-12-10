package validation

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	emailUtil "github.com/u104rak1/pocgo/pkg/email"
)

func ValidUserName(name string) error {
	return v.Validate(name, v.Required, v.Length(userDomain.NameMinLength, userDomain.NameMaxLength))
}

func ValidUserEmail(email string) error {
	return v.Validate(email, v.Required, v.By(func(value interface{}) error {
		if !emailUtil.IsValid(value.(string)) {
			return v.NewError("invalid_email", userDomain.ErrInvalidEmail.Error())
		}
		return nil
	}))
}

func ValidUserPassword(password string) error {
	return v.Validate(password, v.Required, v.Length(authDomain.PasswordMinLength, authDomain.PasswordMaxLength))
}
