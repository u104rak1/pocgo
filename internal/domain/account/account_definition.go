package account_domain

import "errors"

var (
	ErrInvalidID             = errors.New("invalid account id")
	ErrInvalidAccountName    = errors.New("account name must be between 1 and 10 characters")
	ErrPasswordInvalidLength = errors.New("account password must be 4 characters")
	ErrNegativeAmount        = errors.New("amount cannot be negative")
	ErrUnsupportedCurrency   = errors.New("unsupported currency")
	ErrDifferentCurrency     = errors.New("different currency")
	ErrInsufficientFunds     = errors.New("insufficient funds")
)