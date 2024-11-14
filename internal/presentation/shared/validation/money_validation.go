package validation

import (
	"math"

	v "github.com/go-ozzo/ozzo-validation/v4"
	moneyVO "github.com/ucho456job/pocgo/internal/domain/value_object/money"
)

func ValidJPY(amount float64) error {
	if err := v.Validate(amount, v.Min(0.0)); err != nil {
		return err
	}
	if math.Round(amount) != amount {
		return moneyVO.ErrInvalidJPYPrecision
	}
	return nil
}

func ValidUSD(amount float64) error {
	if err := v.Validate(amount, v.Min(0.0)); err != nil {
		return err
	}
	if math.Round(amount*100) != amount*100 {
		return moneyVO.ErrInvalidUSDPrecision
	}
	return nil
}

func ValidCurrency(currency string) error {
	return v.Validate(currency, v.Required, v.In(moneyVO.JPY, moneyVO.USD))
}
