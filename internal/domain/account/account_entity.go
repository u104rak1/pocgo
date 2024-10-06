package account_domain

import (
	"time"

	user_domain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
)

type Account struct {
	id            string
	userID        string
	passwordHash  string
	balance       Money
	lastUpdatedAt time.Time
}

func New(id, userID, password string, amount int, currency Currency, lastUpdatedAt time.Time) (*Account, error) {
	var err error
	if err = IsValidID(id); err != nil {
		return nil, err
	}

	if err = user_domain.IsValidID(userID); err != nil {
		return nil, err
	}

	if err = validPassword(password); err != nil {
		return nil, err
	}
	passwordHash := passwordUtil.Encode(password)

	balance, err := NewMoney(amount, currency)
	if err != nil {
		return nil, err
	}

	return &Account{
		id:            id,
		userID:        userID,
		passwordHash:  passwordHash,
		balance:       *balance,
		lastUpdatedAt: lastUpdatedAt,
	}, nil
}

func (a *Account) ChangePassword(new string) error {
	if err := validPassword(new); err != nil {
		return err
	}
	passwordHash := passwordUtil.Encode(new)
	a.passwordHash = passwordHash
	a.UpdateLastUpdatedAt(time.Now())
	return nil
}

func (a *Account) ComparePassword(password string) error {
	return passwordUtil.Compare(a.passwordHash, password)
}

func (a *Account) Withdraw(amount int, currency Currency) error {
	withdrawMoney, err := NewMoney(amount, currency)
	if err != nil {
		return err
	}

	newBalance, err := a.balance.Sub(*withdrawMoney)
	if err != nil {
		return err
	}

	a.balance = *newBalance
	a.UpdateLastUpdatedAt(time.Now())
	return nil
}

func (a *Account) Deposit(amount int, currency Currency) error {
	depositMoney, err := NewMoney(amount, currency)
	if err != nil {
		return err
	}

	newBalance, err := a.balance.Add(*depositMoney)
	if err != nil {
		return err
	}

	a.balance = *newBalance
	a.UpdateLastUpdatedAt(time.Now())
	return nil
}

func (a *Account) UpdateLastUpdatedAt(now time.Time) {
	a.lastUpdatedAt = now
}

type Money struct {
	amount   int
	currency Currency
}

func NewMoney(amount int, currency Currency) (*Money, error) {
	var err error
	if err = validAmount(amount); err != nil {
		return nil, err
	}

	if err = validCurrency(currency); err != nil {
		return nil, err
	}

	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
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
