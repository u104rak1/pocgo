package account

func ValidID(id string) error {
	if id == "" {
		return ErrInvalidID
	}
	return nil
}

func validName(name string) error {
	if len(name) < NameMinLength || len(name) > NameMaxLength {
		return ErrInvalidAccountName
	}
	return nil
}

func validPassword(password string) error {
	if len(password) != PasswordLength {
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
	if currency != JPY {
		return ErrUnsupportedCurrency
	}
	return nil
}
