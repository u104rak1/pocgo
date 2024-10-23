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
		// JPY specific tests
		{
			name:     "Successfully creates a money value object, if valid JPY.",
			amount:   1000,
			currency: money.JPY,
			wantErr:  nil,
		},
		{
			name:     "Error occurs with invalid JPY with 1 decimal point.",
			amount:   1000.1,
			currency: money.JPY,
			wantErr:  money.ErrInvalidJPYPrecision,
		},
		// USD specific tests
		{
			name:     "Successfully creates a money value object, if valid USD.",
			amount:   10,
			currency: money.USD,
			wantErr:  nil,
		},
		{
			name:     "Successfully creates a money value object, if valid USD with 2 decimal points.",
			amount:   10.99,
			currency: money.USD,
			wantErr:  nil,
		},
		{
			name:     "Error occurs with invalid USD with 3 decimal points.",
			amount:   10.001,
			currency: money.USD,
			wantErr:  money.ErrInvalidUSDPrecision,
		},
		// Common tests
		{
			name:     "Error occurs with negative amount.",
			amount:   -1,
			currency: money.JPY,
			wantErr:  money.ErrNegativeAmount,
		},
		{
			name:     "Error occurs with unsupported currency.",
			amount:   1000,
			currency: "EUR",
			wantErr:  money.ErrUnsupportedCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, err := money.New(tt.amount, tt.currency)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
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
			name:    "Successfully adds two money, if the currency is the same.",
			money1:  m1,
			money2:  m2,
			want:    1500,
			wantErr: nil,
		},
		{
			name:    "Error occurs with different currency.",
			money1:  m1,
			money2:  m3,
			want:    0,
			wantErr: money.ErrDifferentCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.money1.Add(*tt.money2)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result.Amount())
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
			name:    "Successfully subtracts two money, if the currency is the same.",
			money1:  m1,
			money2:  m2,
			want:    500,
			wantErr: nil,
		},
		{
			name:    "Error occurs with insufficient balance.",
			money1:  m1,
			money2:  m3,
			want:    0,
			wantErr: money.ErrInsufficientBalance,
		},
		{
			name:    "Error occurs with different currency.",
			money1:  m1,
			money2:  m4,
			want:    0,
			wantErr: money.ErrDifferentCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.money1.Sub(*tt.money2)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result.Amount())
			}
		})
	}
}
