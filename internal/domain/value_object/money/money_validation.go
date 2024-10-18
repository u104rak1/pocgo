package money

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
