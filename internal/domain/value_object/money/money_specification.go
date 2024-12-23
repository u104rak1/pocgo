package money

import "errors"

const (
	JPY = "JPY"
	USD = "USD"
)

var (
	ErrInvalidMoney               = errors.New("invalid money")
	ErrNegativeAmount             = errors.New("amount cannot be negative")
	ErrInvalidUSDPrecision        = errors.New("amount in USD cannot have more than 2 decimal places")
	ErrInvalidJPYPrecision        = errors.New("amount in JPY must not have decimal places")
	ErrUnsupportedCurrency        = errors.New("unsupported currency")
	ErrDifferentCurrencyOperation = errors.New("operation cannot be performed on different currencies")
	ErrInsufficientBalance        = errors.New("insufficient balance")
)

func validAmount(amount float64) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	return nil
}

func validCurrency(currency string) error {
	if !(currency == JPY || currency == USD) {
		return ErrUnsupportedCurrency
	}
	return nil
}
