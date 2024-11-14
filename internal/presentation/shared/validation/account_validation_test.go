package validation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

func TestValidAccountName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			"An empty name is invalid.",
			"",
			"cannot be blank",
		},
		{
			"A name less than 2-characters are invalid.",
			strings.Repeat("a", accountDomain.NameMinLength-1),
			"the length must be between 3 and 20",
		},
		{
			"Valid name must be a minimum 3-characters.",
			strings.Repeat("a", accountDomain.NameMinLength),
			"",
		},
		{
			"Valid name must be a maximum of 20-characters.",
			strings.Repeat("a", accountDomain.NameMaxLength),
			"",
		},
		{
			"A name longer than 21-characters are invalid.",
			strings.Repeat("a", accountDomain.NameMaxLength+1),
			"the length must be between 3 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidAccountName(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidAccountPassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			"An empty password is invalid.",
			"",
			"cannot be blank",
		},
		{
			"A password less than 3-characters are invalid.",
			strings.Repeat("a", accountDomain.PasswordLength-1),
			"the length must be exactly 4",
		},
		{
			"A password must be 4-characters.",
			strings.Repeat("a", accountDomain.PasswordLength),
			"",
		},
		{
			"A password longer than 5-characters are invalid.",
			strings.Repeat("a", accountDomain.PasswordLength+1),
			"the length must be exactly 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidAccountPassword(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}
