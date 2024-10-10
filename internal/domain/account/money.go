package account_domain

type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) (*Money, error) {
	if err := validAmount(amount); err != nil {
		return nil, err
	}

	if err := validCurrency(currency); err != nil {
		return nil, err
	}

	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Sub(other Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, ErrDifferentCurrency
	}
	if m.amount < other.amount {
		return nil, ErrInsufficientFunds
	}
	return &Money{amount: m.amount - other.amount, currency: m.currency}, nil
}

func (m Money) Add(other Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, ErrDifferentCurrency
	}
	return &Money{amount: m.amount + other.amount, currency: m.currency}, nil
}
