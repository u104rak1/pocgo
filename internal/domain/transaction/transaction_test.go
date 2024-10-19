package transaction_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/internal/domain/transaction"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	validID                = ulid.GenerateStaticULID("transaction")
	validAccountID         = ulid.GenerateStaticULID("account")
	validRecieverAccountID = ulid.GenerateStaticULID("accountReceiver")
	validAmount            = 1000.0
	validCurrency          = money.JPY
	validTransactionAt     = time.Now()
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		caseName          string
		id                string
		accountID         string
		receiverAccountID string
		transactionType   string
		amount            float64
		currency          string
		transactionAt     time.Time
		wantErr           error
	}{
		{
			caseName:          "Happy path: return transaction entity, if arguments are valid.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Transfer,
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Edge case: return error, if the ID is invalid.",
			id:                "invalid",
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Transfer,
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           transaction.ErrInvalidTransactionID,
		},
		{
			caseName:          "Edge case: return error, if the AccountID is invalid.",
			id:                validID,
			accountID:         "inavlid",
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Transfer,
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           accountDomain.ErrInvalidID,
		},
		{
			caseName:          "Edge case: return error, if the ReceiverAccountID is invalid.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: "inavlid",
			transactionType:   transaction.Transfer,
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           accountDomain.ErrInvalidID,
		},
		{
			caseName:          "Happy path: return transaction entity, if the transaction type is TRANSFER.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Transfer,
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Happy path: return transaction entity, if the transaction type is DEPOSIT.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Deposit,
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Happy path: return transaction entity, if the transaction type is WITHDRAW.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Withdraw,
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Happy path: return error, if the transaction type is unsupported.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   "UNSUPPORTED",
			amount:            validAmount,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           transaction.ErrUnsupportTransactionType,
		},
		{
			caseName:          "Edge case: return error, if the amount is invalid.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Transfer,
			amount:            -1000,
			currency:          validCurrency,
			transactionAt:     validTransactionAt,
			wantErr:           money.ErrNegativeAmount,
		},
		{
			caseName:          "Edge case: return error, if the currency is unsupported.",
			id:                validID,
			accountID:         validAccountID,
			receiverAccountID: validRecieverAccountID,
			transactionType:   transaction.Transfer,
			amount:            validAmount,
			currency:          "EUR",
			transactionAt:     validTransactionAt,
			wantErr:           money.ErrUnsupportedCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			tx, err := transaction.New(
				tt.id, tt.accountID, tt.receiverAccountID, tt.transactionType, tt.amount, tt.currency, tt.transactionAt,
			)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, tx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tx)
				assert.Equal(t, tt.id, tx.ID())
				assert.Equal(t, tt.accountID, tx.AccountID())
				assert.Equal(t, tt.receiverAccountID, tx.ReceiverAccountID())
				assert.Equal(t, tt.transactionType, tx.TransactionType())
				assert.Equal(t, tt.amount, tx.TransferAmount().Amount())
				assert.Equal(t, tt.currency, tx.TransferAmount().Currency())
				assert.Equal(t, tt.transactionAt, tx.TransactionAt())
			}
		})
	}
}
