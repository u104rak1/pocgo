package validation

import (
	"math"

	v "github.com/go-ozzo/ozzo-validation/v4"
	moneyVO "github.com/ucho456job/pocgo/internal/domain/value_object/money"
)

type MoneyValidation struct {
	Field   string
	Message string
}

func ValidCurrency(currency string) error {
	return v.Validate(currency, v.Required, v.In(moneyVO.JPY, moneyVO.USD))
}

func ValidAmount(currency string, amount float64) error {
	switch currency {
	case moneyVO.JPY:
		return validJPY(amount)
	case moneyVO.USD:
		return validUSD(amount)
	default:
		return moneyVO.ErrUnsupportedCurrency
	}
}

func validJPY(amount float64) error {
	if err := v.Validate(amount, v.Min(0.0)); err != nil {
		return err
	}
	if math.Round(amount) != amount {
		return moneyVO.ErrInvalidJPYPrecision
	}
	return nil
}

func validUSD(amount float64) error {
	if err := v.Validate(amount, v.Min(0.0)); err != nil {
		return err
	}
	if math.Round(amount*100) != amount*100 {
		return moneyVO.ErrInvalidUSDPrecision
	}
	return nil
}
