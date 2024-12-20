package transaction_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

func TestNewTransaction(t *testing.T) {
	var (
		accountID         = idVO.NewAccountIDForTest("account")
		receiverAccountID = idVO.NewAccountIDForTest("accountReceiver")
		amount            = 1000.0
		currency          = moneyVO.JPY
		transactionAt     = timer.GetFixedDate()
	)

	tests := []struct {
		caseName          string
		accountID         idVO.AccountID
		receiverAccountID *idVO.AccountID
		operationType     string
		amount            float64
		currency          string
		transactionAt     time.Time
		errMsg            string
	}{
		{
			caseName:          "Positive: 取引を作成できる",
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "",
		},
		{
			caseName:          "Positive: 振り込み取引を作成できる",
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
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     "UNSUPPORTED",
			amount:            amount,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            "unsupported transaction type",
		},
		{
			caseName:          "Negative: Money値オブジェクト作成時にエラーが返る場合はエラーが返る",
			accountID:         accountID,
			receiverAccountID: &receiverAccountID,
			operationType:     transactionDomain.Transfer,
			amount:            -1000,
			currency:          currency,
			transactionAt:     transactionAt,
			errMsg:            moneyVO.ErrNegativeAmount.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			tx, err := transactionDomain.New(
				tt.accountID, tt.receiverAccountID, tt.operationType, tt.amount, tt.currency, tt.transactionAt,
			)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, tx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tx)
				assert.NotEmpty(t, tx.ID())
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

func TestReconstruct(t *testing.T) {
	var (
		transactionID     = ulid.GenerateStaticULID("transaction")
		accountID         = ulid.GenerateStaticULID("account")
		receiverAccountID = ulid.GenerateStaticULID("accountReceiver")
		operationType     = "TRANSFER"
		amount            = 1000.0
		currency          = moneyVO.JPY
		transactionAt     = timer.GetFixedDate()
	)
	t.Run("Positive: 取引を再構築できる", func(t *testing.T) {
		tx, err := transactionDomain.Reconstruct(transactionID, accountID, &receiverAccountID, operationType, amount, currency, transactionAt)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.Equal(t, transactionID, tx.IDString())
		assert.Equal(t, accountID, tx.AccountIDString())
		assert.Equal(t, receiverAccountID, tx.ReceiverAccountIDString())
		assert.Equal(t, operationType, tx.OperationType())
		assert.Equal(t, amount, tx.TransferAmount().Amount())
		assert.Equal(t, currency, tx.TransferAmount().Currency())
		assert.Equal(t, transactionAt, tx.TransactionAt())
		assert.Equal(t, timer.GetFixedDateString(), tx.TransactionAtString())
	})
}
