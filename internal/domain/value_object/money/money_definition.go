package money

import "errors"

const (
	JPY = "JPY"
	USD = "USD"
)

var (
	ErrNegativeAmount      = errors.New("amount cannot be negative")
	ErrInvalidUSDPrecision = errors.New("USD amount must have two decimal places")
	ErrInvalidJPYPrecision = errors.New("JPY amount must have no decimal places")
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrDifferentCurrency   = errors.New("cannot add different currencies")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
