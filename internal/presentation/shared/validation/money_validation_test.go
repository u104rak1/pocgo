package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	moneyVO "github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

func TestValidJPY(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantErr string
	}{
		{
			"A negative amount is invalid.",
			-1.0,
			"must be no less than 0",
		},
		{
			"Zero is valid.",
			0.0,
			"",
		},
		{
			"A positive amount is valid.",
			1.0,
			"",
		},
		{
			"A JPY amount with invalid precision is invalid.",
			100.5,
			moneyVO.ErrInvalidJPYPrecision.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidJPY(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidUSD(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantErr string
	}{
		{
			"A negative amount is invalid.",
			-1.0,
			"must be no less than 0",
		},
		{
			"Zero is valid.",
			0.0,
			"",
		},
		{
			"A positive amount is valid.",
			1.0,
			"",
		},
		{
			"USD is allowed up to 2 decimal places.",
			100.12,
			"",
		},
		{
			"USD is not allowed more than 2 decimal places.",
			100.123,
			moneyVO.ErrInvalidUSDPrecision.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUSD(tt.input)
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
		name    string
		input   string
		wantErr string
	}{
		{
			"An empty currency is invalid.",
			"",
			"cannot be blank",
		},
		{
			"An unsupported currency is invalid.",
			"EUR",
			"must be a valid value",
		},
		{
			"A valid currency (JPY) is accepted.",
			moneyVO.JPY,
			"",
		},
		{
			"A valid currency (USD) is accepted.",
			moneyVO.USD,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
