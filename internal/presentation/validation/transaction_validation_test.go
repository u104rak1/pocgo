package validation_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/domain/transaction"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
)

func TestValidTransactionOperationType(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: Depositは有効",
			input:    transaction.Deposit,
			errMsg:   "",
		},
		{
			caseName: "Positive: Withdrawは有効",
			input:    transaction.Withdrawal,
			errMsg:   "",
		},
		{
			caseName: "Positive: Transferは有効",
			input:    transaction.Transfer,
			errMsg:   "",
		},
		{
			caseName: "Negative: 空文字列は無効",
			input:    "",
			errMsg:   "cannot be blank",
		},
		{
			caseName: "Negative: サポートされていない操作タイプは無効",
			input:    "invalid-type",
			errMsg:   "must be a valid value",
		},
		{
			caseName: "Negative: 小文字が含まれている場合は無効",
			input:    "dEPOSIT",
			errMsg:   "must be a valid value",
		},
		{
			caseName: "Negative: 全て小文字の場合は無効",
			input:    "deposit",
			errMsg:   "must be a valid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidTransactionOperationType(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidTransactionOperationTypes(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: 単一の有効な操作タイプ",
			input:    transaction.Deposit,
			errMsg:   "",
		},
		{
			caseName: "Positive: 複数の有効な操作タイプ",
			input:    transaction.Deposit + "," + transaction.Withdrawal,
			errMsg:   "",
		},
		{
			caseName: "Negative: 無効な操作タイプを含む",
			input:    transaction.Deposit + ",invalid-type",
			errMsg:   "operation_types: validation: must be a valid value",
		},
		{
			caseName: "Negative: 空文字列",
			input:    "",
			errMsg:   "operation_types: cannot be blank.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidTransactionOperationTypes(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidListTransactionsLimit(t *testing.T) {
	tests := []struct {
		caseName string
		input    int
		errMsg   string
	}{
		{
			caseName: "Positive: 有効なリミット",
			input:    10,
			errMsg:   "",
		},
		{
			caseName: "Negative: リミット0は無効",
			input:    0,
			errMsg:   "limit must be greater than 0",
		},
		{
			caseName: "Negative: リミットが大きすぎる場合は無効",
			input:    transaction.ListTransactionsLimit + 1,
			errMsg:   "limit must be less than or equal to " + strconv.Itoa(transaction.ListTransactionsLimit),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidListTransactionsLimit(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}
