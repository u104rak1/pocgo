package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/transaction"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

func TestValidTransactionOperationType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			"An empty operation type is invalid",
			"",
			"cannot be blank",
		},
		{
			"An unsupported operation type is invalid",
			"invalid-type",
			"must be a valid value",
		},
		{
			"A valid operation type (Deposit) is accepted",
			transaction.Deposit,
			"",
		},
		{
			"A valid operation type (Withdraw) is accepted",
			transaction.Withdraw,
			"",
		},
		{
			"A valid operation type (Transfer) is accepted",
			transaction.Transfer,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
