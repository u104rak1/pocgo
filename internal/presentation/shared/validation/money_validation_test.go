package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	moneyVO "github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

func TestValidAmount(t *testing.T) {
	tests := []struct {
		caseName string
		currency string
		amount   float64
		wantErr  string
	}{
		{
			caseName: "A negative JPY is invalid.",
			currency: moneyVO.JPY,
			amount:   -1.0,
			wantErr:  "must be no less than 0",
		},
		{
			caseName: "0 JPY is valid.",
			currency: moneyVO.JPY,
			amount:   0.0,
			wantErr:  "",
		},
		{
			caseName: "A positive JPY is valid.",
			currency: moneyVO.JPY,
			amount:   1.0,
			wantErr:  "",
		},
		{
			caseName: "A JPY amount with invalid precision is invalid.",
			currency: moneyVO.JPY,
			amount:   100.5,
			wantErr:  moneyVO.ErrInvalidJPYPrecision.Error(),
		},
		{
			caseName: "A negative USD is invalid.",
			currency: moneyVO.USD,
			amount:   -1.0,
			wantErr:  "must be no less than 0",
		},
		{
			caseName: "0 USD is valid.",
			currency: moneyVO.USD,
			amount:   0.0,
			wantErr:  "",
		},
		{
			caseName: "A positive USD is valid.",
			currency: moneyVO.USD,
			amount:   1.0,
			wantErr:  "",
		},
		{
			caseName: "A USD is allowed up to 2 decimal places.",
			currency: moneyVO.USD,
			amount:   100.12,
			wantErr:  "",
		},
		{
			caseName: "A USD is not allowed more than 2 decimal places.",
			currency: moneyVO.USD,
			amount:   100.123,
			wantErr:  moneyVO.ErrInvalidUSDPrecision.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidAmount(tt.currency, tt.amount)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidCurrency(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  string
	}{
		{
			caseName: "An empty currency is invalid.",
			input:    "",
			wantErr:  "cannot be blank",
		},
		{
			caseName: "An unsupported currency is invalid.",
			input:    "EUR",
			wantErr:  "must be a valid value",
		},
		{
			caseName: "A valid currency (JPY) is accepted.",
			input:    moneyVO.JPY,
			wantErr:  "",
		},
		{
			caseName: "A valid currency (USD) is accepted.",
			input:    moneyVO.USD,
			wantErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidCurrency(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.wantErr)
			}
		})
	}
}
