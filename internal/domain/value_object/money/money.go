package money

import (
	"math"
)

type Money struct {
	amount   float64
	currency string
}

func New(amount float64, currency string) (*Money, error) {
	if err := validAmount(amount); err != nil {
		return nil, err
	}

	if err := validCurrency(currency); err != nil {
		return nil, err
	}

	switch currency {
	case JPY:
		if math.Round(amount) != amount {
			return nil, ErrInvalidJPYPrecision
		}
		amount = math.Round(amount)
		return &Money{amount: amount, currency: currency}, nil
	case USD:
		if math.Round(amount*100) != amount*100 {
			return nil, ErrInvalidUSDPrecision
		}
		return &Money{amount: amount, currency: currency}, nil
	default:
		return nil, ErrInvalidMoney
	}
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, ErrAddDifferentCurrency
	}
	return &Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Sub(other Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, ErrSubDifferentCurrency
	}
	if m.amount < other.amount {
		return nil, ErrInsufficientBalance
	}
	return &Money{amount: m.amount - other.amount, currency: m.currency}, nil
}
