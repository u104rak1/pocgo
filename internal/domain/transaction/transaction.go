package transaction_domain

import (
	"time"

	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
)

type Transaction struct {
	id                string
	senderAccountID   string
	receiverAccountID string
	transferAmount    accountDomain.Money
	transactionAt     time.Time
}

func New(
	id, senderAccountID, receiverAccountID string,
	amount float64,
	currency string,
	transactionAt time.Time,
) (*Transaction, error) {
	if err := ValidID(id); err != nil {
		return nil, err
	}

	if err := accountDomain.ValidID(senderAccountID); err != nil {
		return nil, err
	}

	if err := accountDomain.ValidID(receiverAccountID); err != nil {
		return nil, err
	}

	transferAmount, err := accountDomain.NewMoney(amount, currency)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		id:                id,
		senderAccountID:   senderAccountID,
		receiverAccountID: receiverAccountID,
		transferAmount:    *transferAmount,
		transactionAt:     transactionAt,
	}, nil
}

func (t *Transaction) ID() string {
	return t.id
}

func (t *Transaction) SenderAccountID() string {
	return t.senderAccountID
}

func (t *Transaction) ReceiverAccountID() string {
	return t.receiverAccountID
}

func (t *Transaction) TransferAmount() accountDomain.Money {
	return t.transferAmount
}

func (t *Transaction) TransactionAt() time.Time {
	return t.transactionAt
}
