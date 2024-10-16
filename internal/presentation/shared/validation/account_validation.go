package validation

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
)

func ValidAccountName(name string) error {
	return v.Validate(name, v.Required, v.Length(accountDomain.NameMinLength, accountDomain.NameMaxLength))
}

func ValidAccountPassword(password string) error {
	return v.Validate(password, v.Required, v.Length(accountDomain.PasswordLength, accountDomain.PasswordLength))
}

func ValidAccountCurrency(currency string) error {
	return v.Validate(currency, v.Required, v.In(accountDomain.JPY))
}
