package account

func ValidID(id string) error {
	if id == "" {
		return ErrInvalidID
	}
	return nil
}

func validName(name string) error {
	const nameMinLength = 1
	const nameMaxLength = 10
	if len(name) < nameMinLength || len(name) > nameMaxLength {
		return ErrInvalidAccountName
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

func validAmount(amount float64) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	return nil
}

func validCurrency(currency string) error {
	JPY := "JPY"
	if currency != JPY {
		return ErrUnsupportedCurrency
	}
	return nil
}
