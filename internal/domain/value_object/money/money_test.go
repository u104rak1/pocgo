package money_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"

	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		wantErr  error
	}{
		// 円固有のテスト
		{
			name:     "Happy path: 有効な日本円のMoneyを作成",
			amount:   1000,
			currency: money.JPY,
			wantErr:  nil,
		},
		{
			name:     "Edge case: 小数点を含む無効な日本円はエラー",
			amount:   1000.1,
			currency: money.JPY,
			wantErr:  money.ErrInvalidJPYPrecision,
		},
		// 米ドル固有のテスト
		{
			name:     "Happy path: 有効な米ドルのMoneyを作成",
			amount:   10,
			currency: money.USD,
			wantErr:  nil,
		},
		{
			name:     "Happy path: 小数点2桁の有効な米ドルのMoneyを作成",
			amount:   10.99,
			currency: money.USD,
			wantErr:  nil,
		},
		{
			name:     "Edge case: 小数点3桁の無効な米ドルはエラー",
			amount:   10.001,
			currency: money.USD,
			wantErr:  money.ErrInvalidUSDPrecision,
		},
		// 共通のテスト
		{
			name:     "Edge case: 金額がマイナスを指定するとエラー",
			amount:   -1,
			currency: money.JPY,
			wantErr:  money.ErrNegativeAmount,
		},
		{
			name:     "Edge case: サポートされていない通貨を指定するとエラー",
			amount:   1000,
			currency: "EUR",
			wantErr:  money.ErrUnsupportedCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := money.New(tt.amount, tt.currency)
			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, m)
				assert.Equal(t, tt.amount, m.Amount())
				assert.Equal(t, tt.currency, m.Currency())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, m)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	m1, _ := money.New(1000, money.JPY)
	m2, _ := money.New(500, money.JPY)
	m3, _ := money.New(10.50, money.USD)

	tests := []struct {
		name    string
		money1  *money.Money
		money2  *money.Money
		want    float64
		wantErr error
	}{
		{
			name:    "Happy path: 通貨が同じなら加算できる",
			money1:  m1,
			money2:  m2,
			want:    1500,
			wantErr: nil,
		},
		{
			name:    "Edge case: 異なる通貨はエラー",
			money1:  m1,
			money2:  m3,
			want:    0,
			wantErr: money.ErrDifferentCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money1.Add(*tt.money2)
			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result.Amount())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			}
		})
	}
}

func TestSub(t *testing.T) {
	m1, _ := money.New(1000, money.JPY)
	m2, _ := money.New(500, money.JPY)
	m3, _ := money.New(2000, money.JPY)
	m4, _ := money.New(10.50, money.USD)

	tests := []struct {
		name    string
		money1  *money.Money
		money2  *money.Money
		want    float64
		wantErr error
	}{
		{
			name:    "Happy path: 通貨が同じで金額が足りる場合は減算できる",
			money1:  m1,
			money2:  m2,
			want:    500,
			wantErr: nil,
		},
		{
			name:    "Edge case: 金額が足りない場合はエラー",
			money1:  m1,
			money2:  m3,
			want:    0,
			wantErr: money.ErrInsufficientBalance,
		},
		{
			name:    "Edge case: 異なる通貨はエラー",
			money1:  m1,
			money2:  m4,
			want:    0,
			wantErr: money.ErrDifferentCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money1.Sub(*tt.money2)
			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result.Amount())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			}
		})
	}
}
