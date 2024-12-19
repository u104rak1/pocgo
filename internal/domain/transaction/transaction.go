package transaction

import (
	"time"

	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

type TransactionID string

type Transaction struct {
	id                TransactionID
	accountID         accountDomain.AccountID
	receiverAccountID *accountDomain.AccountID
	operationType     string
	transferAmount    moneyVO.Money
	transactionAt     time.Time
}

// 取引エンティティを作成します。transactionAtは口座の更新日と同じ値にしたいので、引数で受け取ります。
func New(
	accountID accountDomain.AccountID,
	receiverAccountID *accountDomain.AccountID,
	operationType string,
	amount float64,
	currency string,
	transactionAt time.Time,
) (*Transaction, error) {
	id := TransactionID(ulid.New())
	return newTransaction(id, accountID, receiverAccountID, operationType, amount, currency, transactionAt)
}

func Reconstruct(
	id, accountID string,
	receiverAccountID *string,
	operationType string,
	amount float64,
	currency string,
	transactionAt time.Time,
) (*Transaction, error) {
	transactionID := TransactionID(id)
	aID := accountDomain.AccountID(accountID)

	var raID *accountDomain.AccountID
	if receiverAccountID != nil {
		tempRaID := accountDomain.AccountID(*receiverAccountID)
		raID = &tempRaID
	}

	return newTransaction(transactionID, aID, raID, operationType, amount, currency, transactionAt)
}

func newTransaction(
	id TransactionID,
	accountID accountDomain.AccountID,
	receiverAccountID *accountDomain.AccountID,
	operationType string,
	amount float64,
	currency string,
	transactionAt time.Time,
) (*Transaction, error) {
	if err := validOperationType(operationType); err != nil {
		return nil, err
	}

	transferAmount, err := moneyVO.New(amount, currency)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		id:                id,
		accountID:         accountID,
		receiverAccountID: receiverAccountID,
		operationType:     operationType,
		transferAmount:    *transferAmount,
		transactionAt:     transactionAt,
	}, nil
}

func (t *Transaction) ID() TransactionID {
	return t.id
}

func (t *Transaction) IDString() string {
	return string(t.id)
}

func (t *Transaction) AccountID() accountDomain.AccountID {
	return t.accountID
}

func (t *Transaction) AccountIDString() string {
	return string(t.accountID)
}

func (t *Transaction) ReceiverAccountID() *accountDomain.AccountID {
	return t.receiverAccountID
}

func (t *Transaction) ReceiverAccountIDString() string {
	if t.receiverAccountID == nil {
		return ""
	}
	return string(*t.receiverAccountID)
}

func (t *Transaction) OperationType() string {
	return t.operationType
}

func (t *Transaction) TransferAmount() moneyVO.Money {
	return t.transferAmount
}

func (t *Transaction) TransactionAt() time.Time {
	return t.transactionAt
}

func (t *Transaction) TransactionAtString() string {
	return timer.FormatToISO8601(t.transactionAt)
}
