package transaction_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	"github.com/u104rak1/pocgo/internal/domain/transaction"
	"github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

func TestNewTransaction(t *testing.T) {

	var (
		transactionID     = ulid.GenerateStaticULID("transaction")
		accountID         = ulid.GenerateStaticULID("account")
		receiverAccountID = ulid.GenerateStaticULID("accountReceiver")
		amount            = 1000.0
		currency          = money.JPY
		transactionAt     = timer.GetFixedDate()
		invalidID         = "invalid id"
	)

	tests := []struct {
		caseName          string
		id                string
		accountID         string
		receiverAccountID *string
		operationType     string
		amount            float64
		currency          string
		transactionAt     time.Time
		wantErr           error
	}{
		{
			caseName:          "Successfully creates a transaction.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transaction.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Error occurs with invalid ID.",
			id:                invalidID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transaction.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           transaction.ErrInvalidID,
		},
		{
			caseName:          "Error occurs with invalid AccountID.",
			id:                transactionID,
			accountID:         "inavlid",
			receiverAccountID: &receiverAccountID,
			operationType:     transaction.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           accountDomain.ErrInvalidID,
		},
		{
			caseName:          "Error occurs with invalid ReceiverAccountID.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &invalidID,
			operationType:     transaction.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           accountDomain.ErrInvalidID,
		},
		{
			caseName:          "Successfully creates a transaction with TRANSFER transaction type.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transaction.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Successfully creates a transaction with DEPOSIT transaction type.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: nil,
			operationType:     transaction.Deposit,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Successfully creates a transaction with WITHDRAW transaction type.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: nil,
			operationType:     transaction.Withdraw,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           nil,
		},
		{
			caseName:          "Error occurs with unsupported transaction type.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     "UNSUPPORTED",
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           transaction.ErrUnsupportedType,
		},
		{
			caseName:          "Error occurs with invalid amount.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transaction.Transfer,
			amount:            -1000,
			currency:          currency,
			transactionAt:     transactionAt,
			wantErr:           money.ErrNegativeAmount,
		},
		{
			caseName:          "Error occurs with unsupported currency.",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transaction.Transfer,
			amount:            amount,
			currency:          "EUR",
			transactionAt:     transactionAt,
			wantErr:           money.ErrUnsupportedCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			tx, err := transaction.New(
				tt.id, tt.accountID, tt.receiverAccountID, tt.operationType, tt.amount, tt.currency, tt.transactionAt,
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
				assert.Equal(t, tt.operationType, tx.OperationType())
				assert.Equal(t, tt.amount, tx.TransferAmount().Amount())
				assert.Equal(t, tt.currency, tx.TransferAmount().Currency())
				assert.Equal(t, tt.transactionAt, tx.TransactionAt())
				assert.Equal(t, timer.GetFixedDateString(), tx.TransactionAtString())
			}
		})
	}
}
