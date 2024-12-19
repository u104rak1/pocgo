package account

import (
	"time"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

type AccountID string

type Account struct {
	id           AccountID
	userID       userDomain.UserID
	name         string
	passwordHash string
	balance      moneyVO.Money
	updatedAt    time.Time
}

// 口座エンティティを作成します。新規で作成するのでパスワードの検証とハッシュ化を行います。
func New(userID userDomain.UserID, amount float64, name, password, currency string) (*Account, error) {
	id := AccountID(ulid.New())

	if err := validPassword(password); err != nil {
		return nil, err
	}
	passwordHash, err := passwordUtil.Encode(password)
	if err != nil {
		return nil, err
	}

	updatedAt := timer.Now()

	return newAccount(id, name, passwordHash, currency, userID, amount, updatedAt)
}

// データベースから口座を再構築します。パスワードは既にエンコードされているため、検証は行われません。
func Reconstruct(id, userID, name, passwordHash, currency string, amount float64, updatedAt time.Time) (*Account, error) {
	return newAccount(AccountID(id), name, passwordHash, currency, userDomain.UserID(userID), amount, updatedAt)
}

func newAccount(id AccountID, name, passwordHash, currency string, userID userDomain.UserID, amount float64, updatedAt time.Time) (*Account, error) {
	if err := validName(name); err != nil {
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

func (a *Account) ID() AccountID {
	return a.id
}

func (a *Account) IDString() string {
	return string(a.id)
}

func (a *Account) UserID() userDomain.UserID {
	return a.userID
}

func (a *Account) UserIDString() string {
	return string(a.userID)
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
