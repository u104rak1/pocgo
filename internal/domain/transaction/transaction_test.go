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
			caseName:          "Happy path: 有効なTransactionを作成",
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
			caseName:          "Edge case: 無効なIDを指定するとエラー",
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
			caseName:          "Edge case: 無効なAccountIDを指定するとエラー",
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
			caseName:          "Edge case: 無効なReceiverAccountIDを指定するとエラー",
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
			caseName:          "Happy path: Transaction typeがTRANSFERなら成功",
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
			caseName:          "Happy path: Transaction typeがDEPOSITなら成功",
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
			caseName:          "Happy path: Transaction typeがWITHDRAWなら成功",
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
			caseName:          "Happy path: サポートされていないTransaction typeは失敗",
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
			caseName:          "Edge case: 無効な金額を指定するとエラー",
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
			caseName:          "Edge case: サポートされていない通貨を指定するとエラー",
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
