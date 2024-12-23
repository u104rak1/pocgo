package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
)

func TestValidAmount(t *testing.T) {
	tests := []struct {
		caseName string
		currency string
		amount   float64
		errMsg   string
	}{
		{
			caseName: "Negative: 負のJPYは無効",
			currency: moneyVO.JPY,
			amount:   -1.0,
			errMsg:   "must be no less than 0",
		},
		{
			caseName: "Positive: 0 JPYは有効",
			currency: moneyVO.JPY,
			amount:   0.0,
			errMsg:   "",
		},
		{
			caseName: "Positive: 正のJPYは有効",
			currency: moneyVO.JPY,
			amount:   1.0,
			errMsg:   "",
		},
		{
			caseName: "Negative: JPYの精度が無効",
			currency: moneyVO.JPY,
			amount:   100.5,
			errMsg:   moneyVO.ErrInvalidJPYPrecision.Error(),
		},
		{
			caseName: "Negative: 負のUSDは無効",
			currency: moneyVO.USD,
			amount:   -1.0,
			errMsg:   "must be no less than 0",
		},
		{
			caseName: "Positive: 0 USDは有効",
			currency: moneyVO.USD,
			amount:   0.0,
			errMsg:   "",
		},
		{
			caseName: "Positive: 正のUSDは有効",
			currency: moneyVO.USD,
			amount:   1.0,
			errMsg:   "",
		},
		{
			caseName: "Positive: USDは2桁まで有効",
			currency: moneyVO.USD,
			amount:   100.12,
			errMsg:   "",
		},
		{
			caseName: "Negative: USDは3桁以上は無効",
			currency: moneyVO.USD,
			amount:   100.123,
			errMsg:   moneyVO.ErrInvalidUSDPrecision.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidAmount(tt.currency, tt.amount)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidCurrency(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: JPYは有効",
			input:    moneyVO.JPY,
			errMsg:   "",
		},
		{
			caseName: "Positive: USDは有効",
			input:    moneyVO.USD,
			errMsg:   "",
		},
		{
			caseName: "Negative: 空文字列は無効",
			input:    "",
			errMsg:   "cannot be blank",
		},
		{
			caseName: "Negative: サポートされていない通貨は無効",
			input:    "EUR",
			errMsg:   "must be a valid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidCurrency(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}
