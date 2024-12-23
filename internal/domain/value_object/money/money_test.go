package money_test

import (
	"github.com/stretchr/testify/assert"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"

	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		errMsg   string
	}{
		{
			name:     "Positive: 有効なJPYの場合は、金額が作成できる",
			amount:   1000,
			currency: "JPY",
			errMsg:   "",
		},
		{
			name:     "Negative: 小数点以下のJPYの場合はエラーが返る",
			amount:   1000.1,
			currency: "JPY",
			errMsg:   "amount in JPY must not have decimal places",
		},
		{
			name:     "Positive: 有効なUSDの場合は、金額が作成できる",
			amount:   10,
			currency: "USD",
			errMsg:   "",
		},
		{
			name:     "Positive: 小数点第2位までのUSDの場合は、金額が作成できる",
			amount:   10.99,
			currency: "USD",
			errMsg:   "",
		},
		{
			name:     "Negative: 小数点第3位のUSDの場合はエラーが返る",
			amount:   10.001,
			currency: "USD",
			errMsg:   "amount in USD cannot have more than 2 decimal places",
		},
		{
			name:     "Negative: 金額がマイナスの場合はエラーが返る",
			amount:   -1,
			currency: "JPY",
			errMsg:   "amount cannot be negative",
		},
		{
			name:     "Negative: サポートされていない通貨の場合はエラーが返る",
			amount:   1000,
			currency: "EUR",
			errMsg:   "unsupported currency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, err := moneyVO.New(tt.amount, tt.currency)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, m)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, m)
				assert.Equal(t, tt.amount, m.Amount())
				assert.Equal(t, tt.currency, m.Currency())
			}
		})
	}
}

func TestAdd(t *testing.T) {
	m1, _ := moneyVO.New(1000, "JPY")
	m2, _ := moneyVO.New(500, "JPY")
	m3, _ := moneyVO.New(10.50, "USD")

	tests := []struct {
		name   string
		money1 *moneyVO.Money
		money2 *moneyVO.Money
		want   float64
		errMsg string
	}{
		{
			name:   "Positive: 通貨が同じ場合は、金額が加算できる",
			money1: m1,
			money2: m2,
			want:   1500,
			errMsg: "",
		},
		{
			name:   "Negative: 通貨が異なる場合はエラーが返る",
			money1: m1,
			money2: m3,
			want:   0,
			errMsg: "operation cannot be performed on different currencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.money1.Add(*tt.money2)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result.Amount())
			}
		})
	}
}

func TestSub(t *testing.T) {
	m1, _ := moneyVO.New(1000, "JPY")
	m2, _ := moneyVO.New(500, "JPY")
	m3, _ := moneyVO.New(2000, "JPY")
	m4, _ := moneyVO.New(10.50, "USD")

	tests := []struct {
		name   string
		money1 *moneyVO.Money
		money2 *moneyVO.Money
		want   float64
		errMsg string
	}{
		{
			name:   "Positive: 通貨が同じ場合は、金額が減算できる",
			money1: m1,
			money2: m2,
			want:   500,
			errMsg: "",
		},
		{
			name:   "Negative: 金額が足りない場合はエラーが返る",
			money1: m1,
			money2: m3,
			want:   0,
			errMsg: "insufficient balance",
		},
		{
			name:   "Negative: 通貨が異なる場合はエラーが返る",
			money1: m1,
			money2: m4,
			want:   0,
			errMsg: "operation cannot be performed on different currencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.money1.Sub(*tt.money2)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result.Amount())
			}
		})
	}
}
