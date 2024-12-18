package account

import (
	"time"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
	"github.com/u104rak1/pocgo/pkg/timer"
)

type Account struct {
	id           string
	userID       string
	name         string
	passwordHash string
	balance      moneyVO.Money
	updatedAt    time.Time
}

// 口座エンティティを作成します。新規で作成するのでパスワードの検証とハッシュ化を行います。
func New(id, userID, name, password string, amount float64, currency string, updatedAt time.Time) (*Account, error) {
	if err := validPassword(password); err != nil {
		return nil, err
	}
	passwordHash, err := passwordUtil.Encode(password)
	if err != nil {
		return nil, err
	}

	return newAccount(id, userID, name, passwordHash, amount, currency, updatedAt)
}

// データベースから口座を再構築します。パスワードは既にエンコードされているため、検証は行われません。
func Reconstruct(id, userID, name, passwordHash string, amount float64, currency string, updatedAt time.Time) (*Account, error) {
	return newAccount(id, userID, name, passwordHash, amount, currency, updatedAt)
}

func newAccount(id, userID, name, passwordHash string, amount float64, currency string, updatedAt time.Time) (*Account, error) {
	if err := ValidID(id); err != nil {
		return nil, err
	}

	if err := validName(name); err != nil {
		return nil, err
	}

	if err := userDomain.ValidID(userID); err != nil {
		return nil, err
	}

	balance, err := moneyVO.New(amount, currency)
	if err != nil {
		return nil, err
	}

	return &Account{
		id:           id,
		userID:       userID,
		name:         name,
		passwordHash: passwordHash,
		balance:      *balance,
		updatedAt:    updatedAt,
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

func (a *Account) PasswordHash() string {
	return a.passwordHash
}

func (a *Account) Balance() moneyVO.Money {
	return a.balance
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) UpdatedAtString() string {
	return timer.FormatToISO8601(a.updatedAt)
}

func (a *Account) ChangeName(new string) error {
	if err := validName(new); err != nil {
		return err
	}
	a.name = new
	return nil
}

func (a *Account) ChangePassword(new string) error {
	if err := validPassword(new); err != nil {
		return err
	}
	passwordHash, err := passwordUtil.Encode(new)
	if err != nil {
		return err
	}
	a.passwordHash = passwordHash
	return nil
}

func (a *Account) ComparePassword(password string) error {
	if err := passwordUtil.Compare(a.passwordHash, password); err != nil {
		return ErrUnmatchedPassword
	}
	return nil
}

func (a *Account) Withdraw(amount float64, currency string) error {
	withdrawMoney, err := moneyVO.New(amount, currency)
	if err != nil {
		return err
	}

	newBalance, err := a.balance.Sub(*withdrawMoney)
	if err != nil {
		return err
	}

	a.balance = *newBalance
	return nil
}

func (a *Account) Deposit(amount float64, currency string) error {
	depositMoney, err := moneyVO.New(amount, currency)
	if err != nil {
		return err
	}

	newBalance, err := a.balance.Add(*depositMoney)
	if err != nil {
		return err
	}

	a.balance = *newBalance
	return nil
}

func (a *Account) ChangeUpdatedAt(now time.Time) {
	a.updatedAt = now
}
