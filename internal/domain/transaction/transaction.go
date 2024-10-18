package transaction

import (
	"time"

	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
)

type Transaction struct {
	id                string
	accountID         string
	receiverAccountID string
	transactionType   string
	transferAmount    money.Money
	transactionAt     time.Time
}

func New(
	id, accountID, receiverAccountID, transactionType string,
	amount float64,
	currency string,
	transactionAt time.Time,
) (*Transaction, error) {
	if err := ValidID(id); err != nil {
		return nil, err
	}

	if err := accountDomain.ValidID(accountID); err != nil {
		return nil, err
	}

	if err := accountDomain.ValidID(receiverAccountID); err != nil {
		return nil, err
	}

	if err := validTransactionType(transactionType); err != nil {
		return nil, err
	}

	transferAmount, err := money.New(amount, currency)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		id:                id,
		accountID:         accountID,
		receiverAccountID: receiverAccountID,
		transactionType:   transactionType,
		transferAmount:    *transferAmount,
		transactionAt:     transactionAt,
	}, nil
}

func (t *Transaction) ID() string {
	return t.id
}

func (t *Transaction) AccountID() string {
	return t.accountID
}

func (t *Transaction) ReceiverAccountID() string {
	return t.receiverAccountID
}

func (t *Transaction) TransactionType() string {
	return t.transactionType
}

func (t *Transaction) TransferAmount() money.Money {
	return t.transferAmount
}

func (t *Transaction) TransactionAt() time.Time {
	return t.transactionAt
}
