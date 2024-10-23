package transaction_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/internal/domain/transaction"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	validID                = ulid.GenerateStaticULID("transaction")
	validAccountID         = ulid.GenerateStaticULID("account")
	validRecieverAccountID = ulid.GenerateStaticULID("accountReceiver")
	validAmount            = 1000.0
	validCurrency          = money.JPY
	validTransactionAt     = timer.Now()
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
			caseName:          "Successfully creates a transaction.",
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
			caseName:          "Error occurs with invalid ID.",
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
			caseName:          "Error occurs with invalid AccountID.",
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
			caseName:          "Error occurs with invalid ReceiverAccountID.",
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
			caseName:          "Successfully creates a transaction with TRANSFER transaction type.",
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
			caseName:          "Successfully creates a transaction with DEPOSIT transaction type.",
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
			caseName:          "Successfully creates a transaction with WITHDRAW transaction type.",
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
			caseName:          "Error occurs with unsupported transaction type.",
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
			caseName:          "Error occurs with invalid amount.",
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
			caseName:          "Error occurs with unsupported currency.",
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
