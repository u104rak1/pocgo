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
	accountID         string
	receiverAccountID *string
	operationType     string
	transferAmount    moneyVO.Money
	transactionAt     time.Time
}

func New(
	accountID string,
	receiverAccountID *string,
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
	return newTransaction(transactionID, accountID, receiverAccountID, operationType, amount, currency, transactionAt)
}

func newTransaction(
	id TransactionID,
	accountID string,
	receiverAccountID *string,
	operationType string,
	amount float64,
	currency string,
	transactionAt time.Time,
) (*Transaction, error) {
	if err := accountDomain.ValidID(accountID); err != nil {
		return nil, err
	}

	if receiverAccountID != nil {
		if err := accountDomain.ValidID(*receiverAccountID); err != nil {
			return nil, err
		}
	}

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

func (t *Transaction) AccountID() string {
	return t.accountID
}

func (t *Transaction) ReceiverAccountID() *string {
	return t.receiverAccountID
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
