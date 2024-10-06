package account_domain

import "errors"

var (
	ErrInvalidID             = errors.New("invalid account id")
	ErrPasswordInvalidLength = errors.New("account password must be 4 characters")
	ErrNegativeAmount        = errors.New("amount cannot be negative")
	ErrUnsupportedCurrency   = errors.New("unsupported currency")
	ErrDifferentCurrency     = errors.New("different currency")
	ErrInsufficientFunds     = errors.New("insufficient funds")
)

func IsValidID(id string) error {
	if id == "" {
		return ErrInvalidID
	}
	return nil
}

func validPassword(password string) error {
	const passwordLength = 4
	if len(password) == passwordLength {
		return ErrPasswordInvalidLength
	}
	return nil
}

func validAmount(amount int) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	return nil
}

type Currency string

const (
	JPY Currency = "JPY"
)

func validCurrency(currency Currency) error {
	if currency != JPY {
		return ErrUnsupportedCurrency
	}
	return nil
}
