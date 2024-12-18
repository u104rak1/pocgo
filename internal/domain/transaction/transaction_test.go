package transaction_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

func TestNewTransaction(t *testing.T) {

	var (
		transactionID     = ulid.GenerateStaticULID("transaction")
		accountID         = ulid.GenerateStaticULID("account")
		receiverAccountID = ulid.GenerateStaticULID("accountReceiver")
		amount            = 1000.0
		currency          = moneyVO.JPY
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
		errMsg            string
	}{
		{
			caseName:          "Positive: 取引を作成できる",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "",
		},
		{
			caseName:          "Negative: 無効なIDの場合はエラーが返る",
			id:                invalidID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "transaction id must be a valid ULID",
		},
		{
			caseName:          "Negative: 無効な口座IDの場合はエラーが返る",
			id:                transactionID,
			accountID:         "inavlid",
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            accountDomain.ErrInvalidID.Error(),
		},
		{
			caseName:          "Negative: 無効な受取口座IDの場合はエラーが返る",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &invalidID,
			operationType:     transactionDomain.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            accountDomain.ErrInvalidID.Error(),
		},
		{
			caseName:          "Positive: 振替取引を作成できる",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "",
		},
		{
			caseName:          "Positive: 入金取引を作成できる",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: nil,
			operationType:     transactionDomain.Deposit,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "",
		},
		{
			caseName:          "Positive: 出金取引を作成できる",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: nil,
			operationType:     transactionDomain.Withdraw,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "",
		},
		{
			caseName:          "Negative: サポートされていない取引タイプの場合はエラーが返る",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     "UNSUPPORTED",
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "unsupported transaction type",
		},
		{
			caseName:          "Negative: 無効な金額の場合はエラーが返る",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            -1000,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            moneyVO.ErrNegativeAmount.Error(),
		},
		{
			caseName:          "Negative: サポートされていない通貨の場合はエラーが返る",
			id:                transactionID,
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            amount,
			currency:          "EUR",
			transactionAt:     transactionAt,
			errMsg:            moneyVO.ErrUnsupportedCurrency.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			tx, err := transactionDomain.New(
				tt.id, tt.accountID, tt.receiverAccountID, tt.operationType, tt.amount, tt.currency, tt.transactionAt,
			)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
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
