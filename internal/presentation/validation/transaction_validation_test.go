package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/domain/transaction"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
)

func TestValidTransactionOperationType(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  string
	}{
		{
			caseName: "An empty operation type is invalid",
			input:    "",
			wantErr:  "cannot be blank",
		},
		{
			caseName: "An unsupported operation type is invalid",
			input:    "invalid-type",
			wantErr:  "must be a valid value",
		},
		{
			caseName: "A valid operation type (Deposit) is accepted",
			input:    transaction.Deposit,
			wantErr:  "",
		},
		{
			caseName: "A valid operation type (Withdraw) is accepted",
			input:    transaction.Withdraw,
			wantErr:  "",
		},
		{
			caseName: "A valid operation type (Transfer) is accepted",
			input:    transaction.Transfer,
			wantErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidTransactionOperationType(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.wantErr)
			}
		})
	}
}
