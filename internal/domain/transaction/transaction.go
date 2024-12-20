package transaction

import (
	"time"

	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/pkg/timer"
)

type Transaction struct {
	id                idVO.TransactionID
	accountID         idVO.AccountID
	receiverAccountID *idVO.AccountID
	operationType     string
	transferAmount    moneyVO.Money
	transactionAt     time.Time
}

// 取引エンティティを作成します。transactionAtは口座の更新日と同じ値にしたいので、引数で受け取ります。
func New(
	accountID idVO.AccountID,
	receiverAccountID *idVO.AccountID,
	operationType string,
	amount float64,
	currency string,
	transactionAt time.Time,
) (*Transaction, error) {
	id := idVO.NewTransactionID()
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
	tID, err := idVO.TransactionIDFromString(id)
	if err != nil {
		return nil, err
	}

	aID, err := idVO.AccountIDFromString(accountID)
	if err != nil {
		return nil, err
	}

	var raID *idVO.AccountID
	if receiverAccountID != nil {
		tmpID, err := idVO.AccountIDFromString(*receiverAccountID)
		if err != nil {
			return nil, err
		}
		raID = &tmpID
	}

	return newTransaction(tID, aID, raID, operationType, amount, currency, transactionAt)
}

func newTransaction(
	id idVO.TransactionID,
	accountID idVO.AccountID,
	receiverAccountID *idVO.AccountID,
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

func (t *Transaction) ID() idVO.TransactionID {
	return t.id
}

func (t *Transaction) IDString() string {
	return t.id.String()
}

func (t *Transaction) AccountID() idVO.AccountID {
	return t.accountID
}

func (t *Transaction) AccountIDString() string {
	return t.accountID.String()
}

func (t *Transaction) ReceiverAccountID() *idVO.AccountID {
	return t.receiverAccountID
}

func (t *Transaction) ReceiverAccountIDString() *string {
	if t.receiverAccountID == nil {
		return nil
	}
	receiverAccountID := t.receiverAccountID.String()
	return &receiverAccountID
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
