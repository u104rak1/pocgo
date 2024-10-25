package validation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

func TestValidUserName(t *testing.T) {
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
			strings.Repeat("a", userDomain.NameMinLength-1),
			"the length must be between 3 and 20",
		},
		{
			"Valid name must be a minimum 3-characters.",
			strings.Repeat("a", userDomain.NameMinLength),
			"",
		},
		{
			"Valid name must be a maximum of 20-characters.",
			strings.Repeat("a", userDomain.NameMaxLength),
			"",
		},
		{
			"A name longer than 21-characters are invalid.",
			strings.Repeat("a", userDomain.NameMaxLength+1),
			"the length must be between 3 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserName(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidUserEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			"An empty email is invalid.",
			"",
			"cannot be blank",
		},
		{
			"An invalid email format is invalid.",
			"invalid",
			"the email format is invalid",
		},
		{
			"A valid email format is valid.",
			"test@example.com",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserEmail(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidUserPassword(t *testing.T) {
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
			"A password less than 7-characters are invalid.",
			strings.Repeat("a", authDomain.PasswordMinLength-1),
			"the length must be between 8 and 20",
		},
		{
			"Valid password must be a minimum 8-characters.",
			strings.Repeat("a", authDomain.PasswordMinLength),
			"",
		},
		{
			"Valid password must be a maximum of 20-characters.",
			strings.Repeat("a", authDomain.PasswordMaxLength),
			"",
		},
		{
			"A password longer than 21-characters are invalid.",
			strings.Repeat("a", authDomain.PasswordMaxLength+1),
			"the length must be between 8 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserPassword(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}
