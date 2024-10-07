package account_domain

import (
	"time"

	user_domain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
)

type Account struct {
	id            string
	userID        string
	name          string
	passwordHash  string
	balance       Money
	lastUpdatedAt time.Time
}

func New(id, userID, name, password string, amount float64, currency string, lastUpdatedAt time.Time) (*Account, error) {
	var err error
	if err = validID(id); err != nil {
		return nil, err
	}

	if err = validName(name); err != nil {
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
		name:          name,
		passwordHash:  passwordHash,
		balance:       *balance,
		lastUpdatedAt: lastUpdatedAt,
	}, nil
}

func (a *Account) ID() string {
	return a.id
}

func (a *Account) UserID() string {
	return a.userID
}

func (a *Account) Name() string {
	return a.name
}

func (a *Account) Balance() Money {
	return a.balance
}

func (a *Account) LastUpdatedAt() time.Time {
	return a.lastUpdatedAt
}

func (a *Account) ChangeName(new string) error {
	if err := validName(new); err != nil {
		return err
	}
	a.name = new
	a.UpdateLastUpdatedAt(time.Now())
	return nil
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

func (a *Account) Withdraw(amount float64, currency string) error {
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

func (a *Account) Deposit(amount float64, currency string) error {
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
